// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"testing"
	"time"

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

func TestPrimaryMouseButtonEntersThenLaunchesGameplay(t *testing.T) {
	a := New(nil)
	a.model = game.NewModel(80, 24)
	a.keyRelease = true
	a.burnerKeys = make(map[burnerKey]struct{})

	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonPrimary, tcell.ModNone))
	if !a.started || !a.mouseBurner || a.model.BurnerOn() {
		t.Fatalf("mouse press = started:%t held:%t burner:%t", a.started, a.mouseBurner, a.model.BurnerOn())
	}
	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonNone, tcell.ModNone))
	if a.mouseBurner || a.model.BurnerOn() {
		t.Fatalf("mouse release = held:%t burner:%t", a.mouseBurner, a.model.BurnerOn())
	}
	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonPrimary, tcell.ModNone))
	if !a.mouseBurner || !a.model.BurnerOn() {
		t.Fatalf("launch click = held:%t burner:%t", a.mouseBurner, a.model.BurnerOn())
	}
}

func TestMouseAndKeyboardBurnerStatesAreIndependent(t *testing.T) {
	a := &App{
		model:      game.NewModel(80, 24),
		started:    true,
		keyRelease: true,
		burnerKeys: make(map[burnerKey]struct{}),
	}
	spaceDown := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, true, tcell.KeySpace, 1)
	spaceUp := tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1)
	a.handleEvent(spaceDown)
	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonPrimary, tcell.ModNone))
	a.handleEvent(spaceUp)
	if !a.model.BurnerOn() {
		t.Fatal("keyboard release extinguished held mouse burner")
	}
	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonNone, tcell.ModNone))
	if a.model.BurnerOn() {
		t.Fatal("burner remained on after both inputs were released")
	}
}

