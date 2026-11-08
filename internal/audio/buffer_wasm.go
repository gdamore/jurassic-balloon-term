// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

//go:build js && wasm

package audio

import "time"

const (
	// WebAudio requests PCM from Go on the browser main thread. Leave it twice
	// the native callback window so a 30 Hz render does not starve the worklet.
	// This is exactly 4096 frames at 48 kHz after Oto converts it to bytes.
	otoOutputBufferSize = 85_334 * time.Microsecond
	otoPlayerBufferSize = 8_192
)
