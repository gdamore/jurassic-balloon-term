// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import "time"

const worldScrollSpeed = 3.2

// Stage describes an authored stretch of the journey. Checkpoints are
// distances from the stage start; persistence and restart behavior will be
// layered on once stage encounters are authored.
type Stage struct {
	Name            string
	TargetDuration  time.Duration
	GoalDistance    float64
	Checkpoints     []float64
	TreeDistances   []float64
	FernDistances   []float64
	GrazerDistances []float64
}

var stages = []Stage{
	withGrazers(
		withFerns(
			withTrees(newStage("Fern Valley", 2*time.Minute, time.Minute), 0.52, 0.62, 0.72, 0.82, 0.91),
			0.06, 0.14, 0.23, 0.33, 0.43,
		),
		0.10, 0.125, 0.28, 0.305, 0.45, 0.475,
	),
	withGrazers(
		withFerns(
			withTrees(newStage("Red-Eye Run", 150*time.Second, 75*time.Second), 0.16, 0.34, 0.52, 0.70, 0.88),
			0.05, 0.22, 0.42, 0.63, 0.82,
		),
		0.10, 0.125, 0.43, 0.455, 0.76, 0.785,
	),
	newStage("Thunderhead Pass", 3*time.Minute, time.Minute, 2*time.Minute),
	newStage("Long-Neck Grove", 210*time.Second, 70*time.Second, 140*time.Second),
	newStage("Carnivore Country", 4*time.Minute, time.Minute, 2*time.Minute, 3*time.Minute),
}

func withGrazers(stage Stage, fractions ...float64) Stage {
	stage.GrazerDistances = make([]float64, len(fractions))
	for index, fraction := range fractions {
		stage.GrazerDistances[index] = stage.GoalDistance * fraction
	}
	return stage
}

func withFerns(stage Stage, fractions ...float64) Stage {
	stage.FernDistances = make([]float64, len(fractions))
	for index, fraction := range fractions {
		stage.FernDistances[index] = stage.GoalDistance * fraction
	}
	return stage
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
	stage.FernDistances = append([]float64(nil), stage.FernDistances...)
	stage.GrazerDistances = append([]float64(nil), stage.GrazerDistances...)
	return stage
}
