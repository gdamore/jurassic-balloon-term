// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"testing"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

func TestCampaignStagePlan(t *testing.T) {
	want := []struct {
		name        string
		duration    time.Duration
		distance    float64
		checkpoints []float64
	}{
		{"Fern Valley", 2 * time.Minute, 384, []float64{192}},
		{"Red-Eye Run", 150 * time.Second, 480, []float64{240}},
		{"Thunderhead Pass", 3 * time.Minute, 756, []float64{252, 504}},
		{"Long-Neck Grove", 210 * time.Second, 672, []float64{224, 448}},
		{"Carnivore Country", 4 * time.Minute, 768, []float64{192, 384, 576}},
	}
	if StageCount() != len(want) {
		t.Fatalf("StageCount() = %d, want %d", StageCount(), len(want))
	}
	for index, expected := range want {
		stage := StageAt(index)
		if stage.Name != expected.name || stage.TargetDuration != expected.duration || stage.GoalDistance != expected.distance {
			t.Errorf("stage %d = %+v, want %q %v %.0f", index, stage, expected.name, expected.duration, expected.distance)
		}
		if len(stage.Checkpoints) != len(expected.checkpoints) {
			t.Fatalf("stage %d checkpoints = %v, want %v", index, stage.Checkpoints, expected.checkpoints)
		}
		for checkpoint := range expected.checkpoints {
			if stage.Checkpoints[checkpoint] != expected.checkpoints[checkpoint] {
				t.Errorf("stage %d checkpoint %d = %v, want %v", index, checkpoint, stage.Checkpoints[checkpoint], expected.checkpoints[checkpoint])
			}
		}
	}
}

func TestThunderheadPassUsesFasterCoherentScroll(t *testing.T) {
	m := NewStageModel(80, 24, 2)
	if m.stage.ScrollSpeed != stormScrollSpeed || m.stage.ScrollSpeed <= worldScrollSpeed {
		t.Fatalf("storm scroll speed = %v, want %v above %v", m.stage.ScrollSpeed, stormScrollSpeed, worldScrollSpeed)
	}
	if m.stage.GoalDistance != m.stage.ScrollSpeed*m.stage.TargetDuration.Seconds() {
		t.Fatalf("storm goal distance = %v, want speed × duration %v", m.stage.GoalDistance, m.stage.ScrollSpeed*m.stage.TargetDuration.Seconds())
	}
	if m.terrain[1].speed != -stormScrollSpeed {
		t.Fatalf("storm ground speed = %v, want %v", m.terrain[1].speed, -stormScrollSpeed)
	}
	treeSpeed, _ := m.trees[0].sprite.Velocity()
	finishSpeed, _ := m.finish.Velocity()
	if treeSpeed != -stormScrollSpeed || finishSpeed != -stormScrollSpeed {
		t.Fatalf("storm object speeds = tree:%v finish:%v", treeSpeed, finishSpeed)
	}
}

func TestStageDefinitionsAreCopied(t *testing.T) {
	stage := StageAt(0)
	stage.Checkpoints[0] = -1
	stage.TreeDistances[0] = -1
	stage.FernDistances[0] = -1
	stage.GrazerDistances[0] = -1
	canonical := StageAt(0)
	if canonical.Checkpoints[0] < 0 || canonical.TreeDistances[0] < 0 || canonical.FernDistances[0] < 0 || canonical.GrazerDistances[0] < 0 {
		t.Fatal("caller mutated canonical stage slices")
	}
}

