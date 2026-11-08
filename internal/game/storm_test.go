// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func TestStormTelegraphsBeforeStrike(t *testing.T) {
	m := NewModel(80, 24)
	storm := m.storm
	if storm.cloud.Collidable() {
		t.Fatal("storm cloud body is collidable")
	}
	if !storm.bolt.Collidable() {
		t.Fatal("lightning bolt is not collidable")
	}
	if storm.bolt.Visible() {
		t.Fatal("lightning is visible before telegraph")
	}

	storm.updatePhase(stormPhaseDurations[stormGather])
	if storm.phase != stormRain || storm.cloud.FrameName() != "rain-1" || storm.bolt.Visible() {
		t.Fatalf("rain phase = (%v, %q, bolt visible %v)", storm.phase, storm.cloud.FrameName(), storm.bolt.Visible())
	}
	storm.updatePhase(stormPhaseDurations[stormRain])
	if storm.phase != stormStrike || !storm.bolt.Visible() {
		t.Fatalf("strike phase = %v, bolt visible = %v", storm.phase, storm.bolt.Visible())
	}
	storm.updatePhase(stormPhaseDurations[stormStrike])
	if storm.phase != stormCooldown || storm.bolt.Visible() {
		t.Fatalf("cooldown phase = %v, bolt visible = %v", storm.phase, storm.bolt.Visible())
	}
}

func TestVisibleLightningStrikeProducesGameEvent(t *testing.T) {
	m := NewModel(80, 24)
	m.storm.setPosition(50, 0)

	m.storm.phaseElapsed = stormPhaseDurations[stormGather] - time.Millisecond
	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 0 {
		t.Fatalf("rain telegraph events = %v, want none", events)
	}
	m.storm.phaseElapsed = stormPhaseDurations[stormRain] - time.Millisecond
	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 1 || events[0] != EventLightningStruck {
		t.Fatalf("strike events = %v, want lightning strike", events)
	}
}

func TestOffscreenLightningStrikeIsSilent(t *testing.T) {
	m := NewModel(80, 24)
	m.storm.setPosition(float64(m.Width+20), 0)

	m.storm.phase = stormRain
	m.storm.phaseElapsed = stormPhaseDurations[stormRain] - time.Millisecond
	m.storm.applyPhase()
	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 0 {
		t.Fatalf("offscreen strike events = %v, want none", events)
	}
}

func TestLightningHasCompactFootprint(t *testing.T) {
	m := NewModel(80, 24)
	frame, ok := m.storm.bolt.ActiveFrame()
	if ok {
		t.Fatal("lightning is visible before strike")
	}
	if err := m.storm.bolt.SetFrame("strike"); err != nil {
		t.Fatal(err)
	}
	frame, ok = m.storm.bolt.ActiveFrame()
	if !ok {
		t.Fatal("lightning strike frame is unavailable")
	}
	if frame.Width != lightningWidth || frame.Height != 5 {
		t.Fatalf("lightning size = %dx%d, want %dx5", frame.Width, frame.Height, lightningWidth)
	}
}

func TestStormCloudBodyIsSafe(t *testing.T) {
	m := NewModel(80, 24)
	balloonX, balloonY := m.balloon.Position()
	m.storm.setPosition(balloonX, balloonY)

	m.Update(time.Millisecond)
	if m.Popped() {
		t.Fatal("safe storm cloud body popped balloon")
	}
}

func TestStormCloudHasShadedCrownAndUnderside(t *testing.T) {
	storm := newStormHazard(20, 3, 0.1)
	frame, _ := storm.cloud.ActiveFrame()
	if crown := frame.CellAt(4, 1); crown.Style != StyleStormCloudHighlight {
		t.Fatalf("storm crown style = %d, want %d", crown.Style, StyleStormCloudHighlight)
	}
	if body := frame.CellAt(4, 2); body.Style != StyleStormCloud {
		t.Fatalf("storm body style = %d, want %d", body.Style, StyleStormCloud)
	}
	if underside := frame.CellAt(4, 3); underside.Style != StyleStormCloudShadow {
		t.Fatalf("storm underside style = %d, want %d", underside.Style, StyleStormCloudShadow)
	}
}

func TestVisibleLightningPopsBalloon(t *testing.T) {
	m := NewModel(80, 30)
	balloonX, balloonY := m.balloon.Position()
	m.storm.phase = stormStrike
	m.storm.phaseElapsed = 0
	m.storm.applyPhase()
	m.storm.setPosition(balloonX-float64((stormCloudWidth-lightningWidth)/2), balloonY-5)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("visible lightning did not pop balloon")
	}
	if got := m.balloon.FrameName(); got != "pop-1" {
		t.Fatalf("balloon frame after strike = %q, want pop-1", got)
	}
	if m.storm.bolt.Visible() {
		t.Fatal("lightning remains visible after popping balloon")
	}
	if m.balloon.Layer() != popLayer {
		t.Fatalf("pop layer = %d, want %d", m.balloon.Layer(), popLayer)
	}
	if m.BurnerOn() || m.VerticalVelocity() != 0 {
		t.Fatal("flight controls remain active after pop")
	}

	m.Update(700 * time.Millisecond)
	if m.balloon.Visible() {
		t.Fatal("balloon pop animation remains visible after completion")
	}
	if !m.stopped {
		t.Fatal("simulation did not stop after pop animation")
	}
	cloudX, _ := m.storm.cloud.Position()
	elapsed := m.Elapsed
	m.Update(time.Second)
	cloudAfter, _ := m.storm.cloud.Position()
	if cloudAfter != cloudX || m.Elapsed != elapsed {
		t.Fatalf("simulation advanced after stop: cloud %v -> %v, elapsed %v -> %v", cloudX, cloudAfter, elapsed, m.Elapsed)
	}
}
