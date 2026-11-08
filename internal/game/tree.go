// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	treeLayer  = 9
	treeHeight = 10
)

type treeTrack struct {
	sprite   *sprite.Sprite
	distance float64
}

func newTree(index int, distance float64) treeTrack {
	tree := sprite.New(fmt.Sprintf("tree-%d", index+1))
	left := treeFrame(
		"  ▄▄   ",
		" ▟██▙  ",
		"▟████▙ ",
		"██████▌",
		"▀████▀ ",
		"  ▐▌   ",
		"  ▐▌   ",
		"  ▐▌   ",
		"  ▐▌   ",
		"       ",
	)
	left.Duration, left.Next = 650*time.Millisecond, "right"
	right := treeFrame(
		"   ▄▄  ",
		"  ▟██▙ ",
		" ▟████▙",
		"▐██████",
		" ▀████▀",
		"  ▐▌   ",
		"  ▐▌   ",
		"  ▐▌   ",
		"  ▐▌   ",
		"       ",
	)
	right.Duration, right.Next = 650*time.Millisecond, "left"
	for name, frame := range map[string]sprite.Frame{"left": left, "right": right} {
		if err := tree.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := tree.SetFrame("left"); err != nil {
		panic(err)
	}
	tree.SetLayer(treeLayer)
	tree.SetVelocity(-worldScrollSpeed, 0)
	return treeTrack{sprite: tree, distance: distance}
}

func treeFrame(rows ...string) sprite.Frame {
	frame := styledFrame(StyleTreeCanopy, rows...)
	for y := 5; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune != 0 {
				frame.Cells[index].Style = StyleTreeTrunk
			}
		}
	}
	for _, point := range [][2]int{{2, 1}, {3, 1}, {1, 2}, {2, 2}, {1, 3}} {
		x, y := point[0], point[1]
		index := y*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = StyleTreeHighlight
		}
	}
	return frame
}
