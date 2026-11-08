// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"testing"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func TestKeyReleaseCapability(t *testing.T) {
	for _, test := range []struct {
		protocol tcell.KeyProtocol
		want     bool
	}{
		{protocol: tcell.KittyKeyboard, want: true},
		{protocol: tcell.Win32Keyboard, want: true},
		{protocol: tcell.XTermKeyboard, want: false},
		{protocol: tcell.LegacyKeyboard, want: false},
	} {
		if got := supportsKeyRelease(test.protocol); got != test.want {
			t.Errorf("supportsKeyRelease(%v) = %v, want %v", test.protocol, got, test.want)
		}
	}
}

func TestSkyColorVariesWithAltitude(t *testing.T) {
	top := skyColor(0, 24)
	bottom := skyColor(23, 24)
	if top == bottom {
		t.Fatal("sky background does not vary with altitude")
	}
	if got := skyColor(-10, 24); got != top {
		t.Fatalf("sky color above screen = %v, want top color %v", got, top)
	}
	if got := skyColor(100, 24); got != bottom {
		t.Fatalf("sky color below screen = %v, want bottom color %v", got, bottom)
	}
}

func TestPartialBlockUsesUnderlyingLayerColor(t *testing.T) {
	sky := color.NewRGBColor(10, 20, 30)
	ridge := color.NewRGBColor(40, 50, 60)
	style := tcell.StyleDefault.Foreground(ridge).Background(sky)
	if got := underlyingColor("█", style, sky); got != ridge {
		t.Fatalf("background over ridge = %v, want %v", got, ridge)
	}
	if got := underlyingColor(" ", style, sky); got != sky {
		t.Fatalf("background over sky = %v, want %v", got, sky)
	}
	for _, block := range []rune{'▄', '▟', '▙'} {
		if !isPartialBlock(block) {
			t.Errorf("%q is not recognized as a partial block", block)
		}
	}
}

func TestRetryAndQuitKeys(t *testing.T) {
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, " ", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, "r", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEnter, "", tcell.ModNone),
	} {
		if !isRetryEvent(event) {
			t.Errorf("%v is not recognized as retry", event)
		}
	}
	for _, event := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyRune, "q", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone),
	} {
		if !isQuitEvent(event) {
			t.Errorf("%v is not recognized as quit", event)
		}
	}
}

func TestAdvancedBurnerInputTracksPressAndRelease(t *testing.T) {
	a := &App{
		model:      game.NewModel(80, 24),
		started:    true,
		keyRelease: true,
		burnerKeys: make(map[burnerKey]struct{}),
	}
	press := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	release := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1)

	a.handleEvent(press)
	if !a.model.BurnerOn() {
		t.Fatal("burner is off after key press")
	}
	a.handleEvent(release)
	if a.model.BurnerOn() {
		t.Fatal("burner is on after key release")
	}
}

func TestAdvancedBurnerInputTracksMultipleKeys(t *testing.T) {
	a := &App{
		model:      game.NewModel(80, 24),
		started:    true,
		keyRelease: true,
		burnerKeys: make(map[burnerKey]struct{}),
	}
	spaceDown := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	upDown := tcell.NewEventKeyEx(tcell.KeyUp, "", tcell.ModNone, true, tcell.KeyUp, 1)
	spaceUp := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1)
	upUp := tcell.NewEventKeyEx(tcell.KeyUp, "", tcell.ModNone, false, tcell.KeyUp, 1)

	for _, event := range []*tcell.EventKey{spaceDown, upDown, spaceUp} {
		a.handleEvent(event)
	}
	if !a.model.BurnerOn() {
		t.Fatal("releasing one of two held burner keys turned burner off")
	}
	a.handleEvent(upUp)
	if a.model.BurnerOn() {
		t.Fatal("burner is on after all burner keys were released")
	}
}

func TestLegacyBurnerInputStartsPulse(t *testing.T) {
	a := &App{model: game.NewModel(80, 24), started: true}
	press := tcell.NewEventKey(tcell.KeyRune, "w", tcell.ModNone)

	a.handleEvent(press)
	if !a.model.BurnerOn() {
		t.Fatal("burner is off after legacy key press")
	}
	if a.burnerUntil.IsZero() {
		t.Fatal("legacy key press did not set pulse deadline")
	}
}

type recordingSound struct {
	burner       []bool
	pops         int
	pterosaurs   int
	lightning    int
	groundStrike int
	completed    int
	music        []bool
	effects      []bool
	stages       []int
}

