// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"testing"
	"time"
)

func TestPterosaurFlapsWithDynamicMask(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 6, 0.5)
	glide, _ := pterosaur.sprite.ActiveFrame()
	if glide.Width != 15 || glide.Height != pterosaurHeight {
		t.Fatalf("pterosaur size = %dx%d, want 15x%d", glide.Width, glide.Height, pterosaurHeight)
	}
	if glide.CellAt(10, 1).Rune == 0 {
		t.Fatal("gliding wing cell is transparent")
	}
	if eye := glide.CellAt(2, 1); eye.Rune != '●' || eye.Style != StylePterosaurEye {
		t.Fatalf("eye cell = %+v, want body-embedded eye", eye)
	}
	if beak := glide.CellAt(0, 1); beak.Style != StylePterosaurBeak {
		t.Fatalf("beak style = %d, want %d", beak.Style, StylePterosaurBeak)
	}
	if beakBase := glide.CellAt(1, 1); beakBase.Rune != '▰' || beakBase.Style != StylePterosaurBeak {
		t.Fatalf("beak base = %+v, want thick beak styling", beakBase)
	}
	if crest := glide.CellAt(4, 0); crest.Style != StylePterosaurCrest {
		t.Fatalf("crest style = %d, want %d", crest.Style, StylePterosaurCrest)
	}
	if wing := glide.CellAt(9, 1); wing.Style != StylePterosaurWing {
		t.Fatalf("near wing style = %d, want %d", wing.Style, StylePterosaurWing)
	}
	if tip := glide.CellAt(10, 1); tip.Style != StylePterosaurWingHighlight {
		t.Fatalf("wing-tip style = %d, want %d", tip.Style, StylePterosaurWingHighlight)
	}
	if torso := glide.CellAt(5, 1); torso.Rune != '▀' || torso.Style != StylePterosaur {
		t.Fatalf("torso = %+v, want half-height body", torso)
	}

	pterosaur.sprite.Update(pterosaurGlideTime)
	if got := pterosaur.sprite.FrameName(); got != "up" {
		t.Fatalf("frame after glide = %q, want up", got)
	}
	up, _ := pterosaur.sprite.ActiveFrame()
	if up.CellAt(11, 1).Rune != 0 {
		t.Fatal("gliding wing remains in raised-wing collision mask")
	}
	pterosaur.sprite.Update(pterosaurWingUpTime + pterosaurWingDownTime)
	if got := pterosaur.sprite.FrameName(); got != "call" {
		t.Fatalf("frame after quick flap = %q, want open-mouth call", got)
	}
	call, _ := pterosaur.sprite.ActiveFrame()
	if call.CellAt(0, 0).Rune != '◂' || call.CellAt(0, 2).Rune != '◂' {
		t.Fatal("ordinary pterosaur call frame does not open both jaws")
	}
	pterosaur.sprite.Update(pterosaurCallTime)
	if got := pterosaur.sprite.FrameName(); got != "glide" {
		t.Fatalf("frame after call = %q, want glide", got)
	}
}

func TestPterosaurPathIsBounded(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 8, 0.5)
	for range 100 {
		pterosaur.updatePath(33 * time.Millisecond)
		_, y := pterosaur.sprite.Position()
		if math.Abs(y-pterosaur.baseY) > pterosaurBobHeight+1e-9 {
			t.Fatalf("pterosaur y = %v outside bob range around %v", y, pterosaur.baseY)
		}
	}
}

func TestRedEyePterosaurHasDistinctPaletteAndSpeed(t *testing.T) {
	pterosaur := newProfiledPterosaurHazard(redEyePterosaurProfile, 20, 8, 0.5)
	frame, _ := pterosaur.sprite.ActiveFrame()
	if frame.Width != 18 || frame.Height != redEyePterosaurHeight {
		t.Fatalf("red-eye size = %dx%d, want 18x%d", frame.Width, frame.Height, redEyePterosaurHeight)
	}

	if eye := frame.CellAt(2, 1); eye.Rune != '●' || eye.Style != StyleRedEyePterosaurEye {
		t.Fatalf("eye cell = %+v, want red-eye styling", eye)
	}
	if body := frame.CellAt(3, 1); body.Style != StyleRedEyePterosaur {
		t.Fatalf("body style = %d, want %d", body.Style, StyleRedEyePterosaur)
	}
	if crest := frame.CellAt(4, 0); crest.Style != StyleRedEyePterosaurCrest {
		t.Fatalf("crest style = %d, want %d", crest.Style, StyleRedEyePterosaurCrest)
	}
	if wing := frame.CellAt(9, 1); wing.Style != StyleRedEyePterosaurWing {
		t.Fatalf("wing style = %d, want %d", wing.Style, StyleRedEyePterosaurWing)
	}
	velocity, _ := pterosaur.sprite.Velocity()
	if velocity != 2*pterosaurSpeed {
		t.Fatalf("cruise velocity = %v, want 2x ordinary speed %v", velocity, 2*pterosaurSpeed)
	}
	if velocity != redEyePterosaurSpeed {
		t.Fatalf("cruise velocity = %v, want profile speed %v", velocity, redEyePterosaurSpeed)
	}
	if pterosaur.profile.boostSpeed != redEyeBoostSpeed {
		t.Fatalf("boost velocity = %v, want capped boost %v", pterosaur.profile.boostSpeed, redEyeBoostSpeed)
	}
}

