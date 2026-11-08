// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

//go:build !js || !wasm

package audio

import "time"

const (
	// Native audio callbacks comfortably sustain this small buffer, keeping
	// gameplay cues responsive without sacrificing clean playback.
	otoOutputBufferSize = 42_667 * time.Microsecond
	otoPlayerBufferSize = 4_096
)
