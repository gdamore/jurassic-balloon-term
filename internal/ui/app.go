// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v3"
)

const updateStep = 33 * time.Millisecond
const burnerPulse = 300 * time.Millisecond

// App coordinates terminal events, fixed simulation updates, and rendering.
type App struct {
	screen      tcell.Screen
	model       *game.Model
	sound       Sound
	started     bool
	music       bool
	effects     bool
	keyRelease  bool
	burnerKeys  map[burnerKey]struct{}
	burnerUntil time.Time
}

type burnerKey struct {
	key      tcell.Key
	physical tcell.Key
	str      string
}

func New(screen tcell.Screen, sounds ...Sound) *App {
	app := &App{screen: screen, music: true, effects: true}
	if len(sounds) != 0 {
		app.sound = sounds[0]
	}
	return app
}

func (a *App) Run() error {
	if err := a.screen.Init(); err != nil {
		return err
	}
	defer a.screen.Fini()

	a.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	a.screen.HideCursor()
	width, height := a.screen.Size()
	a.model = game.NewStageModel(width, height, 0)
	if a.sound != nil {
		a.sound.Start()
	}
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
			if !a.started {
				continue
			}
			if !a.keyRelease && a.model.BurnerOn() && time.Now().After(a.burnerUntil) {
				a.setBurner(false)
			}
			a.model.Update(updateStep)
			a.handleGameEvents()
			a.continueCompletedStage()
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
		if stageIndex, selected := selectedStage(event); selected {
			a.resetStage(stageIndex)
			a.started = true
			return false
		}
		if !a.started {
			a.handleStartEvent(event)
			return false
		}
		if a.model.Completed() {
			if event.Pressed() && isRetryEvent(event) {
				next := a.model.StageIndex() + 1
				if next >= game.StageCount() {
					next = 0
				}
				a.resetStage(next)
			}
			return false
		}
		if a.model.Popped() {
			if event.Pressed() && isRetryEvent(event) {
				a.resetStage(a.model.StageIndex())
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
				if a.sound != nil {
					a.sound.SetPaused(a.model.Paused)
				}
			}
		}
	}
	return false
}

func (a *App) handleStartEvent(event *tcell.EventKey) {
	if !event.Pressed() {
		return
	}
	if event.Key() == tcell.KeyRune {
		switch event.Str() {
		case "m", "M":
			a.music = !a.music
			if a.sound != nil {
				a.sound.SetMusicEnabled(a.music)
			}
			return
		case "s", "S":
			a.effects = !a.effects
			if a.sound != nil {
				a.sound.SetEffectsEnabled(a.effects)
			}
			return
		}
	}
	if event.Key() == tcell.KeyEnter || isBurnerEvent(event) {
		a.started = true
		if isBurnerEvent(event) {
			a.handleBurnerEvent(event)
		}
	}
}

func (a *App) resetStage(stageIndex int) {
	width, height := a.model.Width, a.model.Height
	if a.screen != nil {
		width, height = a.screen.Size()
	}
	a.model = game.NewStageModel(width, height, stageIndex)
	a.burnerKeys = make(map[burnerKey]struct{})
	a.burnerUntil = time.Time{}
	if a.sound != nil {
		a.sound.SetBurner(false)
		a.sound.Retry()
		a.sound.SetStage(stageIndex)
	}
}

func (a *App) continueCompletedStage() {
	if !a.model.Completed() || a.model.StageIndex()+1 >= game.StageCount() {
		return
	}
	a.model = a.model.ContinueToStage(a.model.StageIndex() + 1)
	if a.sound != nil {
		a.sound.SetStage(a.model.StageIndex())
	}
}

func selectedStage(event *tcell.EventKey) (int, bool) {
	if !event.Pressed() || event.Key() != tcell.KeyRune {
		return 0, false
	}
	runes := []rune(event.Str())
	if len(runes) != 1 || runes[0] < '1' || runes[0] > '9' {
		return 0, false
	}
	index := int(runes[0] - '1')
	return index, index < game.StageCount()
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
			a.setBurner(true)
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
	a.setBurner(len(a.burnerKeys) != 0)
}

func (a *App) setBurner(enabled bool) {
	if a.model.BurnerOn() == enabled {
		return
	}
	a.model.SetBurner(enabled)
	if a.sound != nil {
		a.sound.SetBurner(enabled)
	}
}

func (a *App) handleGameEvents() {
	if a.sound == nil {
		return
	}
	for _, event := range a.model.Events() {
		a.handleGameEvent(event)
	}
}

func (a *App) handleGameEvent(event game.Event) {
	if a.sound == nil {
		return
	}
	switch event {
	case game.EventBalloonPopped:
		a.sound.SetBurner(false)
		a.sound.BalloonPopped()
	case game.EventPterosaurLunged:
		a.sound.PterosaurLunged()
	case game.EventLightningStruck:
		a.sound.LightningStruck()
	case game.EventStageCompleted:
		if a.model.StageIndex()+1 >= game.StageCount() {
			a.sound.SetBurner(false)
			a.sound.StageCompleted()
		}
	}
}