func TestRedEyePterosaurStartsInSecondStage(t *testing.T) {
	first := NewStageModel(80, 30, 0)
	if first.redEye != nil {
		t.Fatal("red-eye pterosaur appears in Fern Valley")
	}

	second := NewStageModel(80, 30, 1)
	if second.redEye == nil {
		t.Fatal("Red-Eye Run has no red-eye pterosaur")
	}
	if second.pterosaur == nil {
		t.Fatal("Red-Eye Run lost the ordinary pterosaur")
	}
}

func TestRedEyePterosaurCollisionPopsBalloon(t *testing.T) {
	m := NewStageModel(80, 30, 1)
	makeAirborne(m)
	balloonX, balloonY := m.balloon.Position()
	m.redEye.baseY = balloonY
	m.redEye.phase = 0
	m.redEye.sprite.SetPosition(balloonX+0.5, balloonY)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("red-eye pterosaur contact did not pop balloon")
	}
}

func TestPterosaurFlapProducesDecayingSpeedBoost(t *testing.T) {
	pterosaur := newPterosaurHazard(20, 8, 0.5)

	pterosaur.sprite.Update(pterosaurGlideTime)
	pterosaur.updatePath(0)
	pterosaur.sprite.Update(pterosaurWingUpTime)
	pterosaur.updatePath(0)
	pterosaur.sprite.Update(pterosaurWingDownTime)
	pterosaur.updatePath(0)

	boosted, _ := pterosaur.sprite.Velocity()
	if boosted != pterosaurBoostSpeed {
		t.Fatalf("velocity after flap = %v, want %v", boosted, pterosaurBoostSpeed)
	}
	_, lungeY := pterosaur.sprite.Position()
	lungePhase := pterosaur.phase

	pterosaur.updatePath(pterosaurBoostTime / 2)
	decaying, _ := pterosaur.sprite.Velocity()
	if decaying <= boosted || decaying >= pterosaurSpeed {
		t.Fatalf("decaying velocity = %v, want between %v and %v", decaying, boosted, pterosaurSpeed)
	}
	_, y := pterosaur.sprite.Position()
	if y != lungeY || pterosaur.phase != lungePhase {
		t.Fatalf("lunge changed vertical path: y = %v, phase = %v", y, pterosaur.phase)
	}

	pterosaur.updatePath(pterosaurBoostTime / 2)
	cruising, _ := pterosaur.sprite.Velocity()
	if cruising != pterosaurSpeed {
		t.Fatalf("velocity after boost = %v, want %v", cruising, pterosaurSpeed)
	}
	_, y = pterosaur.sprite.Position()
	if y != lungeY || pterosaur.phase != lungePhase {
		t.Fatalf("last lunge step changed vertical path: y = %v, phase = %v", y, pterosaur.phase)
	}

	pterosaur.updatePath(time.Millisecond)
	if pterosaur.phase == lungePhase {
		t.Fatal("vertical path did not resume after lunge")
	}
}

func TestPterosaurCollisionPopsBalloon(t *testing.T) {
	m := NewModel(80, 30)
	makeAirborne(m)
	balloonX, balloonY := m.balloon.Position()
	m.pterosaur.baseY = balloonY
	m.pterosaur.phase = 0
	m.pterosaur.sprite.SetPosition(balloonX+0.5, balloonY)

	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("pterosaur contact did not pop balloon")
	}
	if events := m.Events(); len(events) != 1 || events[0] != EventBalloonPopped {
		t.Fatalf("collision events = %v, want balloon pop", events)
	}
}

func TestPterosaurLungeProducesGameEvent(t *testing.T) {
	m := NewModel(80, 30)
	makeAirborne(m)
	m.pterosaur.lastFrame = "down"
	m.pterosaur.sprite.SetPosition(40, 6)
	if err := m.pterosaur.sprite.SetFrame("call"); err != nil {
		t.Fatal(err)
	}

	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 1 || events[0] != EventPterosaurLunged {
		t.Fatalf("lunge events = %v, want pterosaur lunge", events)
	}
	if m.pterosaur.sprite.FrameName() != "call" {
		t.Fatalf("cry began on frame %q, want open-mouth call", m.pterosaur.sprite.FrameName())
	}

	m.pterosaur.lastFrame = "down"
	m.pterosaur.sprite.SetPosition(float64(m.Width+1), 6)
	m.Update(time.Millisecond)
	if events := m.Events(); len(events) != 0 {
		t.Fatalf("offscreen lunge events = %v, want none", events)
	}
}

