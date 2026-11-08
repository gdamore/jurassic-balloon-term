// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sound

import (
	"testing"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/audio"
)

func TestCatalogValid(t *testing.T) {
	if err := Catalog().Validate(); err != nil {
		t.Fatalf("Catalog().Validate() returned %v", err)
	}
}

func TestGroundStrikeUsesDedicatedRecordedCue(t *testing.T) {
	catalog := Catalog()
	id := catalog.Cues[audio.CueGroundLightningStrike]
	if id != soundGroundStrike {
		t.Fatalf("ground strike cue = %q, want %q", id, soundGroundStrike)
	}
	sound := catalog.Sounds[id]
	if len(sound.Variants) != 1 || sound.MaxVoices != 1 || sound.Cooldown < time.Second {
		t.Fatalf("ground strike mix = variants:%d voices:%d cooldown:%v", len(sound.Variants), sound.MaxVoices, sound.Cooldown)
	}
	duration := time.Duration(len(sound.Variants[0].PCM)/audio.DefaultChannels) * time.Second / audio.DefaultSampleRate
	if duration < 2*time.Second || duration > 3*time.Second {
		t.Fatalf("ground strike duration = %v, want focused 2–3 second transient", duration)
	}
}

func TestClearMusicIsLongQuietLoop(t *testing.T) {
	music := Catalog().Sounds[soundClear]
	if !music.Loop || music.Gain > 0.45 || len(music.Variants) != 1 {
		t.Fatalf("clear music = loop:%t gain:%v variants:%d", music.Loop, music.Gain, len(music.Variants))
	}
	duration := time.Duration(len(music.Variants[0].PCM)/audio.DefaultChannels) * time.Second / audio.DefaultSampleRate
	if duration < 3*time.Minute {
		t.Fatalf("clear music duration = %v, want a long-form loop", duration)
	}
}

func TestPterosaurCallUsesSpacedRecordedVariants(t *testing.T) {
	sound := Catalog().Sounds[soundPterosaur]
	if sound.Gain > 0.7 || sound.MaxVoices != 1 || sound.Cooldown < 2*time.Second {
		t.Fatalf("pterosaur mix = gain:%v voices:%d cooldown:%v", sound.Gain, sound.MaxVoices, sound.Cooldown)
	}
	if len(sound.Variants) != 2 {
		t.Fatalf("pterosaur variants = %d, want 2", len(sound.Variants))
	}
	for index, variant := range sound.Variants {
		duration := time.Duration(len(variant.PCM)/audio.DefaultChannels) * time.Second / audio.DefaultSampleRate
		if duration < 1700*time.Millisecond || duration > 2200*time.Millisecond {
			t.Errorf("variant %d duration = %v, want an isolated two-second call", index, duration)
		}
	}
}

func TestSauropodCallUsesDedicatedDeepBellow(t *testing.T) {
	catalog := Catalog()
	if got := catalog.Cues[audio.CueDistantSauropod]; got != soundSauropod {
		t.Fatalf("sauropod cue = %q, want %q", got, soundSauropod)
	}
	call := catalog.Sounds[soundSauropod]
	if call.Bus != audio.BusEffects || call.Gain < 0.4 || call.MaxVoices != 1 || call.Cooldown < 3*time.Second {
		t.Fatalf("sauropod mix = bus:%v gain:%v voices:%d cooldown:%v", call.Bus, call.Gain, call.MaxVoices, call.Cooldown)
	}
	if len(call.Variants) != 1 {
		t.Fatalf("sauropod variants = %d, want 1", len(call.Variants))
	}
	duration := time.Duration(len(call.Variants[0].PCM)/audio.DefaultChannels) * time.Second / audio.DefaultSampleRate
	if duration < 1500*time.Millisecond || duration > 2*time.Second {
		t.Fatalf("sauropod duration = %v, want a focused 1.5–2 second bellow", duration)
	}
}

func TestDayAmbienceUsesJungleBedWithoutCreatureEffects(t *testing.T) {
	catalog := Catalog()
	scene := catalog.Ambience[audio.EnvironmentPrehistoricDay]
	if scene.Bed != soundJungle {
		t.Fatalf("day ambience bed = %q, want %q", scene.Bed, soundJungle)
	}
	for _, id := range scene.Random {
		if id == soundDinosaur1 || id == soundDinosaur2 || id == soundSauropod || id == soundPterosaurFar {
			t.Fatalf("day ambience still schedules creature effect %q", id)
		}
	}
	if gain := catalog.Sounds[soundJungle].Gain; gain > 0.20 {
		t.Fatalf("jungle ambience gain = %v, want quiet background bed", gain)
	}
}

func rms(samples []float32) float64 {
	var sum float64
	for _, sample := range samples {
		sum += float64(sample * sample)
	}
	return sum / float64(len(samples))
}
