// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

// Package sound provides the game's embedded CC0 sound library.
package sound

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/audio"
)

const (
	soundClear         audio.SoundID = "music.clear"
	soundThreat        audio.SoundID = "music.threat"
	soundFailure       audio.SoundID = "music.failure"
	soundVictory       audio.SoundID = "music.victory"
	soundJungle        audio.SoundID = "ambience.jungle"
	soundStorm         audio.SoundID = "ambience.storm"
	soundBurner        audio.SoundID = "continuous.burner"
	soundBasket        audio.SoundID = "effect.basket-creak"
	soundBasketAmbient audio.SoundID = "ambience.basket-creak"
	soundPop           audio.SoundID = "effect.balloon-pop"
	soundDinosaur1     audio.SoundID = "effect.dinosaur-distant-1"
	soundDinosaur2     audio.SoundID = "effect.dinosaur-distant-2"
	soundSauropod      audio.SoundID = "effect.sauropod-distant"
	soundCarnivore     audio.SoundID = "effect.carnivore"
	soundPterosaur     audio.SoundID = "effect.pterosaur"
	soundPterosaurFar  audio.SoundID = "ambience.pterosaur-distant"
	soundLightning     audio.SoundID = "effect.lightning-crack"
	soundThunderFar    audio.SoundID = "effect.thunder-distant"
	soundThunder       audio.SoundID = "effect.thunder-near"
	soundPickup        audio.SoundID = "effect.pickup"
	soundConfirm       audio.SoundID = "ui.confirm"
)

