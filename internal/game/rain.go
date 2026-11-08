// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const rainFieldLayer = 6

func newRainField(width, height int) *sprite.Sprite {
	rain := sprite.New("storm-rain")
	for index := range 3 {
		name := rainFrameName(index)
		frame := rainFieldFrame(width, height, index)
		frame.Duration = 120 * time.Millisecond
		frame.Next = rainFrameName((index + 1) % 3)
		if err := rain.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := rain.SetFrame(rainFrameName(0)); err != nil {
		panic(err)
	}
	rain.SetLayer(rainFieldLayer)
	rain.SetCollidable(false)
	return rain
}

func rainFrameName(index int) string {
	return [...]string{"rain-1", "rain-2", "rain-3"}[index]
}

func rainFieldFrame(width, height, phase int) sprite.Frame {
	width, height = max(width, 1), max(height, 1)
	frame := sprite.Frame{Width: width, Height: height, Cells: make([]sprite.Cell, width*height)}
	for y := 1; y < height-1; y++ {
		start := (y*7 + phase*4) % 13
		for x := start; x < width; x += 13 {
			glyph := '╲'
			if (x+y+phase)%3 == 0 {
				glyph = '│'
			}
			frame.Cells[y*width+x] = sprite.Cell{Rune: glyph, Style: StyleRain}
		}
	}
	return frame
}
