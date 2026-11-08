// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	fernLayer  = 7
	fernHeight = 5
)

type fernTrack struct {
	sprite   *sprite.Sprite
	distance float64
}

func newFern(index int, distance float64) fernTrack {
	fern := sprite.New(fmt.Sprintf("fern-%d", index+1))
	left := fernFrame(
		"   ▄   ",
		"▗▄▙█▟▄▖",
		" ▝▟█▙▘ ",
		"  ▝▌▘  ",
		"       ",
	)
	left.Duration, left.Next = 750*time.Millisecond, "right"
	right := fernFrame(
		"  ▄    ",
		" ▗▙█▟▄▖",
		"▝▟██▙▘ ",
		"  ▝▌▘  ",
		"       ",
	)
	right.Duration, right.Next = 750*time.Millisecond, "left"
	for name, frame := range map[string]sprite.Frame{"left": left, "right": right} {
		if err := fern.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := fern.SetFrame("left"); err != nil {
		panic(err)
	}
	fern.SetLayer(fernLayer)
	fern.SetCollidable(false)
	fern.SetVelocity(-worldScrollSpeed, 0)
	return fernTrack{sprite: fern, distance: distance}
}

func fernFrame(rows ...string) sprite.Frame {
	frame := styledFrame(StyleFernLeaf, rows...)
	for _, point := range [][2]int{{3, 0}, {1, 1}, {2, 1}, {2, 2}} {
		x, y := point[0], point[1]
		index := y*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = StyleFernHighlight
		}
	}
	for x := 0; x < frame.Width; x++ {
		index := 3*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = StyleFernStem
		}
	}
	return frame
}