// Catalog returns the audio assets and semantic cue mappings used by the game.
func Catalog() audio.Catalog {
	format := audio.Format{SampleRate: audio.DefaultSampleRate, Channels: audio.DefaultChannels}
	return audio.Catalog{
		Format: format,
		Sounds: map[audio.SoundID]audio.Sound{
			soundClear: {
				Variants: []audio.Clip{recorded("music_clear.ogg")},
				Bus:      audio.BusMusic, Gain: 0.72, Loop: true,
			},
			soundThreat: {
				Variants: []audio.Clip{recorded("music_threat.ogg")},
				Bus:      audio.BusMusic, Gain: 0.65, Loop: true,
			},
			soundFailure: {
				Variants: []audio.Clip{recorded("music_failure.ogg")},
				Bus:      audio.BusMusic, Gain: 0.75,
			},
			soundVictory: {
				Variants: []audio.Clip{recorded("music_victory.ogg")},
				Bus:      audio.BusMusic, Gain: 0.75,
			},
			soundJungle: {
				Variants: []audio.Clip{recorded("ambience_jungle.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.28, Loop: true,
			},
			soundStorm: {
				Variants: []audio.Clip{recorded("ambience_rain.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.45, Loop: true,
			},
			soundBurner: {
				Variants: []audio.Clip{{PCM: burnerLoop()}},
				Bus:      audio.BusEffects, Loop: true,
			},
			soundBasket: {
				Variants: []audio.Clip{
					recorded("effect_basket_1.ogg"),
					recorded("effect_basket_2.ogg"),
				},
				Bus: audio.BusEffects, Gain: 0.55, MaxVoices: 1, Cooldown: 150 * time.Millisecond,
			},
			soundBasketAmbient: {
				Variants: []audio.Clip{
					recorded("effect_basket_1.ogg"),
					recorded("effect_basket_2.ogg"),
				},
				Bus: audio.BusAmbience, Gain: 0.18, MaxVoices: 1, Cooldown: 2 * time.Second,
			},
			soundPop: {
				Variants: []audio.Clip{recorded("effect_balloon_pop.ogg")},
				Bus:      audio.BusEffects, Gain: 0.8, MaxVoices: 1,
			},
			soundDinosaur1: {
				Variants: []audio.Clip{recorded("effect_dinosaur_1.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.2, MaxVoices: 1, Cooldown: 2 * time.Second,
			},
			soundDinosaur2: {
				Variants: []audio.Clip{recorded("effect_dinosaur_2.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.18, MaxVoices: 1, Cooldown: 2 * time.Second,
			},
			soundSauropod: {
				Variants: []audio.Clip{recorded("effect_sauropod.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.2, MaxVoices: 1, Cooldown: 3 * time.Second,
			},
			soundCarnivore: {
				Variants: []audio.Clip{recorded("effect_carnivore.ogg")},
				Bus:      audio.BusEffects, Gain: 0.6, MaxVoices: 1,
			},
			soundPterosaur: {
				Variants: []audio.Clip{
					{PCM: crowSquawk(760, 330, 320*time.Millisecond, 0.42, 29)},
					{PCM: crowSquawk(860, 380, 270*time.Millisecond, 0.40, 47)},
				},
				Bus: audio.BusEffects, Gain: 0.78, MaxVoices: 1, Cooldown: 400 * time.Millisecond,
			},
			soundPterosaurFar: {
				Variants: []audio.Clip{
					{PCM: crowSquawk(440, 210, 380*time.Millisecond, 0.26, 71)},
					{PCM: crowSquawk(485, 225, 330*time.Millisecond, 0.24, 83)},
				},
				Bus: audio.BusAmbience, Gain: 0.12, MaxVoices: 1, Cooldown: 3 * time.Second,
			},
			soundLightning: {
				Variants: []audio.Clip{recorded("effect_lightning_crack.ogg")},
				Bus:      audio.BusEffects, Gain: 0.95, MaxVoices: 2,
			},
			soundThunderFar: {
				Variants: []audio.Clip{recorded("effect_thunder.ogg")},
				Bus:      audio.BusAmbience, Gain: 0.4, MaxVoices: 1, Cooldown: 2 * time.Second,
			},
			soundThunder: {
				Variants: []audio.Clip{recorded("effect_thunder.ogg")},
				Bus:      audio.BusEffects, Gain: 0.85, MaxVoices: 1,
			},
			soundPickup: {
				Variants: []audio.Clip{{PCM: chirp(440, 1100, 250*time.Millisecond, 0.25)}},
				Bus:      audio.BusEffects, MaxVoices: 2,
			},
			soundConfirm: {
				Variants: []audio.Clip{{PCM: tone(660, 90*time.Millisecond, 0.18)}},
				Bus:      audio.BusUI, MaxVoices: 2,
			},
		},
		Music: map[audio.Music]audio.SoundID{
			audio.MusicClearSkies: soundClear,
			audio.MusicThreat:     soundThreat,
			audio.MusicFailure:    soundFailure,
			audio.MusicVictory:    soundVictory,
		},
		Continuous: map[audio.Continuous]audio.SoundID{
			audio.ContinuousBurner: soundBurner,
		},
		Cues: map[audio.Cue]audio.SoundID{
			audio.CueBalloonPop:       soundPop,
			audio.CuePterosaurCall:    soundPterosaur,
			audio.CueDistantDinosaur:  soundDinosaur1,
			audio.CueDistantSauropod:  soundSauropod,
			audio.CueDistantCarnivore: soundCarnivore,
			audio.CueThunderDistant:   soundThunderFar,
			audio.CueThunderNear:      soundThunder,
			audio.CuePickup:           soundPickup,
			audio.CueMenuConfirm:      soundConfirm,
			audio.CueBasketCreak:      soundBasket,
			audio.CueLightningCrack:   soundLightning,
		},
		Ambience: map[audio.Environment]audio.Soundscape{
			audio.EnvironmentPrehistoricDay: {
				Bed:      soundJungle,
				Random:   []audio.SoundID{soundBasketAmbient},
				MinDelay: 8 * time.Second, MaxDelay: 18 * time.Second,
			},
			audio.EnvironmentStorm: {
				Bed: soundStorm,
			},
		},
	}
}

func melody(notes []float64, noteLength time.Duration, gain float64) []float32 {
	framesPerNote := durationFrames(noteLength)
	pcm := make([]float32, len(notes)*framesPerNote*2)
	for noteIndex, frequency := range notes {
		for frame := range framesPerNote {
			phase := float64(frame) / audio.DefaultSampleRate
			envelope := edgeEnvelope(frame, framesPerNote, 0.08)
			sample := (math.Sin(2*math.Pi*frequency*phase) + 0.28*math.Sin(4*math.Pi*frequency*phase)) * gain * envelope
			putStereo(pcm, noteIndex*framesPerNote+frame, sample)
		}
	}
	return pcm
}

func threatLoop() []float32 {
	frames := durationFrames(2 * time.Second)
	pcm := make([]float32, frames*2)
	for frame := range frames {
		t := float64(frame) / audio.DefaultSampleRate
		pulse := 0.35 + 0.65*math.Pow(0.5+0.5*math.Sin(2*math.Pi*3*t), 4)
		sample := (math.Sin(2*math.Pi*73.42*t) + 0.45*math.Sin(2*math.Pi*77.78*t)) * 0.14 * pulse
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func victory() []float32 {
	notes := []float64{261.63, 329.63, 392, 523.25}
	pcm := melody(notes, 150*time.Millisecond, 0.24)
	chordFrames := durationFrames(600 * time.Millisecond)
	chord := make([]float32, chordFrames*2)
	for frame := range chordFrames {
		t := float64(frame) / audio.DefaultSampleRate
		envelope := math.Exp(-2.2 * t)
		sample := 0.11 * envelope * (math.Sin(2*math.Pi*261.63*t) + math.Sin(2*math.Pi*329.63*t) + math.Sin(2*math.Pi*392*t))
		putStereo(chord, frame, sample)
	}
	return append(pcm, chord...)
}

func tone(frequency float64, duration time.Duration, gain float64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	for frame := range frames {
		t := float64(frame) / audio.DefaultSampleRate
		putStereo(pcm, frame, math.Sin(2*math.Pi*frequency*t)*gain*edgeEnvelope(frame, frames, 0.15))
	}
	return pcm
}

func chirp(start, end float64, duration time.Duration, gain float64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	phase := 0.0
	for frame := range frames {
		progress := float64(frame) / float64(frames)
		frequency := start + (end-start)*progress
		phase += 2 * math.Pi * frequency / audio.DefaultSampleRate
		putStereo(pcm, frame, math.Sin(phase)*gain*edgeEnvelope(frame, frames, 0.12))
	}
	return pcm
}

func crowSquawk(start, end float64, duration time.Duration, gain float64, seed int64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(seed))
	phase := 0.0
	rasp := 0.0
	for frame := range frames {
		progress := float64(frame) / float64(frames)
		frequency := start * math.Pow(end/start, progress)
		frequency *= 1 + 0.035*math.Sin(progress*math.Pi*13)
		phase += 2 * math.Pi * frequency / audio.DefaultSampleRate
		rasp += (random.Float64()*2 - 1 - rasp) * 0.3

		// A crow-like call has a hard onset, a short falling throat tone,
		// rough harmonics, and a noisy tail rather than a long clean scream.
		attack := min(progress/0.015, 1)
		release := min((1-progress)/0.08, 1)
		decay := math.Exp(-2.8 * progress)
		throat := 0.45*math.Sin(phase) + 0.30*math.Sin(phase*2.07) + 0.18*math.Sin(phase*3.91) + 0.10*math.Sin(phase*6.2)
		roughness := 0.7 + 0.3*math.Sin(progress*math.Pi*17)
		tail := math.Pow(progress, 1.7) * rasp * 0.32
		onset := rasp * math.Exp(-22*progress) * 0.45
		sample := (throat*roughness + onset + tail) * gain * attack * decay * release
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func noiseBed(duration time.Duration, gain, smoothing float64, seed int64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(seed))
	filtered := 0.0
	for frame := range frames {
		filtered += (random.Float64()*2 - 1 - filtered) * smoothing
		wobble := 0.65 + 0.35*math.Sin(2*math.Pi*0.17*float64(frame)/audio.DefaultSampleRate)
		putStereo(pcm, frame, filtered*gain*wobble)
	}
	return pcm
}

func burnerLoop() []float32 {
	frames := durationFrames(900 * time.Millisecond)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(17))
	filtered := 0.0
	for frame := range frames {
		filtered += (random.Float64()*2 - 1 - filtered) * 0.08
		t := float64(frame) / audio.DefaultSampleRate
		sample := filtered*0.25 + math.Sin(2*math.Pi*47*t)*0.035
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func basketCreak(start, end float64, duration time.Duration, gain float64, seed int64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(seed))
	phase := 0.0
	friction := 0.0
	for frame := range frames {
		progress := float64(frame) / float64(frames)
		frequency := start * math.Pow(end/start, progress)
		frequency *= 1 + 0.035*math.Sin(progress*math.Pi*7)
		phase += 2 * math.Pi * frequency / audio.DefaultSampleRate
		friction += (random.Float64()*2 - 1 - friction) * 0.12

		// Uneven stick-slip pulses keep the result closer to wicker and rope
		// under strain than a clean descending tone.
		stickSlip := 0.45 + 0.55*math.Pow(math.Abs(math.Sin(progress*math.Pi*5.5)), 3)
		envelope := math.Pow(math.Sin(math.Pi*progress), 0.55)
		sample := (0.62*math.Sin(phase) + 0.22*math.Sin(phase*2.7) + 0.22*friction) * gain * stickSlip * envelope
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func growl(start, end float64, duration time.Duration, gain float64, seed int64) []float32 {
	frames := durationFrames(duration)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(seed))
	phase := 0.0
	noise := 0.0
	for frame := range frames {
		progress := float64(frame) / float64(frames)
		frequency := start + (end-start)*progress + 4*math.Sin(progress*math.Pi*5)
		phase += 2 * math.Pi * frequency / audio.DefaultSampleRate
		noise += (random.Float64()*2 - 1 - noise) * 0.04
		envelope := math.Sin(math.Pi * progress)
		sample := (math.Sin(phase) + 0.45*math.Sin(phase*2.03) + noise*0.7) * gain * envelope
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func pop(seed int64) []float32 {
	frames := durationFrames(180 * time.Millisecond)
	pcm := make([]float32, frames*2)
	random := rand.New(rand.NewSource(seed))
	for frame := range frames {
		progress := float64(frame) / float64(frames)
		envelope := math.Exp(-18 * progress)
		sample := ((random.Float64()*2-1)*0.8 + math.Sin(2*math.Pi*110*float64(frame)/audio.DefaultSampleRate)*0.3) * envelope
		putStereo(pcm, frame, sample)
	}
	return pcm
}

func durationFrames(duration time.Duration) int {
	return int(duration * audio.DefaultSampleRate / time.Second)
}

func edgeEnvelope(frame, frames int, fraction float64) float64 {
	edge := max(1, int(float64(frames)*fraction))
	if frame < edge {
		return float64(frame) / float64(edge)
	}
	if frame >= frames-edge {
		return float64(frames-frame-1) / float64(edge)
	}
	return 1
}

func putStereo(pcm []float32, frame int, sample float64) {
	pcm[frame*2] = float32(sample)
	pcm[frame*2+1] = float32(sample)
}
