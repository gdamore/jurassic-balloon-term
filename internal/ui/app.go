// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v2"
)

const updateStep = 33 * time.Millisecond

// App coordinates terminal events, fixed simulation updates, and rendering.
type App struct {
	screen tcell.Screen
	model  *game.Model
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

	events := make(chan tcell.Event, 16)
	stopEvents := make(chan struct{})
	defer close(stopEvents)
	go a.screen.ChannelEvents(events, stopEvents)

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
			a.model.Update(updateStep)
			a.draw()
		}
	}
}

func (a *App) handleEvent(event tcell.Event) bool {
	switch event := event.(type) {
	case *tcell.EventResize:
		width, height := event.Size()
		a.model.Resize(width, height)
		a.screen.Sync()
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEscape:
			return true
		case tcell.KeyUp:
			a.model.Move(0, -1)
		case tcell.KeyDown:
			a.model.Move(0, 1)
		case tcell.KeyLeft:
			a.model.Move(-1, 0)
		case tcell.KeyRight:
			a.model.Move(1, 0)
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				return true
			case 'w', 'W':
				a.model.Move(0, -1)
			case 's', 'S':
				a.model.Move(0, 1)
			case 'a', 'A':
				a.model.Move(-1, 0)
			case 'd', 'D':
				a.model.Move(1, 0)
			case ' ':
				a.model.TogglePause()
			}
		}
	}
	return false
}
