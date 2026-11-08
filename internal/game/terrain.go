// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const terrainChunkWidth = 32

type terrainLayer struct {
	name        string
	chunks      []*sprite.Sprite
	random      *rand.Rand
	speed       float64
	layer       int
	maxHeight   int
	minHeight   int
	lastHeight  int
	surface     int
	fill        int
	hazard      bool
	chunkID     int
	target      int
	segmentLeft int
	slopePhase  int
}

func newTerrainLayer(name string, seed int64, speed float64, layer, minHeight, maxHeight, surface, fill int, hazard bool) *terrainLayer {
	return &terrainLayer{
		name:       name,
		random:     rand.New(rand.NewSource(seed)),
		speed:      speed,
		layer:      layer,
		minHeight:  minHeight,
		maxHeight:  maxHeight,
		lastHeight: minHeight,
		surface:    surface,
		fill:       fill,
		hazard:     hazard,
	}
}

func (t *terrainLayer) initialize(scene *sprite.Scene, width, height int) {
	for x := 0.0; x < float64(width+terrainChunkWidth); x += terrainChunkWidth {
		t.addChunk(scene, x, height)
	}
}

func (t *terrainLayer) addChunk(scene *sprite.Scene, x float64, screenHeight int) {
	t.chunkID++
	frame := t.nextFrame()
	chunk := sprite.New(fmt.Sprintf("%s-%d", t.name, t.chunkID))
	if err := chunk.AddFrame("terrain", frame); err != nil {
		panic(err)
	}
	if err := chunk.SetFrame("terrain"); err != nil {
		panic(err)
	}
	chunk.SetLayer(t.layer)
	chunk.SetCollidable(t.hazard)
	chunk.SetPosition(x, float64(screenHeight-t.maxHeight))
	chunk.SetVelocity(t.speed, 0)
	t.chunks = append(t.chunks, chunk)
	scene.Add(chunk)
}

func (t *terrainLayer) nextFrame() sprite.Frame {
	heights := make([]int, terrainChunkWidth)
	height := t.lastHeight
	for x := range heights {
		if x != 0 {
			if t.segmentLeft == 0 {
				if t.minHeight == t.maxHeight {
					t.target = height
				} else {
					t.target = t.minHeight + t.random.Intn(t.maxHeight-t.minHeight+1)
				}
				if t.target == height && t.minHeight != t.maxHeight {
					if height < t.maxHeight {
						t.target++
					} else {
						t.target--
					}
				}
				distance := abs(t.target - height)
				t.segmentLeft = distance*3 + 6 + t.random.Intn(8)
				t.slopePhase = 0
			}
			t.segmentLeft--
			if height != t.target {
				t.slopePhase++
				if t.slopePhase == 3 {
					if height < t.target {
						height++
					} else {
						height--
					}
					t.slopePhase = 0
				}
			}
		}
		height = min(max(height, t.minHeight), t.maxHeight)
		heights[x] = height
	}

	// Every fourth far-terrain chunk carries a broad, cratered volcano.
	if t.name == "far-terrain" && t.chunkID%4 == 2 {
		center := terrainChunkWidth / 2
		for x := range heights {
			distance := abs(x - center)
			addition := max(0, 8-(distance+1)/2)
			heights[x] = min(heights[x]+addition, t.maxHeight)
		}
		rim := min(heights[center-2], heights[center+2])
		crater := max(t.minHeight, rim-2)
		heights[center] = crater
		heights[center-1] = min(heights[center-1], crater+1)
		heights[center+1] = min(heights[center+1], crater+1)
	}
	t.lastHeight = heights[len(heights)-1]

	frame := sprite.Frame{
		Width:  terrainChunkWidth,
		Height: t.maxHeight,
		Cells:  make([]sprite.Cell, terrainChunkWidth*t.maxHeight),
	}
	for x, columnHeight := range heights {
		top := t.maxHeight - columnHeight
		surface := '▄'
		if x > 0 {
			surface = terrainSurface(columnHeight, heights[x-1])
		}
		frame.Cells[top*terrainChunkWidth+x] = sprite.Cell{Rune: surface, Style: t.surface}
		for y := top + 1; y < t.maxHeight; y++ {
			frame.Cells[y*terrainChunkWidth+x] = sprite.Cell{Rune: '█', Style: t.fill}
		}
	}
	return frame
}

func terrainSurface(current, previous int) rune {
	switch {
	case current > previous:
		return '▟'
	case current < previous:
		return '▙'
	default:
		return '▄'
	}
}

func (t *terrainLayer) resize(height int) {
	for _, chunk := range t.chunks {
		x, _ := chunk.Position()
		chunk.SetPosition(x, float64(height-t.maxHeight))
	}
}

func (t *terrainLayer) replenish(scene *sprite.Scene, width, height int) {
	for len(t.chunks) != 0 {
		first := t.chunks[0]
		x, _ := first.Position()
		if x+terrainChunkWidth >= 0 {
			break
		}
		scene.Remove(first)
		t.chunks = t.chunks[1:]
	}

	right := 0.0
	if len(t.chunks) != 0 {
		x, _ := t.chunks[len(t.chunks)-1].Position()
		right = x + terrainChunkWidth
	}
	for right < float64(width+terrainChunkWidth) {
		t.addChunk(scene, right, height)
		right += terrainChunkWidth
	}
}

func (t *terrainLayer) contains(candidate *sprite.Sprite) bool {
	if !t.hazard {
		return false
	}
	for _, chunk := range t.chunks {
		if chunk == candidate {
			return true
		}
	}
	return false
}

func (t *terrainLayer) surfaceYAt(screenX int) (int, bool) {
	for _, chunk := range t.chunks {
		frame, visible := chunk.ActiveFrame()
		if !visible {
			continue
		}
		position := chunk.CellPosition()
		localX := screenX - position.X
		if localX < 0 || localX >= frame.Width {
			continue
		}
		for y := 0; y < frame.Height; y++ {
			if frame.CellAt(localX, y).Rune != 0 {
				return position.Y + y, true
			}
		}
	}
	return 0, false
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
