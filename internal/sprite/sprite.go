// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sprite

import (
	"fmt"
	"time"
)

// Point is an integer cell coordinate.
type Point struct {
	X int
	Y int
}

// Rect is an inclusive rectangle in scene coordinates.
type Rect struct {
	Min Point
	Max Point
}

// Sprite is an animated, moving collection of terminal cells.
type Sprite struct {
	name         string
	frames       map[string]Frame
	frame        string
	frameElapsed time.Duration
	x, y         float64
	origin       Point
	vx, vy       float64
	layer        int
}

func New(name string) *Sprite {
	return &Sprite{name: name, frames: make(map[string]Frame)}
}

func (s *Sprite) Name() string { return s.name }

// AddFrame registers an animation frame. All frames in a sprite must have the
// same dimensions so its collision bounds do not jump during animation.
func (s *Sprite) AddFrame(name string, frame Frame) error {
	if name == "" {
		return fmt.Errorf("frame name must not be empty")
	}
	if frame.Width <= 0 || frame.Height <= 0 {
		return fmt.Errorf("frame %q has invalid size %dx%d", name, frame.Width, frame.Height)
	}
	if len(frame.Cells) != frame.Width*frame.Height {
		return fmt.Errorf("frame %q has %d cells, want %d", name, len(frame.Cells), frame.Width*frame.Height)
	}
	for existingName, existing := range s.frames {
		if existing.Width != frame.Width || existing.Height != frame.Height {
			return fmt.Errorf("frame %q is %dx%d; frame %q is %dx%d", name, frame.Width, frame.Height, existingName, existing.Width, existing.Height)
		}
		break
	}
	s.frames[name] = frame
	return nil
}

// SetFrame activates a frame and restarts its animation timer. An empty name
// hides the sprite.
func (s *Sprite) SetFrame(name string) error {
	if name != "" {
		if _, ok := s.frames[name]; !ok {
			return fmt.Errorf("unknown frame %q", name)
		}
	}
	s.frame = name
	s.frameElapsed = 0
	return nil
}

func (s *Sprite) FrameName() string { return s.frame }

func (s *Sprite) ActiveFrame() (Frame, bool) {
	frame, ok := s.frames[s.frame]
	return frame, ok
}

func (s *Sprite) Visible() bool {
	_, ok := s.ActiveFrame()
	return ok
}

func (s *Sprite) Hide() { _ = s.SetFrame("") }

func (s *Sprite) SetOrigin(x, y int) { s.origin = Point{X: x, Y: y} }

func (s *Sprite) SetPosition(x, y float64) { s.x, s.y = x, y }

func (s *Sprite) Position() (float64, float64) { return s.x, s.y }

func (s *Sprite) CellPosition() Point {
	return Point{X: int(s.x), Y: int(s.y)}
}

// SetVelocity sets motion in terminal cells per second.
func (s *Sprite) SetVelocity(x, y float64) { s.vx, s.vy = x, y }

func (s *Sprite) Velocity() (float64, float64) { return s.vx, s.vy }

func (s *Sprite) SetLayer(layer int) { s.layer = layer }

func (s *Sprite) Layer() int { return s.layer }

// Update advances movement and animation by an explicit simulation step.
func (s *Sprite) Update(step time.Duration) {
	if step <= 0 {
		return
	}
	seconds := step.Seconds()
	s.x += s.vx * seconds
	s.y += s.vy * seconds

	s.frameElapsed += step
	// The cap protects against pathological animation graphs while allowing a
	// large simulation step to cross several ordinary frames.
	for range 1024 {
		frame, ok := s.ActiveFrame()
		if !ok || frame.Duration <= 0 || s.frameElapsed < frame.Duration {
			return
		}
		s.frameElapsed -= frame.Duration
		if frame.Next == "" {
			s.Hide()
			return
		}
		if _, ok := s.frames[frame.Next]; !ok {
			s.Hide()
			return
		}
		s.frame = frame.Next
	}
}

func (s *Sprite) Bounds() (Rect, bool) {
	frame, ok := s.ActiveFrame()
	if !ok {
		return Rect{}, false
	}
	position := s.CellPosition()
	min := Point{X: position.X - s.origin.X, Y: position.Y - s.origin.Y}
	return Rect{Min: min, Max: Point{X: min.X + frame.Width - 1, Y: min.Y + frame.Height - 1}}, true
}

// CellAt returns the sprite cell at a scene coordinate.
func (s *Sprite) CellAt(x, y int) Cell {
	frame, ok := s.ActiveFrame()
	if !ok {
		return Cell{}
	}
	position := s.CellPosition()
	return frame.CellAt(x-position.X+s.origin.X, y-position.Y+s.origin.Y)
}

func (s *Sprite) CollidesAt(x, y int) bool { return s.CellAt(x, y).Rune != 0 }