func TestDistanceCompletesStageAndStopsSimulation(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	makeAirborne(m)
	m.balloon.SetPosition(20, 2)
	m.distance = m.stage.GoalDistance - worldScrollSpeed*10*time.Millisecond.Seconds()
	m.positionFinishFlag()

	m.Update(10 * time.Millisecond)
	if !m.Completed() || m.Popped() {
		t.Fatalf("state after goal = completed:%t popped:%t", m.Completed(), m.Popped())
	}
	if m.Distance() != m.stage.GoalDistance {
		t.Fatalf("distance = %v, want %v", m.Distance(), m.stage.GoalDistance)
	}
	if events := m.Events(); len(events) != 1 || events[0] != EventStageCompleted {
		t.Fatalf("completion events = %v, want stage completed", events)
	}
	flagX, _ := m.finish.Position()
	balloonX, _ := m.balloon.Position()
	if math.Abs(flagX-balloonX) > 1e-9 {
		t.Fatalf("finish x = %v, balloon x = %v", flagX, balloonX)
	}

	elapsed := m.Elapsed
	distance := m.Distance()
	m.Update(time.Second)
	if m.Elapsed != elapsed || m.Distance() != distance {
		t.Fatalf("completed stage advanced: elapsed %v→%v distance %v→%v", elapsed, m.Elapsed, distance, m.Distance())
	}
}

func TestFinishFlagTracksRemainingDistanceAndResize(t *testing.T) {
	m := NewStageModel(80, 30, 2)
	balloonX, _ := m.balloon.Position()
	flagX, flagY := m.finish.Position()
	if flagX != balloonX+m.stage.GoalDistance {
		t.Fatalf("initial finish x = %v, want %v", flagX, balloonX+m.stage.GoalDistance)
	}
	if flagY != 30-finishHeight {
		t.Fatalf("initial finish y = %v, want %d", flagY, 30-finishHeight)
	}

	m.Resize(100, 40)
	_, flagY = m.finish.Position()
	if flagY != 40-finishHeight {
		t.Fatalf("resized finish y = %v, want %d", flagY, 40-finishHeight)
	}
}

func TestContinueToStagePreservesFlightState(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	makeAirborne(m)
	m.balloon.SetPosition(24, 7.25)
	m.heat = 0.72
	m.vy = -2.4
	m.SetBurner(true)
	m.completeStage()

	if m.VerticalVelocity() != -2.4 || !m.BurnerOn() {
		t.Fatalf("intermediate completion interrupted flight: velocity=%v burner=%t", m.VerticalVelocity(), m.BurnerOn())
	}
	next := m.ContinueToStage(1)
	if next.terrain[0].coneVariant != 1 {
		t.Fatalf("continued far-terrain cone variant = %d, want 1", next.terrain[0].coneVariant)
	}
	if next.StageIndex() != 1 || next.Grounded() || next.Completed() {
		t.Fatalf("continued stage state = index:%d grounded:%t completed:%t", next.StageIndex(), next.Grounded(), next.Completed())
	}
	if next.BalloonPosition().X != 24 || next.BalloonY() != 7.25 || next.VerticalVelocity() != -2.4 || next.Heat() != 0.72 || !next.BurnerOn() {
		t.Fatalf("continued flight = position:%v y:%v velocity:%v heat:%v burner:%t", next.BalloonPosition(), next.BalloonY(), next.VerticalVelocity(), next.Heat(), next.BurnerOn())
	}
}

func TestStageTransitionCarriesVisibleFlagCloudAndPterosaur(t *testing.T) {
	previous := NewStageModel(80, 30, 0)
	makeAirborne(previous)
	previous.balloon.SetPosition(24, 10)
	previous.finish.SetPosition(18, 22)
	cloud := previous.clouds[0].sprite
	cloud.SetPosition(8, 3)
	hazard := previous.pterosaur
	hazard.setBaseY(4)
	hazard.sprite.SetPosition(52, 4)
	previous.completeStage()

	next := previous.ContinueToStage(1)
	for name, carried := range map[string]*sprite.Sprite{
		"flag":      previous.finish,
		"cloud":     cloud,
		"pterosaur": hazard.sprite,
	} {
		if !modelHasSprite(next, carried) {
			t.Fatalf("visible outgoing %s was cleared at the stage boundary", name)
		}
	}
	if len(next.departingPterosaurs) != 1 || next.departingPterosaurs[0] != hazard {
		t.Fatal("outgoing pterosaur behavior was not preserved")
	}
	for _, incoming := range next.clouds {
		x, _ := incoming.sprite.Position()
		if x < float64(next.Width) {
			t.Fatalf("replacement cloud starts at x=%v inside the playfield", x)
		}
	}

	oldX, _ := hazard.sprite.Position()
	oldPhase := hazard.phase
	next.Update(100 * time.Millisecond)
	newX, _ := hazard.sprite.Position()
	if newX >= oldX || hazard.phase == oldPhase {
		t.Fatalf("outgoing pterosaur did not continue flying: x %v -> %v, phase %v -> %v", oldX, newX, oldPhase, hazard.phase)
	}
}

