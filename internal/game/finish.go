// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	finishLayer  = 9
	finishHeight = 8
)

func newFinishFlag() *sprite.Sprite {
	flag := sprite.New("finish-flag")
	left := finishFrame(
		"  ▟▀▀▙",
		"  █★ █",
		"  ▜▄▄▛",
		"  │   ",
		"  │   ",
		"  │   ",
		" ╱│╲  ",
		"      ",
	)
	left.Duration, left.Next = 220*time.Millisecond, "right"
	right := finishFrame(
		"  ▟▀▙ ",
		"  █ ★█",
		"  ▜▄▛ ",
		"  │   ",
		"  │   ",
		"  │   ",
		" ╱│╲  ",
		"      ",
	)
	right.Duration, right.Next = 220*time.Millisecond, "left"
	for name, frame := range map[string]sprite.Frame{"left": left, "right": right} {
		if err := flag.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := flag.SetFrame("left"); err != nil {
		panic(err)
	}
	flag.SetLayer(finishLayer)
	flag.SetCollidable(false)
	flag.SetVelocity(-worldScrollSpeed, 0)
	return flag
}

func finishFrame(rows ...string) sprite.Frame {
	frame := sprite.FrameFromStrings(rows...)
	for index := range frame.Cells {
		if frame.Cells[index].Rune == 0 {
			continue
		}
		switch frame.Cells[index].Rune {
		case '│', '╱', '╲':
			frame.Cells[index].Style = StyleFinishPole
		default:
			frame.Cells[index].Style = StyleFinishFlag
		}
	}
	return frame
}
