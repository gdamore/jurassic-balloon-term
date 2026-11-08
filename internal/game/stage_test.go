// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"testing"
	"time"
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
		{"Thunderhead Pass", 3 * time.Minute, 576, []float64{192, 384}},
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

func TestStageDefinitionsAreCopied(t *testing.T) {
	stage := StageAt(0)
	stage.Checkpoints[0] = -1
	stage.TreeDistances[0] = -1
	canonical := StageAt(0)
	if canonical.Checkpoints[0] < 0 || canonical.TreeDistances[0] < 0 {
		t.Fatal("caller mutated canonical stage slices")
	}
}

func TestDistanceCompletesStageAndStopsSimulation(t *testing.T) {
	m := NewStageModel(80, 30, 0)
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
