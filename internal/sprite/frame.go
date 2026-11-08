// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

// Package sprite implements terminal-independent animated sprites and
// cell-accurate collision detection. Its design is derived from the sprite
// engine in github.com/gdamore/proxima5.
package sprite

import "time"

// Cell is one drawable sprite cell. A zero Rune is transparent. Style is a
// renderer-defined palette index, keeping this package independent of tcell.
type Cell struct {
	Rune        rune
	Style       int
	NoCollision bool
}

// Frame is one image in a sprite animation.
type Frame struct {
	Width    int
	Height   int
	Cells    []Cell
	Duration time.Duration
	Next     string
}

// FrameFromStrings creates a frame from text rows. Spaces are transparent.
func FrameFromStrings(rows ...string) Frame {
	height := len(rows)
	width := 0
	for _, row := range rows {
		width = max(width, len([]rune(row)))
	}

	frame := Frame{Width: width, Height: height, Cells: make([]Cell, width*height)}
	for y, row := range rows {
		for x, r := range []rune(row) {
			if r != ' ' {
				frame.Cells[y*width+x].Rune = r
			}
		}
	}
	return frame
}

// CellAt returns the cell at local frame coordinates.
func (f Frame) CellAt(x, y int) Cell {
	if x < 0 || y < 0 || x >= f.Width || y >= f.Height {
		return Cell{}
	}
	index := y*f.Width + x
	if index >= len(f.Cells) {
		return Cell{}
	}
	return f.Cells[index]
}
