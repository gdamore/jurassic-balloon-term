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

func TestPterosaurCallIsShortCrowLikeCue(t *testing.T) {
	sound := Catalog().Sounds[soundPterosaur]
	if sound.Gain < 0.7 || sound.MaxVoices != 1 || sound.Cooldown < 400*time.Millisecond {
		t.Fatalf("pterosaur mix = gain:%v voices:%d cooldown:%v", sound.Gain, sound.MaxVoices, sound.Cooldown)
	}
	if len(sound.Variants) != 2 {
		t.Fatalf("pterosaur variants = %d, want 2", len(sound.Variants))
	}
	for index, variant := range sound.Variants {
		duration := time.Duration(len(variant.PCM)/audio.DefaultChannels) * time.Second / audio.DefaultSampleRate
		if duration < 250*time.Millisecond || duration > 400*time.Millisecond {
			t.Errorf("variant %d duration = %v, want a short squawk", index, duration)
		}
	}
}

func TestCrowSquawkIsDeterministicAndDecays(t *testing.T) {
	first := crowSquawk(760, 330, 320*time.Millisecond, 0.42, 29)
	second := crowSquawk(760, 330, 320*time.Millisecond, 0.42, 29)
	if len(first) == 0 || len(first) != len(second) {
		t.Fatalf("squawk lengths = %d and %d", len(first), len(second))
	}
	for index := range first {
		if first[index] != second[index] {
			t.Fatalf("squawk differs at sample %d", index)
		}
	}
	attackFrames := durationFrames(80 * time.Millisecond)
	tailFrames := durationFrames(80 * time.Millisecond)
	if rms(first[:attackFrames*audio.DefaultChannels]) <= rms(first[len(first)-tailFrames*audio.DefaultChannels:]) {
		t.Fatal("squawk tail does not decay below its attack")
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
