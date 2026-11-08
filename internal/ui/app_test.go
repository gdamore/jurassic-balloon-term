// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"testing"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v3"
)

func TestKeyReleaseCapability(t *testing.T) {
	for _, test := range []struct {
		protocol tcell.KeyProtocol
		want     bool
	}{
		{protocol: tcell.KittyKeyboard, want: true},
		{protocol: tcell.Win32Keyboard, want: true},
		{protocol: tcell.XTermKeyboard, want: false},
		{protocol: tcell.LegacyKeyboard, want: false},
	} {
		if got := supportsKeyRelease(test.protocol); got != test.want {
			t.Errorf("supportsKeyRelease(%v) = %v, want %v", test.protocol, got, test.want)
		}
	}
}

func TestSkyColorVariesWithAltitude(t *testing.T) {
	top := skyColor(0, 24)
	bottom := skyColor(23, 24)
	if top == bottom {
		t.Fatal("sky background does not vary with altitude")
	}
	if got := skyColor(-10, 24); got != top {
		t.Fatalf("sky color above screen = %v, want top color %v", got, top)
	}
	if got := skyColor(100, 24); got != bottom {
		t.Fatalf("sky color below screen = %v, want bottom color %v", got, bottom)
	}
}

func TestRetryAndQuitKeys(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, " ", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, "r", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEnter, "", tcell.ModNone),
	} {
		if !isRetryEvent(event) {
			t.Errorf("%v is not recognized as retry", event)
		}
	}
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, "q", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone),
	} {
		if !isQuitEvent(event) {
			t.Errorf("%v is not recognized as quit", event)
		}
	}
}

func TestAdvancedBurnerInputTracksPressAndRelease(t *testing.T) {
	a := &App{
		model:      game.NewModel(80, 24),
		keyRelease: true,
		burnerKeys: make(map[burnerKey]struct{}),
	}
	press := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	release := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1)

	a.handleEvent(press)
	if !a.model.BurnerOn() {
		t.Fatal("burner is off after key press")
	}
	a.handleEvent(release)
	if a.model.BurnerOn() {
		t.Fatal("burner is on after key release")
	}
}

func TestAdvancedBurnerInputTracksMultipleKeys(t *testing.T) {
	a := &App{
		model:      game.NewModel(80, 24),
		keyRelease: true,
		burnerKeys: make(map[burnerKey]struct{}),
	}
	spaceDown := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	upDown := tcell.NewEventKeyEx(tcell.KeyUp, "", tcell.ModNone, true, tcell.KeyUp, 1)
	spaceUp := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1)
	upUp := tcell.NewEventKeyEx(tcell.KeyUp, "", tcell.ModNone, false, tcell.KeyUp, 1)

	for _, event := range []*tcell.EventKey{spaceDown, upDown, spaceUp} {
		a.handleEvent(event)
	}
	if !a.model.BurnerOn() {
		t.Fatal("releasing one of two held burner keys turned burner off")
	}
	a.handleEvent(upUp)
	if a.model.BurnerOn() {
		t.Fatal("burner is on after all burner keys were released")
	}
}

func TestLegacyBurnerInputStartsPulse(t *testing.T) {
	a := &App{model: game.NewModel(80, 24)}
	press := tcell.NewEventKey(tcell.KeyRune, "w", tcell.ModNone)

	a.handleEvent(press)
	if !a.model.BurnerOn() {
		t.Fatal("burner is off after legacy key press")
	}
	if a.burnerUntil.IsZero() {
		t.Fatal("legacy key press did not set pulse deadline")
	}
}