func TestOtherMouseButtonsDoNotAffectBurner(t *testing.T) {
	a := &App{model: game.NewModel(80, 24), started: true}
	a.handleEvent(tcell.NewEventMouse(10, 10, tcell.ButtonSecondary, tcell.ModNone))
	if a.mouseBurner || a.model.BurnerOn() {
		t.Fatal("secondary mouse button activated burner")
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
	volumes      []float64
	stages       []int
	paused       []bool
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
func (s *recordingSound) SetMasterVolume(volume float64) {
	s.volumes = append(s.volumes, volume)
}
func (s *recordingSound) SetStage(stage int)    { s.stages = append(s.stages, stage) }
func (*recordingSound) Retry()                  {}
func (s *recordingSound) SetPaused(paused bool) { s.paused = append(s.paused, paused) }

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

func TestIntermediateStageCompletionShowsOverlayWithoutSound(t *testing.T) {
	a := &App{model: game.NewModel(80, 24)}
	a.handleGameEvent(game.EventStageCompleted)
	if a.overlayTitle != "STAGE CLEAR" || a.overlayDetail != game.StageAt(1).Name {
		t.Fatalf("stage overlay = %q / %q", a.overlayTitle, a.overlayDetail)
	}
	if !a.overlayUntil.After(time.Now()) {
		t.Fatal("stage overlay has no future expiry")
	}
}

func TestFractionalBlock(t *testing.T) {
	tests := map[float64]rune{0: ' ', 0.01: '▁', 0.5: '▄', 1: '█'}
	for value, want := range tests {
		if got := fractionalBlock(value); got != want {
			t.Errorf("fractionalBlock(%v) = %q, want %q", value, got, want)
		}
	}
}

func TestTitleCloudScrollsAndWraps(t *testing.T) {
	start := titleCloudX(120, 12, 0.5, 1.25, 0)
	later := titleCloudX(120, 12, 0.5, 1.25, 2*time.Second)
	if later >= start {
		t.Fatalf("title cloud moved from %d to %d, want movement left", start, later)
	}
	wrapped := titleCloudX(40, 12, 0, 10, 10*time.Second)
	if wrapped < -12 || wrapped >= 58 {
		t.Fatalf("wrapped title cloud x = %d, outside cycle", wrapped)
	}
}

func TestTitleTerrainStaysVariedAndInBounds(t *testing.T) {
	seen := map[int]bool{}
	for x := 0; x < 160; x++ {
		height := titleTerrainHeight(x, 17.5, 2, 7, 0.4)
		if height < 2 || height > 7 {
			t.Fatalf("title terrain height = %d, want 2..7", height)
		}
		seen[height] = true
	}
	if len(seen) < 3 {
		t.Fatalf("title terrain used only %d heights: %v", len(seen), seen)
	}
}

func TestStartScreenTogglesAudioAndEntersGroundedStage(t *testing.T) {
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
	if !a.started || a.model.BurnerOn() || !a.model.Grounded() {
		t.Fatalf("space start = started:%t burner:%t", a.started, a.model.BurnerOn())
	}
	a.handleEvent(tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModNone, false, tcell.KeySpace, 1))
	a.handleEvent(space)
	if !a.model.BurnerOn() {
		t.Fatal("second space press did not fire burner")
	}
}

func TestGameplayAudioControlsToggleAndClampVolume(t *testing.T) {
	sound := &recordingSound{}
	a := New(nil, sound)
	a.model = game.NewModel(80, 24)
	a.started = true

	for _, key := range []string{"m", "s", "-", "-", "+"} {
		a.handleEvent(tcell.NewEventKey(tcell.KeyRune, key, tcell.ModNone))
	}
	if a.music || a.effects {
		t.Fatalf("audio toggles = music:%t sound:%t, want both off", a.music, a.effects)
	}
	if a.volume != 9 {
		t.Fatalf("volume steps = %d, want 9", a.volume)
	}
	if len(sound.volumes) != 3 || sound.volumes[0] != 0.9 || sound.volumes[1] != 0.8 || sound.volumes[2] != 0.9 {
		t.Fatalf("master volumes = %v, want [0.9 0.8 0.9]", sound.volumes)
	}
	if a.overlayTitle != "VOLUME 90%" || !a.overlayUntil.After(time.Now()) {
		t.Fatalf("volume overlay = %q until %v", a.overlayTitle, a.overlayUntil)
	}

	for range volumeSteps + 2 {
		a.handleEvent(tcell.NewEventKey(tcell.KeyRune, "-", tcell.ModNone))
	}
	if a.volume != 0 || sound.volumes[len(sound.volumes)-1] != 0 {
		t.Fatalf("minimum volume = step:%d gain:%v", a.volume, sound.volumes[len(sound.volumes)-1])
	}
}

func TestQuitKeyReturnsGameplayToTitleAndQuitsFromTitle(t *testing.T) {
	for _, key := range []*tcell.EventKey{
		tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, "q", tcell.ModNone),
	} {
		sound := &recordingSound{}
		a := New(nil, sound)
		a.model = game.NewStageModel(80, 24, 2)
		a.model.Paused = true
		a.started = true
		a.music = false

		if quit := a.handleEvent(key); quit {
			t.Fatalf("%v quit from gameplay", key)
		}
		if a.started || a.model.StageIndex() != 0 || a.model.Paused {
			t.Fatalf("%v title return = started:%t stage:%d paused:%t", key, a.started, a.model.StageIndex(), a.model.Paused)
		}
		if a.music {
			t.Fatalf("%v title return discarded audio settings", key)
		}
		if len(sound.paused) != 1 || sound.paused[0] {
			t.Fatalf("%v pause calls = %v, want [false]", key, sound.paused)
		}
		if quit := a.handleEvent(key); !quit {
			t.Fatalf("%v did not quit from title", key)
		}
	}
}

func TestCreditsOpenAndReturnWithoutStartingOrQuitting(t *testing.T) {
	a := New(nil)
	a.model = game.NewModel(80, 24)
	c := tcell.NewEventKey(tcell.KeyRune, "c", tcell.ModNone)
	escape := tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone)

	if quit := a.handleEvent(c); quit || !a.credits || a.started {
		t.Fatalf("credits open = quit:%t credits:%t started:%t", quit, a.credits, a.started)
	}
	if quit := a.handleEvent(tcell.NewEventKey(tcell.KeyRune, " ", tcell.ModNone)); quit || a.started {
		t.Fatalf("space from credits = quit:%t started:%t", quit, a.started)
	}
	if quit := a.handleEvent(escape); quit || a.credits {
		t.Fatalf("credits escape = quit:%t credits:%t", quit, a.credits)
	}
	if quit := a.handleEvent(escape); !quit {
		t.Fatal("escape did not quit after returning to main title")
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
