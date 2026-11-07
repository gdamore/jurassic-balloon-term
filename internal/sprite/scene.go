// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sprite

import (
	"slices"
	"sort"
	"time"
)

// Collision describes a cell-accurate overlap between two sprites.
type Collision struct {
	A       *Sprite
	B       *Sprite
	Contact Point
}

// Scene owns sprites, updates them, orders them for rendering, and detects
// collisions after each update.
type Scene struct {
	sprites  []*Sprite
	handlers []func(Collision)
}

func NewScene() *Scene { return &Scene{} }

func (s *Scene) Add(sprite *Sprite) {
	if sprite == nil || slices.Contains(s.sprites, sprite) {
		return
	}
	s.sprites = append(s.sprites, sprite)
}

func (s *Scene) Remove(sprite *Sprite) {
	for index, candidate := range s.sprites {
		if candidate == sprite {
			s.sprites = slices.Delete(s.sprites, index, index+1)
			return
		}
	}
}

// Sprites returns a stable layer-ordered snapshot for rendering. Sprites in
// the same layer retain insertion order.
func (s *Scene) Sprites() []*Sprite {
	result := slices.Clone(s.sprites)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Layer() < result[j].Layer()
	})
	return result
}

// OnCollision registers a callback invoked after each colliding pair is
// detected. Update also returns all collisions for callers that prefer polling.
func (s *Scene) OnCollision(handler func(Collision)) {
	if handler != nil {
		s.handlers = append(s.handlers, handler)
	}
}

func (s *Scene) Update(step time.Duration) []Collision {
	for _, sprite := range s.sprites {
		sprite.Update(step)
	}

	collisions := make([]Collision, 0)
	ordered := s.Sprites()
	for i := len(ordered) - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			contact, ok := Collide(ordered[i], ordered[j])
			if !ok {
				continue
			}
			collision := Collision{A: ordered[i], B: ordered[j], Contact: contact}
			collisions = append(collisions, collision)
			for _, handler := range s.handlers {
				handler(collision)
			}
		}
	}
	return collisions
}

// Collide checks bounding boxes first, then tests only opaque cells in their
// intersection. It returns the first contact point in row-major order.
func Collide(a, b *Sprite) (Point, bool) {
	if a == nil || b == nil || a == b {
		return Point{}, false
	}
	aBounds, aVisible := a.Bounds()
	bBounds, bVisible := b.Bounds()
	if !aVisible || !bVisible {
		return Point{}, false
	}

	minX := max(aBounds.Min.X, bBounds.Min.X)
	minY := max(aBounds.Min.Y, bBounds.Min.Y)
	maxX := min(aBounds.Max.X, bBounds.Max.X)
	maxY := min(aBounds.Max.Y, bBounds.Max.Y)
	if minX > maxX || minY > maxY {
		return Point{}, false
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if a.CollidesAt(x, y) && b.CollidesAt(x, y) {
				return Point{X: x, Y: y}, true
			}
		}
	}
	return Point{}, false
}
