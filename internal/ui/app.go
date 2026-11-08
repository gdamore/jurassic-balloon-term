// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v3"
)

const updateStep = 33 * time.Millisecond
const burnerPulse = 200 * time.Millisecond

// App coordinates terminal events, fixed simulation updates, and rendering.
type App struct {
	screen      tcell.Screen
	model       *game.Model
	keyRelease  bool
	burnerKeys  map[burnerKey]struct{}
	burnerUntil time.Time
}

type burnerKey struct {
	key      tcell.Key
	physical tcell.Key
	str      string
}

func New(screen tcell.Screen) *App {
	return &App{screen: screen}
}

func (a *App) Run() error {
	if err := a.screen.Init(); err != nil {
		return err
	}
	defer a.screen.Fini()

	a.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	a.screen.HideCursor()
	width, height := a.screen.Size()
	a.model = game.NewModel(width, height)
	protocol := a.screen.KeyboardProtocol()
	a.keyRelease = supportsKeyRelease(protocol)
	a.burnerKeys = make(map[burnerKey]struct{})

	events := a.screen.EventQ()

	ticker := time.NewTicker(updateStep)
	defer ticker.Stop()
	a.draw()

	for {
		select {
		case event, ok := <-events:
			if !ok {
				return nil
			}
			quit := a.handleEvent(event)
			if quit {
				return nil
			}
			a.draw()
		case <-ticker.C:
			if !a.keyRelease && a.model.BurnerOn() && time.Now().After(a.burnerUntil) {
				a.model.SetBurner(false)
			}
			a.model.Update(updateStep)
			a.draw()
		}
	}
}

func supportsKeyRelease(protocol tcell.KeyProtocol) bool {
	return protocol == tcell.KittyKeyboard || protocol == tcell.Win32Keyboard
}

func (a *App) handleEvent(event tcell.Event) bool {
	switch event := event.(type) {
	case *tcell.EventResize:
		width, height := event.Size()
		a.model.Resize(width, height)
		a.screen.Sync()
	case *tcell.EventKey:
		if event.Pressed() && isQuitEvent(event) {
			return true
		}
		if a.model.Popped() {
			if event.Pressed() && isRetryEvent(event) {
				width, height := a.screen.Size()
				a.model = game.NewModel(width, height)
				a.burnerKeys = make(map[burnerKey]struct{})
				a.burnerUntil = time.Time{}
			}
			return false
		}
		if isBurnerEvent(event) {
			a.handleBurnerEvent(event)
			return false
		}
		if !event.Pressed() {
			return false
		}
		switch event.Key() {
		case tcell.KeyEscape:
			return true
		case tcell.KeyRune:
			switch event.Str() {
			case "q", "Q":
				return true
			case "p", "P":
				a.model.TogglePause()
			}
		}
	}
	return false
}

func isQuitEvent(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyEscape {
		return true
	}
	return event.Key() == tcell.KeyRune && (event.Str() == "q" || event.Str() == "Q")
}

func isRetryEvent(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyEnter {
		return true
	}
	if event.Key() != tcell.KeyRune {
		return false
	}
	return event.Str() == " " || event.Str() == "r" || event.Str() == "R"
}

func isBurnerEvent(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyUp {
		return true
	}
	if event.Key() != tcell.KeyRune {
		return false
	}
	return event.Str() == " " || event.Str() == "w" || event.Str() == "W"
}

func (a *App) handleBurnerEvent(event *tcell.EventKey) {
	if !a.keyRelease {
		if event.Pressed() {
			a.model.SetBurner(true)
			a.burnerUntil = time.Now().Add(burnerPulse)
		}
		return
	}

	key := burnerKey{key: event.Key(), physical: event.Physical(), str: event.Str()}
	if event.Pressed() {
		a.burnerKeys[key] = struct{}{}
	} else {
		delete(a.burnerKeys, key)
	}
	a.model.SetBurner(len(a.burnerKeys) != 0)
}
