// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package audio

import (
	"errors"
	"testing"
)

func TestCatalogValidation(t *testing.T) {
	format := Format{SampleRate: 48_000, Channels: 2}
	tests := []struct {
		name    string
		catalog Catalog
	}{
		{
			name: "missing semantic sound",
			catalog: Catalog{
				Format: format,
				Sounds: map[SoundID]Sound{},
				Cues:   map[Cue]SoundID{CueBalloonPop: "missing"},
			},
		},
		{
			name: "partial frame",
			catalog: Catalog{
				Format: format,
				Sounds: map[SoundID]Sound{
					"bad": {Bus: BusEffects, Variants: []Clip{{PCM: []float32{1}}}},
				},
			},
		},
		{
			name: "invalid loop",
			catalog: Catalog{
				Format: format,
				Sounds: map[SoundID]Sound{
					"bad": {Bus: BusEffects, Variants: []Clip{{PCM: stereo(1), LoopStart: 1, LoopEnd: 1}}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.catalog.Validate(); err == nil {
				t.Fatal("Validate() succeeded, want an error")
			}
		})
	}

	valid := testCatalog()
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() returned %v", err)
	}
	invalidFormat := valid
	invalidFormat.Format = Format{}
	if err := invalidFormat.Validate(); !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("Validate() error = %v, want ErrInvalidFormat", err)
	}
}
