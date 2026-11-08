// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

// Event reports a semantic gameplay occurrence to presentation layers.
type Event uint8

const (
	EventBalloonPopped Event = iota
	EventPterosaurLunged
	EventLightningStruck
	EventGroundLightningStruck
	EventStageCompleted
)

// Events returns the events produced by the most recent simulation update.
func (m *Model) Events() []Event {
	return append([]Event(nil), m.events...)
}
