// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Engine is the semantic audio director. Its cue methods are safe to call from
// any goroutine and do not perform device I/O.
type Engine struct {
	mu      sync.Mutex
	backend Backend
	catalog Catalog
	options Options
	mixer   *mixer
	random  *rand.Rand
	closed  bool

	musicVoice  voiceID
	music       Music
	ambienceBed voiceID
	continuous  map[Continuous]voiceID
	active      map[SoundID][]voiceID
	lastPlayed  map[SoundID]time.Time
	environment Environment
	lastAmbient SoundID

	ambienceWake chan struct{}
	done         chan struct{}
	wg           sync.WaitGroup
}

// New constructs an engine and starts its backend. Catalog clips must already
// be decoded into the requested output format.
func New(backend Backend, catalog Catalog, options Options) (*Engine, error) {
	if backend == nil {
		return nil, fmt.Errorf("audio: nil backend")
	}
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	options = options.normalized()
	format := Format{SampleRate: options.SampleRate, Channels: options.Channels}
	if format != catalog.Format {
		return nil, fmt.Errorf("audio: catalog format %+v does not match output format %+v", catalog.Format, format)
	}

	mix := newMixer(format)
	mix.setBusGain(BusMusic, options.MusicGain)
	mix.setBusGain(BusAmbience, options.AmbienceGain)
	mix.setBusGain(BusEffects, options.EffectsGain)
	mix.setBusGain(BusUI, options.UIGain)

	e := &Engine{
		backend:      backend,
		catalog:      catalog,
		options:      options,
		mixer:        mix,
		random:       rand.New(rand.NewSource(options.RandomSeed)),
		continuous:   make(map[Continuous]voiceID),
		active:       make(map[SoundID][]voiceID),
		lastPlayed:   make(map[SoundID]time.Time),
		ambienceWake: make(chan struct{}, 1),
		done:         make(chan struct{}),
	}
	if err := backend.Start(&pcmReader{mixer: mix}, format); err != nil {
		return nil, fmt.Errorf("audio: start backend: %w", err)
	}
	e.wg.Add(1)
	go e.runAmbience()
	return e, nil
}

// Play emits a one-shot gameplay cue. Missing catalog entries are silent so a
// partially populated sound library remains usable during development.
func (e *Engine) Play(cue Cue) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	if id, ok := e.catalog.Cues[cue]; ok {
		e.playSoundLocked(id, 0)
	}
	return nil
}

// SetMusic crossfades to a semantic music program. MusicNone fades to silence.
func (e *Engine) SetMusic(music Music) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	if e.music == music && (music == MusicNone || e.mixer.active(e.musicVoice)) {
		return nil
	}
	fade := e.frames(e.options.MusicFade)
	if e.musicVoice != 0 {
		e.mixer.stop(e.musicVoice, fade)
		e.musicVoice = 0
	}
	if music == MusicNone {
		e.music = MusicNone
		return nil
	}
	if id, ok := e.catalog.Music[music]; ok {
		e.musicVoice = e.playSoundLocked(id, fade)
		e.music = music
	}
	return nil
}

// SetContinuous starts or stops a stateful sound such as the burner or rain.
// Repeating the current state is harmless.
func (e *Engine) SetContinuous(cue Continuous, enabled bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	voice := e.continuous[cue]
	if enabled {
		if voice != 0 && e.mixer.active(voice) {
			return nil
		}
		if id, ok := e.catalog.Continuous[cue]; ok {
			e.continuous[cue] = e.playSoundLocked(id, e.frames(e.options.ContinuousFade))
		}
		return nil
	}
	if voice != 0 {
		e.mixer.stop(voice, e.frames(e.options.ContinuousFade))
		delete(e.continuous, cue)
	}
	return nil
}

