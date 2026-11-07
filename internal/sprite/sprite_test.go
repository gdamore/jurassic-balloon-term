// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sprite

import (
	"testing"
	"time"
)

func mustFrame(t *testing.T, sprite *Sprite, name string, frame Frame) {
	t.Helper()
	if err := sprite.AddFrame(name, frame); err != nil {
		t.Fatalf("AddFrame(%q): %v", name, err)
	}
}

func TestAnimationAdvancesAndLoops(t *testing.T) {
	s := New("balloon")
	a := FrameFromStrings("A")
	a.Duration, a.Next = 100*time.Millisecond, "b"
	b := FrameFromStrings("B")
	b.Duration, b.Next = 100*time.Millisecond, "a"
	mustFrame(t, s, "a", a)
	mustFrame(t, s, "b", b)
	if err := s.SetFrame("a"); err != nil {
		t.Fatal(err)
	}

	s.Update(250 * time.Millisecond)
	if got := s.FrameName(); got != "a" {
		t.Fatalf("frame after 250ms = %q, want a", got)
	}

	s.Update(50 * time.Millisecond)
	if got := s.FrameName(); got != "b" {
		t.Fatalf("frame after 300ms = %q, want b", got)
	}
}

func TestVelocityUsesCellsPerSecond(t *testing.T) {
	s := New("moving")
	mustFrame(t, s, "visible", FrameFromStrings("X"))
	_ = s.SetFrame("visible")
	s.SetPosition(1, 2)
	s.SetVelocity(4, -2)
	s.Update(500 * time.Millisecond)

	x, y := s.Position()
	if x != 3 || y != 1 {
		t.Fatalf("position = (%v, %v), want (3, 1)", x, y)
	}
}

func TestCollisionIgnoresTransparentCells(t *testing.T) {
	a := New("a")
	b := New("b")
	mustFrame(t, a, "on", FrameFromStrings("X "))
	mustFrame(t, b, "on", FrameFromStrings(" X"))
	_ = a.SetFrame("on")
	_ = b.SetFrame("on")

	if _, collided := Collide(a, b); collided {
		t.Fatal("transparent overlap reported as a collision")
	}

	b.SetPosition(-1, 0)
	contact, collided := Collide(a, b)
	if !collided || contact != (Point{X: 0, Y: 0}) {
		t.Fatalf("collision = (%v, %v), want ({0 0}, true)", contact, collided)
	}
}

func TestSceneOrdersLayersAndReportsOneCollisionPerPair(t *testing.T) {
	scene := NewScene()
	low := New("low")
	high := New("high")
	for _, s := range []*Sprite{low, high} {
		mustFrame(t, s, "on", FrameFromStrings("XX"))
		_ = s.SetFrame("on")
	}
	low.SetLayer(1)
	high.SetLayer(9)
	scene.Add(high)
	scene.Add(low)

	callbackCount := 0
	scene.OnCollision(func(Collision) { callbackCount++ })
	collisions := scene.Update(time.Millisecond)
	if len(collisions) != 1 || callbackCount != 1 {
		t.Fatalf("collisions = %d, callbacks = %d; want 1, 1", len(collisions), callbackCount)
	}
	if got := scene.Sprites(); got[0] != low || got[1] != high {
		t.Fatal("sprites are not ordered from low to high layer")
	}
}
