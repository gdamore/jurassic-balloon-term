// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import "time"

const worldScrollSpeed = 3.2

// Stage describes an authored stretch of the journey. Checkpoints are
// distances from the stage start; persistence and restart behavior will be
// layered on once stage encounters are authored.
type Stage struct {
	Name           string
	TargetDuration time.Duration
	GoalDistance   float64
	Checkpoints    []float64
	TreeDistances  []float64
}

var stages = []Stage{
	withTrees(newStage("Fern Valley", 2*time.Minute, time.Minute), 0.52, 0.62, 0.72, 0.82, 0.91),
	newStage("Red-Eye Run", 150*time.Second, 75*time.Second),
	newStage("Thunderhead Pass", 3*time.Minute, time.Minute, 2*time.Minute),
	newStage("Long-Neck Grove", 210*time.Second, 70*time.Second, 140*time.Second),
	newStage("Carnivore Country", 4*time.Minute, time.Minute, 2*time.Minute, 3*time.Minute),
}

func withTrees(stage Stage, fractions ...float64) Stage {
	stage.TreeDistances = make([]float64, len(fractions))
	for index, fraction := range fractions {
		stage.TreeDistances[index] = stage.GoalDistance * fraction
	}
	return stage
}

func newStage(name string, duration time.Duration, checkpoints ...time.Duration) Stage {
	stage := Stage{
		Name:           name,
		TargetDuration: duration,
		GoalDistance:   worldScrollSpeed * duration.Seconds(),
		Checkpoints:    make([]float64, len(checkpoints)),
	}
	for index, checkpoint := range checkpoints {
		stage.Checkpoints[index] = worldScrollSpeed * checkpoint.Seconds()
	}
	return stage
}

// StageCount returns the number of stages in the initial journey.
func StageCount() int { return len(stages) }

// StageAt returns a copy of a stage definition. Out-of-range indexes select
// the first stage so callers always receive a playable definition.
func StageAt(index int) Stage {
	if index < 0 || index >= len(stages) {
		index = 0
	}
	stage := stages[index]
	stage.Checkpoints = append([]float64(nil), stage.Checkpoints...)
	stage.TreeDistances = append([]float64(nil), stage.TreeDistances...)
	return stage
}
