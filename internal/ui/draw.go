// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
	"github.com/gdamore/tcell/v2"
)

func (a *App) draw() {
	a.screen.Clear()
	width, height := a.screen.Size()

	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	infoStyle := tcell.StyleDefault.Foreground(tcell.ColorSilver)
	balloonStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)

	putCentered(a.screen, 0, width, "JURASSIC BALLOON TERM", titleStyle)
	status := "arrows/WASD move  •  space pauses  •  q/esc quits"
	if a.model.Paused {
		status = "PAUSED  •  space resumes  •  q/esc quits"
	}
	putCentered(a.screen, height-1, width, status, infoStyle)

	drawSprites(a.screen, a.model.Sprites(), map[int]tcell.Style{0: balloonStyle})

	if width >= 34 && height >= 12 {
		elapsed := fmt.Sprintf("prototype runtime %0.1fs", a.model.Elapsed.Seconds())
		putCentered(a.screen, height-3, width, elapsed, infoStyle)
	}
	a.screen.Show()
}

func drawSprites(screen tcell.Screen, sprites []*sprite.Sprite, palette map[int]tcell.Style) {
	for _, item := range sprites {
		frame, visible := item.ActiveFrame()
		bounds, bounded := item.Bounds()
		if !visible || !bounded {
			continue
		}
		for y := 0; y < frame.Height; y++ {
			for x := 0; x < frame.Width; x++ {
				cell := frame.CellAt(x, y)
				if cell.Rune == 0 {
					continue
				}
				style, ok := palette[cell.Style]
				if !ok {
					style = tcell.StyleDefault
				}
				screen.SetContent(bounds.Min.X+x, bounds.Min.Y+y, cell.Rune, nil, style)
			}
		}
	}
}

func putCentered(screen tcell.Screen, y, width int, text string, style tcell.Style) {
	if y < 0 {
		return
	}
	x := max((width-len([]rune(text)))/2, 0)
	putString(screen, x, y, text, style)
}

func putString(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	width, height := screen.Size()
	if y < 0 || y >= height || x >= width {
		return
	}
	if x < 0 {
		runes := []rune(text)
		trim := min(-x, len(runes))
		text = string(runes[trim:])
		x = 0
	}
	if x < width {
		screen.PutStrStyled(x, y, strings.TrimRight(text, "\n"), style)
	}
}
