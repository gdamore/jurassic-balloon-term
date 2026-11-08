// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v3"
)

const (
	updateStep      = 33 * time.Millisecond
	burnerPulse     = 300 * time.Millisecond
	stageOverlay    = 1500 * time.Millisecond
	audioOverlay    = 1100 * time.Millisecond
	preferredWidth  = 120
	preferredHeight = 40
	volumeSteps     = 10
)

// App coordinates terminal events, fixed simulation updates, and rendering.
type App struct {
	screen        tcell.Screen
	model         *game.Model
	sound         Sound
	started       bool
	credits       bool
	music         bool
	effects       bool
	volume        int
	keyRelease    bool
	mouseBurner   bool
	burnerKeys    map[burnerKey]struct{}
	burnerUntil   time.Time
	overlayTitle  string
	overlayDetail string
	overlayUntil  time.Time
	overlayAudio  bool
	titleTime     time.Duration
}

type burnerKey struct {
	key      tcell.Key
	physical tcell.Key
	str      string
}

func New(screen tcell.Screen, sounds ...Sound) *App {
	app := &App{screen: screen, music: true, effects: true, volume: volumeSteps}
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
	a.screen.EnableMouse(tcell.MouseButtonEvents)
	width, height := a.screen.Size()
	if width != preferredWidth || height != preferredHeight {
		a.screen.SetSize(preferredWidth, preferredHeight)
		width, height = a.screen.Size()
	}
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
				a.titleTime += updateStep
				a.draw()
				continue
			}
			if !a.keyRelease && !a.mouseBurner && a.model.BurnerOn() && time.Now().After(a.burnerUntil) {
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
	case *tcell.EventMouse:
		a.handleMouseEvent(event)
	case *tcell.EventKey:
		if event.Pressed() && isQuitEvent(event) {
			if a.started {
				a.returnToTitle()
				return false
			}
			if a.credits {
				a.credits = false
				return false
			}
			return true
		}
		if !a.started && a.credits {
			if event.Pressed() && event.Key() == tcell.KeyRune && (event.Str() == "c" || event.Str() == "C") {
				a.credits = false
			}
			return false
		}
		if stageIndex, selected := selectedStage(event); selected {
			a.resetStage(stageIndex)
			a.started = true
			return false
		}
		if a.handleAudioEvent(event) {
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
		case tcell.KeyRune:
			switch event.Str() {
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

func (a *App) returnToTitle() {
	if a.sound != nil {
		a.sound.SetPaused(false)
	}
	a.resetStage(0)
	a.started = false
	a.credits = false
	a.titleTime = 0
}

func (a *App) handleStartEvent(event *tcell.EventKey) {
	if !event.Pressed() {
		return
	}
	if event.Key() == tcell.KeyRune && (event.Str() == "c" || event.Str() == "C") {
		a.credits = true
		return
	}
	if event.Key() == tcell.KeyEnter || isBurnerEvent(event) {
		a.started = true
	}
}

func (a *App) handleMouseEvent(event *tcell.EventMouse) {
	if !a.started && a.credits {
		return
	}
	buttons := event.Buttons()
	if buttons != tcell.ButtonNone && buttons&tcell.ButtonPrimary == 0 {
		return
	}
	pressed := buttons&tcell.ButtonPrimary != 0
	if pressed == a.mouseBurner {
		return
	}
	a.mouseBurner = pressed

	if pressed {
		if !a.started {
			a.started = true
			return
		}
		if a.model.Completed() {
			next := a.model.StageIndex() + 1
			if next >= game.StageCount() {
				next = 0
			}
			a.resetStage(next)
			return
		}
		if a.model.Popped() {
			a.resetStage(a.model.StageIndex())
			return
		}
		a.setBurner(true)
		return
	}

	if !a.started {
		return
	}
	if a.keyRelease {
		a.setBurner(len(a.burnerKeys) != 0)
		return
	}
	a.setBurner(time.Now().Before(a.burnerUntil))
}

func (a *App) handleAudioEvent(event *tcell.EventKey) bool {
	if !event.Pressed() || event.Key() != tcell.KeyRune {
		return false
	}
	switch event.Str() {
	case "m", "M":
		a.music = !a.music
		if a.sound != nil {
			a.sound.SetMusicEnabled(a.music)
		}
		a.showAudioOverlay("MUSIC " + enabledLabel(a.music))
	case "s", "S":
		a.effects = !a.effects
		if a.sound != nil {
			a.sound.SetEffectsEnabled(a.effects)
		}
		a.showAudioOverlay("SOUND " + enabledLabel(a.effects))
	case "+", "=":
		a.setVolume(a.volume + 1)
	case "-", "_":
		a.setVolume(a.volume - 1)
	default:
		return false
	}
	return true
}

func (a *App) setVolume(steps int) {
	a.volume = min(max(steps, 0), volumeSteps)
	if a.sound != nil {
		a.sound.SetMasterVolume(float64(a.volume) / volumeSteps)
	}
	a.showAudioOverlay("VOLUME " + volumeLabel(a.volume))
}

func (a *App) showAudioOverlay(title string) {
	if !a.started {
		return
	}
	a.overlayTitle = title
	a.overlayDetail = "M music  •  S sound  •  −/+ volume"
	a.overlayUntil = time.Now().Add(audioOverlay)
	a.overlayAudio = true
}

func volumeLabel(steps int) string {
	return fmt.Sprintf("%d%%", min(max(steps, 0), volumeSteps)*100/volumeSteps)
}

func (a *App) resetStage(stageIndex int) {
	width, height := a.model.Width, a.model.Height
	if a.screen != nil {
		width, height = a.screen.Size()
	}
	a.model = game.NewStageModel(width, height, stageIndex)
	a.burnerKeys = make(map[burnerKey]struct{})
	a.mouseBurner = false
	a.burnerUntil = time.Time{}
	a.overlayTitle = ""
	a.overlayDetail = ""
	a.overlayUntil = time.Time{}
	a.overlayAudio = false
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
	a.setBurner(len(a.burnerKeys) != 0 || a.mouseBurner)
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
	for _, event := range a.model.Events() {
		a.handleGameEvent(event)
	}
}

func (a *App) handleGameEvent(event game.Event) {
	if event == game.EventStageCompleted && a.model.StageIndex()+1 < game.StageCount() {
		next := game.StageAt(a.model.StageIndex() + 1)
		a.overlayTitle = "STAGE CLEAR"
		a.overlayDetail = next.Name
		a.overlayUntil = time.Now().Add(stageOverlay)
		a.overlayAudio = false
	}
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
	case game.EventGroundLightningStruck:
		a.sound.GroundLightningStruck()
	case game.EventStageCompleted:
		if a.model.StageIndex()+1 >= game.StageCount() {
			a.sound.SetBurner(false)
			a.sound.StageCompleted()
		}
	}
}