func TestAuthoredPterosaurGroupsAreLimitedAndStaggered(t *testing.T) {
	first := NewStageModel(120, 40, 0)
	first.distance = first.stage.PterosaurGroups[0].Distance
	first.spawnScheduledPterosaurGroups()
	if len(first.flock) != 2 {
		t.Fatalf("Fern Valley flock size = %d, want 2", len(first.flock))
	}
	if first.flock[0].sprite.FrameName() != "glide" || first.flock[1].sprite.FrameName() != "delay" {
		t.Fatalf("Fern Valley flock frames = %q, %q, want staggered starts", first.flock[0].sprite.FrameName(), first.flock[1].sprite.FrameName())
	}
	first.spawnScheduledPterosaurGroups()
	if len(first.flock) != 2 {
		t.Fatalf("authored flock spawned twice: size %d", len(first.flock))
	}

	second := NewStageModel(120, 40, 1)
	second.distance = second.stage.PterosaurGroups[0].Distance
	second.spawnScheduledPterosaurGroups()
	if len(second.flock) != 3 {
		t.Fatalf("Red-Eye Run flock size = %d, want 3", len(second.flock))
	}
	redEyes := 0
	for _, hazard := range second.flock {
		if hazard.profile.large {
			redEyes++
		}
	}
	if redEyes != 1 {
		t.Fatalf("Red-Eye Run flock has %d red-eyes, want 1", redEyes)
	}
}

func TestLongNeckGroveSchedulesSparseHighPterosaurs(t *testing.T) {
	grove := NewStageModel(120, 40, 3)
	if grove.stage.Pterosaurs {
		t.Fatal("Long-Neck Grove enables routine random pterosaurs")
	}
	if len(grove.stage.PterosaurGroups) != 3 {
		t.Fatalf("Long-Neck Grove pterosaur encounters = %d, want 3", len(grove.stage.PterosaurGroups))
	}
	for index, group := range grove.stage.PterosaurGroups {
		if group.Count != 1 || group.RedEyes != 0 || group.Altitude <= 0 || group.Altitude > 0.20 {
			t.Errorf("encounter %d = %+v, want one ordinary high-flying pterosaur", index, group)
		}
	}

	grove.distance = grove.stage.PterosaurGroups[0].Distance
	grove.spawnScheduledPterosaurGroups()
	if len(grove.flock) != 1 {
		t.Fatalf("first grove encounter spawned %d pterosaurs, want 1", len(grove.flock))
	}
	if got, want := grove.flock[0].yFraction, grove.stage.PterosaurGroups[0].Altitude; got != want {
		t.Fatalf("grove pterosaur altitude fraction = %v, want %v", got, want)
	}
}

func TestDepartedPterosaurIsReplaced(t *testing.T) {
	m := NewModel(80, 24)
	makeAirborne(m)
	departed := m.pterosaur.sprite
	frame, _ := departed.ActiveFrame()
	departed.SetPosition(float64(-frame.Width-1), 4)

	m.Update(time.Millisecond)
	if m.pterosaur.sprite == departed {
		t.Fatal("departed pterosaur was repositioned instead of replaced")
	}
	x, _ := m.pterosaur.sprite.Position()
	if x <= float64(m.Width) {
		t.Fatalf("replacement pterosaur x = %v, want beyond right edge %d", x, m.Width)
	}
	for _, item := range m.Sprites() {
		if item == departed {
			t.Fatal("departed pterosaur remains in scene")
		}
	}
}

func TestDepartedRedEyePterosaurIsReplacedIndependently(t *testing.T) {
	m := NewStageModel(80, 24, 1)
	makeAirborne(m)
	ordinary := m.pterosaur.sprite
	departed := m.redEye.sprite
	frame, _ := departed.ActiveFrame()
	departed.SetPosition(float64(-frame.Width-1), 4)

	m.Update(time.Millisecond)
	if m.redEye.sprite == departed {
		t.Fatal("departed red-eye pterosaur was not replaced")
	}
	if m.pterosaur.sprite != ordinary {
		t.Fatal("replacing red-eye pterosaur also replaced ordinary pterosaur")
	}
	x, _ := m.redEye.sprite.Position()
	if x <= float64(m.Width) {
		t.Fatalf("replacement red-eye pterosaur x = %v, want beyond right edge %d", x, m.Width)
	}
}
