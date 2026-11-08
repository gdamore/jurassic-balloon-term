// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

//go:build js && wasm

package main

import "syscall/js"

// browserReady enables the launch control before Oto waits for the browser's
// first user gesture. That gesture unlocks WebAudio and lets startup continue.
func browserReady() {
	event := js.Global().Get("Event").New("jurassic-balloon-ready")
	js.Global().Call("dispatchEvent", event)
}
