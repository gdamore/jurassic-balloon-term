// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

func TestTerrainLayersProvideParallaxAndFatalGround(t *testing.T) {
	m := NewModel(100, 30)
	if len(m.terrain) != 2 {
		t.Fatalf("terrain layer count = %d, want 2", len(m.terrain))
	}
	far, ground := m.terrain[0], m.terrain[1]
	if far.speed <= ground.speed {
		t.Fatalf("terrain speeds are not ordered far-to-near: %v, %v", far.speed, ground.speed)
	}
	if !(far.layer < ground.layer && ground.layer < m.balloon.Layer()) {
		t.Fatalf("terrain layers are not ordered far-to-near: %d, %d, balloon %d", far.layer, ground.layer, m.balloon.Layer())
	}
	for _, chunk := range far.chunks {
		if chunk.Collidable() {
			t.Fatalf("background terrain %q is collidable", chunk.Name())
		}
	}
	for _, chunk := range ground.chunks {
		if !chunk.Collidable() {
			t.Fatalf("foreground terrain %q is not collidable", chunk.Name())
		}
	}
}

func TestTerrainChunksStayConnectedAndReplenish(t *testing.T) {
	m := NewModel(80, 24)
	ground := m.terrain[1]
	m.Update(250 * time.Millisecond)
	for index := 1; index < len(ground.chunks); index++ {
		left, _ := ground.chunks[index-1].Position()
		right, _ := ground.chunks[index].Position()
		if right-left != terrainChunkWidth {
			t.Fatalf("terrain gap between chunks %d and %d = %v, want %d", index-1, index, right-left, terrainChunkWidth)
		}
	}

	departed := ground.chunks[0]
	departed.SetPosition(-terrainChunkWidth-1, float64(m.Height-ground.maxHeight))
	ground.replenish(m.scene, m.Width, m.Height)
	for _, chunk := range ground.chunks {
		if chunk == departed {
			t.Fatal("departed terrain chunk remains active")
		}
	}
	lastX, _ := ground.chunks[len(ground.chunks)-1].Position()
	if lastX+terrainChunkWidth < float64(m.Width+terrainChunkWidth) {
		t.Fatal("terrain replenishment did not cover the approach to the screen")
	}
}

func TestFarTerrainIncludesCrateredVolcano(t *testing.T) {
	m := NewModel(80, 24)
	frame, _ := m.terrain[0].chunks[1].ActiveFrame()
	center := columnHeight(frame, terrainChunkWidth/2)
	leftPeak := columnHeight(frame, terrainChunkWidth/2-2)
	rightPeak := columnHeight(frame, terrainChunkWidth/2+2)
	if center >= leftPeak || center >= rightPeak {
		t.Fatalf("volcano lacks crater: heights left=%d center=%d right=%d", leftPeak, center, rightPeak)
	}
}

func TestTerrainUsesQuarterBlocksForSlopes(t *testing.T) {
	if got := terrainSurface(3, 2); got != '▟' {
		t.Fatalf("rising terrain glyph = %q, want ▟", got)
	}
	if got := terrainSurface(2, 3); got != '▙' {
		t.Fatalf("falling terrain glyph = %q, want ▙", got)
	}
	if got := terrainSurface(2, 2); got != '▄' {
		t.Fatalf("level terrain glyph = %q, want ▄", got)
	}
}

func TestTerrainContinuesVarying(t *testing.T) {
	layer := newTerrainLayer("ground", 37, -3.2, groundLayer, 2, 5, StyleGroundSurface, StyleGroundFill, true)
	longest, run := 0, 0
	previous := -1
	for range 12 {
		frame := layer.nextFrame()
		for x := range frame.Width {
			height := columnHeight(frame, x)
			if height == previous {
				run++
			} else {
				run = 1
				previous = height
			}
			longest = max(longest, run)
		}
	}
	if longest > terrainChunkWidth {
		t.Fatalf("terrain remained flat for %d columns", longest)
	}
}

func TestForegroundTerrainPopsBalloon(t *testing.T) {
	m := NewModel(80, 30)
	makeAirborne(m)
	ground := m.terrain[1].chunks[0]
	frame, _ := ground.ActiveFrame()
	terrainX, terrainY := firstOpaqueCell(frame)
	balloonX, balloonY := m.balloon.Position()
	ground.SetPosition(balloonX+4-float64(terrainX), balloonY+1-float64(terrainY))
	ground.SetVelocity(0, 0)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("foreground terrain did not pop balloon")
	}
}

func columnHeight(frame sprite.Frame, x int) int {
	for y := 0; y < frame.Height; y++ {
		if frame.CellAt(x, y).Rune != 0 {
			return frame.Height - y
		}
	}
	return 0
}

func firstOpaqueCell(frame sprite.Frame) (int, int) {
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			if frame.CellAt(x, y).Rune != 0 {
				return x, y
			}
		}
	}
	return 0, 0
}
