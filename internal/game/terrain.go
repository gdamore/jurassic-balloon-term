// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"math/rand"
	"time"

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
	frameHeight int
	minHeight   int
	lastHeight  int
	surface     int
	fill        int
	hazard      bool
	chunkID     int
	target      int
	segmentLeft int
	slopePhase  int
	coneVariant int
}

func newTerrainLayer(name string, seed int64, speed float64, layer, minHeight, maxHeight, surface, fill int, hazard bool) *terrainLayer {
	return &terrainLayer{
		name:        name,
		random:      rand.New(rand.NewSource(seed)),
		speed:       speed,
		layer:       layer,
		minHeight:   minHeight,
		maxHeight:   maxHeight,
		frameHeight: maxHeight,
		lastHeight:  minHeight,
		surface:     surface,
		fill:        fill,
		hazard:      hazard,
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
	if frameHasStyle(frame, StyleFarTerrainLava) {
		glow := frame
		glow.Cells = append([]sprite.Cell(nil), frame.Cells...)
		for index := range glow.Cells {
			if glow.Cells[index].Style == StyleFarTerrainLava {
				glow.Cells[index].Style = StyleFarTerrainLavaGlow
			}
		}
		frame.Duration = time.Duration(620+t.random.Intn(181)) * time.Millisecond
		frame.Next = "terrain-glow"
		glow.Duration = 280 * time.Millisecond
		glow.Next = "terrain"
		if err := chunk.AddFrame("terrain-glow", glow); err != nil {
			panic(err)
		}
	}
	if err := chunk.AddFrame("terrain", frame); err != nil {
		panic(err)
	}
	if err := chunk.SetFrame("terrain"); err != nil {
		panic(err)
	}
	chunk.SetLayer(t.layer)
	chunk.SetCollidable(t.hazard)
	chunk.SetPosition(x, float64(screenHeight-t.frameHeight))
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

	// Starting beyond the opening view, every other far-terrain chunk carries
	// a broad, cratered cinder cone.
	coneCenter := -1
	coneVariant := 0
	if t.name == "far-terrain" && t.chunkID >= 5 && (t.chunkID-5)%2 == 0 {
		coneIndex := (t.chunkID - 5) / 2
		coneVariant = (coneIndex + t.coneVariant) % 3
		centerOffsets := [...]int{-2, 1, 0}
		peakInsets := [...]int{1, 2, 1}
		craterWidths := [...]int{2, 3, 2}
		craterDepths := [...]int{2, 1, 3}
		center := terrainChunkWidth/2 + centerOffsets[coneVariant]
		coneCenter = center
		for x := range heights {
			distance := abs(x - center)
			peak := t.frameHeight - peakInsets[coneVariant]
			rise := peak - t.minHeight
			coneHeight := peak - distance*rise/(terrainChunkWidth/2)
			heights[x] = max(heights[x], coneHeight)
		}
		craterWidth := craterWidths[coneVariant]
		rim := min(heights[center-craterWidth], heights[center+craterWidth])
		crater := max(t.minHeight, rim-craterDepths[coneVariant])
		heights[center] = crater
		heights[center-1] = min(heights[center-1], crater+1)
		heights[center+1] = min(heights[center+1], crater+1)
	}
	t.lastHeight = heights[len(heights)-1]

	frame := sprite.Frame{
		Width:  terrainChunkWidth,
		Height: t.frameHeight,
		Cells:  make([]sprite.Cell, terrainChunkWidth*t.frameHeight),
	}
	for x, columnHeight := range heights {
		top := t.frameHeight - columnHeight
		surface := '▄'
		if x > 0 {
			surface = terrainSurface(columnHeight, heights[x-1])
		}
		frame.Cells[top*terrainChunkWidth+x] = sprite.Cell{Rune: surface, Style: t.surface}
		for y := top + 1; y < t.frameHeight; y++ {
			frame.Cells[y*terrainChunkWidth+x] = sprite.Cell{Rune: '█', Style: t.fill}
		}
	}
	if coneCenter >= 0 {
		direction := 1
		if ((t.chunkID-5)/2+t.coneVariant)%2 != 0 {
			direction = -1
		}
		addLavaChannel(&frame, heights, coneCenter, direction, coneVariant)
	}
	return frame
}

func addLavaChannel(frame *sprite.Frame, heights []int, center, direction, variant int) {
	startDrops := [...]int{1, 2, 1}
	branchLengths := [...]int{5, 4, 6}
	startY := frame.Height - heights[center] + startDrops[variant]
	mainOffsets := []int{0, 1, 2, 1, 0, -1, -2, -1, 0, 1, 0}
	branchOffsets := make([]int, branchLengths[variant])
	for index := range branchOffsets {
		branchOffsets[index] = index
	}
	drawLavaPath(frame, heights, center, direction, startY, mainOffsets)
	drawLavaPath(frame, heights, center, direction, startY+4, branchOffsets)

	junctionX := center
	junctionY := startY + 4
	if junctionY < frame.Height && junctionY > frame.Height-heights[junctionX] {
		index := junctionY*frame.Width + junctionX
		frame.Cells[index] = sprite.Cell{Rune: '█', Style: StyleFarTerrainLava}
	}
}

func drawLavaPath(frame *sprite.Frame, heights []int, center, direction, startY int, offsets []int) {
	for step, offset := range offsets {
		y := startY + step
		if y >= frame.Height {
			break
		}
		x := center + direction*offset
		if x < 0 || x >= frame.Width {
			continue
		}
		if y <= frame.Height-heights[x] {
			continue
		}
		index := y*frame.Width + x
		if frame.Cells[index].Rune == 0 {
			continue
		}
		r := '▌'
		if step+1 < len(offsets) {
			nextX := center + direction*offsets[step+1]
			switch {
			case nextX > x:
				r = '▚'
			case nextX < x:
				r = '▞'
			}
		}
		frame.Cells[index].Rune = r
		frame.Cells[index].Style = StyleFarTerrainLava
	}
}

func frameHasStyle(frame sprite.Frame, style int) bool {
	for _, cell := range frame.Cells {
		if cell.Rune != 0 && cell.Style == style {
			return true
		}
	}
	return false
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
		chunk.SetPosition(x, float64(height-t.frameHeight))
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