func (*recordingSound) Start()                   {}
func (s *recordingSound) SetBurner(on bool)      { s.burner = append(s.burner, on) }
func (s *recordingSound) BalloonPopped()         { s.pops++ }
func (s *recordingSound) PterosaurLunged()       { s.pterosaurs++ }
func (s *recordingSound) LightningStruck()       { s.lightning++ }
func (s *recordingSound) GroundLightningStruck() { s.groundStrike++ }
func (s *recordingSound) StageCompleted()        { s.completed++ }
func (s *recordingSound) SetMusicEnabled(on bool) {
	s.music = append(s.music, on)
}
func (s *recordingSound) SetEffectsEnabled(on bool) {
	s.effects = append(s.effects, on)
}
func (s *recordingSound) SetStage(stage int) { s.stages = append(s.stages, stage) }
func (*recordingSound) Retry()               {}
func (*recordingSound) SetPaused(bool)       {}

func TestSoundTracksBurnerAndGameEvents(t *testing.T) {
	sound := &recordingSound{}
	a := &App{model: game.NewStageModel(80, 24, game.StageCount()-1), sound: sound}

	a.setBurner(true)
	a.setBurner(true)
	a.setBurner(false)
	a.handleGameEvent(game.EventPterosaurLunged)
	a.handleGameEvent(game.EventLightningStruck)
	a.handleGameEvent(game.EventGroundLightningStruck)
	a.handleGameEvent(game.EventBalloonPopped)
	a.handleGameEvent(game.EventStageCompleted)

	if len(sound.burner) != 4 || !sound.burner[0] || sound.burner[1] || sound.burner[2] || sound.burner[3] {
		t.Fatalf("burner sound changes = %v, want [true false false false]", sound.burner)
	}
	if sound.pterosaurs != 1 || sound.lightning != 1 || sound.groundStrike != 1 || sound.pops != 1 || sound.completed != 1 {
		t.Fatalf("cue counts = pterosaur:%d lightning:%d ground:%d pop:%d complete:%d, want 1 each", sound.pterosaurs, sound.lightning, sound.groundStrike, sound.pops, sound.completed)
	}
}

func TestIntermediateStageCompletionDoesNotPlayFinalVictory(t *testing.T) {
	sound := &recordingSound{}
	a := &App{model: game.NewModel(80, 24), sound: sound}
	a.handleGameEvent(game.EventStageCompleted)
	if sound.completed != 0 {
		t.Fatalf("intermediate stage played %d final victory cues", sound.completed)
	}
}

func TestStartScreenTogglesAudioAndStartsBurner(t *testing.T) {
	sound := &recordingSound{}
	a := New(nil, sound)
	a.model = game.NewModel(80, 24)
	a.keyRelease = true
	a.burnerKeys = make(map[burnerKey]struct{})

	a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "m", tcell.ModNone))
	a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "s", tcell.ModNone))
	if len(sound.music) != 1 || sound.music[0] || len(sound.effects) != 1 || sound.effects[0] {
		t.Fatalf("audio toggles = music:%v effects:%v, want [false] each", sound.music, sound.effects)
	}

	space := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	a.handleEvent(space)
	if !a.started || !a.model.BurnerOn() {
		t.Fatalf("space start = started:%t burner:%t", a.started, a.model.BurnerOn())
	}
}

func TestNumberKeysSelectStagesImmediately(t *testing.T) {
	sound := &recordingSound{}
	a := New(nil, sound)
	a.model = game.NewModel(80, 24)
	a.burnerKeys = make(map[burnerKey]struct{})
	a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "2", tcell.ModNone))
	if !a.started || a.model.StageIndex() != 1 {
		t.Fatalf("stage-two selection = started:%t index:%d", a.started, a.model.StageIndex())
	}

	a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "5", tcell.ModNone))
	if a.model.StageIndex() != 4 {
		t.Fatalf("stage-five selection index = %d, want 4", a.model.StageIndex())
	}
	if len(sound.stages) != 2 || sound.stages[0] != 1 || sound.stages[1] != 4 {
		t.Fatalf("audio stage selections = %v, want [1 4]", sound.stages)
	}

	a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "6", tcell.ModNone))
	if a.model.StageIndex() != 4 {
		t.Fatalf("out-of-range stage key changed index to %d", a.model.StageIndex())
	}
}

func TestStageSelectionRequiresKeyPress(t *testing.T) {
	release := tcell.NewEventKeyEx(tcell.KeyRune, "2", tcell.ModNone, false, tcell.KeyRune, 1)
	if index, selected := selectedStage(release); selected {
		t.Fatalf("stage selected on key release: %d", index)
	}
}
