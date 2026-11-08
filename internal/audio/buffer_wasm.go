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

	// Oto reads one complete output callback from the player at a time and
	// zero-fills any samples the player cannot supply. Keep this at least as
	// large as the 4096-frame, stereo float32 callback above; a smaller buffer
	// produces periodic audio followed by silence (the "helicopter" effect).
	otoPlayerBufferSize = 32_768
)