func TestOutgoingPterosaurRemainsHazardous(t *testing.T) {
	previous := NewStageModel(80, 30, 0)
	makeAirborne(previous)
	previous.balloon.SetPosition(24, 10)
	hazard := previous.pterosaur
	hazard.setBaseY(10)
	hazard.sprite.SetPosition(24, 10)
	previous.completeStage()

	next := previous.ContinueToStage(1)
	next.Update(time.Millisecond)
	if !next.Popped() {
		t.Fatal("outgoing pterosaur became harmless at the stage boundary")
	}
}

func TestEnteringThunderheadPassSweepsWeatherInAndRampsWind(t *testing.T) {
	previous := NewStageModel(80, 30, 1)
	makeAirborne(previous)
	previous.balloon.SetPosition(24, 10)
	oldGround := previous.terrain[1]

	next := previous.ContinueToStage(2)
	if next.terrain[1] != oldGround {
		t.Fatal("stage transition replaced scrolling terrain")
	}
	if next.transition == nil || next.scrollSpeed != worldScrollSpeed {
		t.Fatalf("initial transition = %+v speed %v", next.transition, next.scrollSpeed)
	}
	rainX, _ := next.rain.Position()
	ceilingX, _ := next.ceiling.cloud.Position()
	if rainX != 80 || ceilingX != 80 {
		t.Fatalf("incoming weather starts at rain:%v ceiling:%v, want right edge 80", rainX, ceilingX)
	}
	if len(next.departing) == 0 {
		t.Fatal("visible prior-stage scenery did not carry across boundary")
	}

	next.updateStageTransition(time.Second)
	if next.scrollSpeed <= worldScrollSpeed || next.scrollSpeed >= stormScrollSpeed {
		t.Fatalf("mid-transition speed = %v", next.scrollSpeed)
	}
	rainX, _ = next.rain.Position()
	if rainX <= 0 || rainX >= 80 {
		t.Fatalf("mid-transition rain x = %v", rainX)
	}

	next.updateStageTransition(time.Second)
	if next.transition != nil || next.scrollSpeed != stormScrollSpeed {
		t.Fatalf("completed transition = %+v speed %v", next.transition, next.scrollSpeed)
	}
	rainX, _ = next.rain.Position()
	if rainX != 0 {
		t.Fatalf("completed rain x = %v, want 0", rainX)
	}
}

func TestLeavingThunderheadPassSweepsWeatherOutAndEasesWind(t *testing.T) {
	previous := NewStageModel(80, 30, 2)
	makeAirborne(previous)
	previous.balloon.SetPosition(24, 10)
	oldRain, oldCeiling := previous.rain, previous.ceiling

	next := previous.ContinueToStage(3)
	if next.rain != oldRain || next.ceiling != oldCeiling {
		t.Fatal("outgoing weather was replaced at stage boundary")
	}
	next.updateStageTransition(time.Second)
	rainX, _ := next.rain.Position()
	if rainX >= 0 || rainX <= -80 || next.scrollSpeed >= stormScrollSpeed || next.scrollSpeed <= worldScrollSpeed {
		t.Fatalf("mid-exit state = rain:%v speed:%v", rainX, next.scrollSpeed)
	}

	next.updateStageTransition(time.Second)
	if next.rain != nil || next.ceiling != nil || next.scrollSpeed != worldScrollSpeed || next.transition != nil {
		t.Fatalf("completed exit = rain:%v ceiling:%v speed:%v transition:%v", next.rain, next.ceiling, next.scrollSpeed, next.transition)
	}
}

func modelHasSprite(model *Model, target *sprite.Sprite) bool {
	for _, candidate := range model.Sprites() {
		if candidate == target {
			return true
		}
	}
	return false
}
