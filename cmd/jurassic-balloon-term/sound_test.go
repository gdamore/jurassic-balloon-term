// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package main

import (
	"testing"

	"github.com/gdamore/jurassic-balloon-term/internal/audio"
)

func TestMusicForStage(t *testing.T) {
	tests := []struct {
		stage int
		want  audio.Music
	}{
		{stage: 0, want: audio.MusicClearSkies},
		{stage: 1, want: audio.MusicClearSkies},
		{stage: 2, want: audio.MusicThreat},
		{stage: 3, want: audio.MusicClearSkies},
	}
	for _, test := range tests {
		if got := musicForStage(test.stage); got != test.want {
			t.Errorf("musicForStage(%d) = %v, want %v", test.stage, got, test.want)
		}
	}
}