// SetEnvironment replaces the ambient bed and restarts the environment's
// sparse randomized sound scheduler.
func (e *Engine) SetEnvironment(environment Environment) error {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return ErrClosed
	}
	if e.environment == environment {
		e.mu.Unlock()
		return nil
	}
	fade := e.frames(e.options.MusicFade)
	if e.ambienceBed != 0 {
		e.mixer.stop(e.ambienceBed, fade)
		e.ambienceBed = 0
	}
	e.environment = environment
	e.lastAmbient = ""
	if scene, ok := e.catalog.Ambience[environment]; ok && scene.Bed != "" {
		e.ambienceBed = e.playSoundLocked(scene.Bed, fade)
	}
	e.mu.Unlock()
	e.wakeAmbience()
	return nil
}

// SetBusGain changes a user-facing mix control. Values above one are allowed,
// although the final stream is clipped to [-1,1]. Negative and NaN values mute.
func (e *Engine) SetBusGain(bus Bus, gain float32) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	e.mixer.setBusGain(bus, gain)
	return nil
}

// Suspend releases or pauses device playback while preserving engine state.
func (e *Engine) Suspend() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	return e.backend.Suspend()
}

// Resume is also the browser audio unlock operation and should be called from
// the first user gesture on WebAssembly.
func (e *Engine) Resume() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return ErrClosed
	}
	return e.backend.Resume()
}

// Close stops scheduling and closes the backend. It is idempotent.
func (e *Engine) Close() error {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return nil
	}
	e.closed = true
	close(e.done)
	e.mu.Unlock()
	e.wg.Wait()
	return e.backend.Close()
}

func (e *Engine) playSoundLocked(id SoundID, fadeFrames int) voiceID {
	sound, ok := e.catalog.Sounds[id]
	if !ok || len(sound.Variants) == 0 {
		return 0
	}
	now := time.Now()
	if sound.Cooldown > 0 && now.Sub(e.lastPlayed[id]) < sound.Cooldown {
		return 0
	}
	voices := e.active[id][:0]
	for _, voice := range e.active[id] {
		if e.mixer.active(voice) {
			voices = append(voices, voice)
		}
	}
	e.active[id] = voices
	if sound.MaxVoices > 0 && len(voices) >= sound.MaxVoices {
		return 0
	}
	clip := sound.Variants[e.random.Intn(len(sound.Variants))]
	voice := e.mixer.play(clip, sound.Bus, sound.Gain, sound.Loop, fadeFrames)
	e.active[id] = append(voices, voice)
	e.lastPlayed[id] = now
	return voice
}

func (e *Engine) runAmbience() {
	defer e.wg.Done()
	for {
		delay, available := e.nextAmbienceDelay()
		if !available {
			select {
			case <-e.ambienceWake:
				continue
			case <-e.done:
				return
			}
		}
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			e.playAmbientDetail()
		case <-e.ambienceWake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		case <-e.done:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

func (e *Engine) nextAmbienceDelay() (time.Duration, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return 0, false
	}
	scene, ok := e.catalog.Ambience[e.environment]
	if !ok || len(scene.Random) == 0 {
		return 0, false
	}
	span := scene.MaxDelay - scene.MinDelay
	if span <= 0 {
		return scene.MinDelay, true
	}
	return scene.MinDelay + time.Duration(e.random.Int63n(int64(span)+1)), true
}

func (e *Engine) playAmbientDetail() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return
	}
	scene, ok := e.catalog.Ambience[e.environment]
	if !ok || len(scene.Random) == 0 {
		return
	}
	choices := scene.Random
	index := e.random.Intn(len(choices))
	if len(choices) > 1 && choices[index] == e.lastAmbient {
		// Choose uniformly from all other entries without allocating.
		index = (index + 1 + e.random.Intn(len(choices)-1)) % len(choices)
	}
	id := choices[index]
	e.lastAmbient = id
	e.playSoundLocked(id, 0)
}

func (e *Engine) wakeAmbience() {
	select {
	case e.ambienceWake <- struct{}{}:
	default:
	}
}

func (e *Engine) frames(duration time.Duration) int {
	return int(duration * time.Duration(e.options.SampleRate) / time.Second)
}
