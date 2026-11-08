// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	grazerLayer  = 9
	grazerHeight = 4
)

var grazerFootColumns = [...]int{3, 4, 6, 7}

type grazerTrack struct {
	sprite   *sprite.Sprite
	distance float64
	anchored bool
}

func newGrazer(index int, distance float64) grazerTrack {
	grazer := sprite.New(fmt.Sprintf("three-horn-grazer-%d", index+1))
	frame := grazerFrame(
		"▲▲       ",
		"◀█▄▟███▅▂",
		" ▼█▜████▛",
		"   ▀▀ ▀▀ ",
	)
	if err := grazer.AddFrame("still", frame); err != nil {
		panic(err)
	}
	if err := grazer.SetFrame("still"); err != nil {
		panic(err)
	}
	grazer.SetLayer(grazerLayer)
	grazer.SetVelocity(-worldScrollSpeed, 0)
	return grazerTrack{sprite: grazer, distance: distance}
}

func grazerFrame(rows ...string) sprite.Frame {
	frame := styledFrame(StyleGrazerBody, rows...)
	paintGrazer(&frame, StyleGrazerHorn, [][2]int{{0, 0}, {1, 0}, {0, 1}})
	paintGrazer(&frame, StyleGrazerBeak, [][2]int{{1, 2}})
	paintGrazer(&frame, StyleGrazerHead, [][2]int{{1, 1}, {2, 1}, {2, 2}})
	paintGrazer(&frame, StyleGrazerHighlight, [][2]int{{3, 1}, {3, 2}})
	paintGrazer(&frame, StyleGrazerShadow, [][2]int{{4, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 2}})
	paintGrazer(&frame, StyleGrazerLeg, [][2]int{{3, 3}, {4, 3}, {6, 3}, {7, 3}})
	return frame
}

func paintGrazer(frame *sprite.Frame, style int, points [][2]int) {
	for _, point := range points {
		x, y := point[0], point[1]
		index := y*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = style
		}
	}
}
