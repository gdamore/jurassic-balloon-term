// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	gravity          = 5.0
	maximumLift      = 11.0
	heatResponse     = 350 * time.Millisecond
	coolingResponse  = 900 * time.Millisecond
	maximumRiseSpeed = 5.0
	maximumFallSpeed = 7.0
)

// Renderer palette indexes shared by game-created sprite frames.
const (
	StyleBalloon = iota
	StyleBalloonFill
	StyleBalloonHighlightFill
	StyleBalloonShadowFill
	StyleBalloonShade
	StyleFlame
	StylePatch
	StylePilot
	StyleFarCloud
	StyleFarCloudHighlight
	StyleFarCloudShadow
	StyleNearCloud
	StyleNearCloudHighlight
	StyleNearCloudShadow
	StyleStormCloud
	StyleStormCloudHighlight
	StyleStormCloudShadow
	StyleRain
	StyleLightning
	StylePop
	StyleFarTerrain
	StyleGroundSurface
	StyleGroundFill
	StylePterosaur
	StylePterosaurCrest
	StylePterosaurEye
	StylePterosaurWing
	StylePterosaurWingHighlight
	StylePterosaurShadow
	StylePterosaurBeak
	StyleRedEyePterosaur
	StyleRedEyePterosaurCrest
	StyleRedEyePterosaurEye
	StyleRedEyePterosaurWing
	StyleRedEyePterosaurWingHighlight
	StyleRedEyePterosaurShadow
	StyleRedEyePterosaurBeak
	StyleFinishFlag
	StyleFinishPole
	StyleTreeCanopy
	StyleTreeHighlight
	StyleTreeTrunk
	StyleFernLeaf
	StyleFernHighlight
	StyleFernStem
	StyleGrazerBody
	StyleGrazerHighlight
	StyleGrazerShadow
	StyleGrazerHorn
	StyleGrazerHead
	StyleGrazerBeak
	StyleGrazerLeg
)

const (
	farTerrainLayer = 0
	farCloudLayer   = 2
	groundLayer     = 8
	balloonLayer    = 10
	nearCloudLayer  = 20
	popLayer        = 30
)

const balloonPopDuration = 620 * time.Millisecond

type cloudTrack struct {
	sprite    *sprite.Sprite
	depth     cloudDepth
	yFraction float64
}

type cloudDepth int

const (
	farCloud cloudDepth = iota
	nearCloud
)

// Model contains game state without any terminal-specific dependencies.
type Model struct {
	Width   int
	Height  int
	Paused  bool
	Elapsed time.Duration

	scene      *sprite.Scene
	balloon    *sprite.Sprite
	clouds     []cloudTrack
	storm      *stormHazard
	pterosaur  *pterosaurHazard
	redEye     *pterosaurHazard
	finish     *sprite.Sprite
	trees      []treeTrack
	ferns      []fernTrack
	grazers    []grazerTrack
	terrain    []*terrainLayer
	random     *rand.Rand
	cloudID    int
	burner     bool
	heat       float64
	vy         float64
	popped     bool
	stopped    bool
	popTime    time.Duration
	events     []Event
	stage      Stage
	stageIndex int
	distance   float64
	completed  bool
	grounded   bool
}

// NewModel creates a model centered in the supplied play field.
func NewModel(width, height int) *Model {
	return NewStageModel(width, height, 0)
}

// NewStageModel creates a model for one stage of the journey.
func NewStageModel(width, height, stageIndex int) *Model {
	if stageIndex < 0 || stageIndex >= StageCount() {
		stageIndex = 0
	}
	m := &Model{
		scene:      sprite.NewScene(),
		balloon:    newBalloon(),
		finish:     newFinishFlag(),
		random:     rand.New(rand.NewSource(int64(stageIndex + 1))),
		stage:      StageAt(stageIndex),
		stageIndex: stageIndex,
		grounded:   true,
	}
	m.scene.Add(m.balloon)
	m.Resize(width, height)
	m.addTerrain()
	m.placeBalloonOnGround()
	m.addStageFerns()
	m.addStageGrazers()
	m.addStageTrees()
	m.addClouds()
	m.spawnStorm(float64(m.Width + 10))
	m.spawnPterosaur(float64(m.Width + 35))
	if stageIndex >= 1 {
		m.spawnRedEye(float64(m.Width + 90))
	}
	m.positionFinishFlag()
	m.scene.Add(m.finish)
	return m
}

func newBalloon() *sprite.Sprite {
	balloon := sprite.New("balloon")
	off := balloonFrame(
		"  ▗▄▄▄▖  ",
		" ▟     ▙ ",
		"▐  ╳    ▌",
		" ▜     ▛ ",
		"  ▝▄┬▄▘  ",
		"   ╲ ╱   ",
		"   ╰☺╯   ",
	)
	burn1 := balloonFrame(
		"  ▗▄▄▄▖  ",
		" ▟     ▙ ",
		"▐  ╳    ▌",
		" ▜     ▛ ",
		"  ▝▄┴▄▘  ",
		"   ╲▲╱   ",
		"   ╰☺╯   ",
	)
	burn1.Duration, burn1.Next = 100*time.Millisecond, "burn-2"
	burn2 := balloonFrame(
		"  ▗▄▄▄▖  ",
		" ▟     ▙ ",
		"▐  ╳    ▌",
		" ▜     ▛ ",
		"  ▝▄┴▄▘  ",
		"   ╲◆╱   ",
		"   ╰☺╯   ",
	)
	burn2.Duration, burn2.Next = 100*time.Millisecond, "burn-1"
	pop1 := styledFrame(StylePop,
		"    ✶    ",
		" ╲  │  ╱ ",
		"  ╲   ╱  ",
		" ·  ╳  · ",
		"    ☺    ",
		"   ╱ ╲   ",
		"         ",
	)
	pop1.Duration, pop1.Next = 120*time.Millisecond, "pop-2"
	pop2 := styledFrame(StylePop,
		" ·     · ",
		"   ╲ ╱   ",
		"    ✦    ",
		" ╱     ╲ ",
		"    ☺    ",
		"   ╱ ╲   ",
		"         ",
	)
	pop2.Duration, pop2.Next = 150*time.Millisecond, "pop-3"
	pop3 := styledFrame(StylePop,
		"·       ·",
		"         ",
		"  ·   ·  ",
		"         ",
		"    ☺    ",
		"   ╱ ╲   ",
		"         ",
	)
	pop3.Duration, pop3.Next = 350*time.Millisecond, ""
	for name, frame := range map[string]sprite.Frame{
		"off": off, "burn-1": burn1, "burn-2": burn2,
		"pop-1": pop1, "pop-2": pop2, "pop-3": pop3,
	} {
		if err := balloon.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := balloon.SetFrame("off"); err != nil {
		panic(err)
	}
	balloon.SetLayer(balloonLayer)
	return balloon
}

func balloonFrame(rows ...string) sprite.Frame {
	frame := sprite.FrameFromStrings(rows...)
	for _, span := range []struct{ y, min, max int }{
		{1, 2, 6},
		{2, 1, 7},
		{3, 2, 6},
	} {
		for x := span.min; x <= span.max; x++ {
			index := span.y*frame.Width + x
			if frame.Cells[index].Rune == 0 {
				style := StyleBalloonFill
				if x == span.min {
					style = StyleBalloonHighlightFill
				} else if x == span.max {
					style = StyleBalloonShadowFill
				}
				frame.Cells[index] = sprite.Cell{Rune: ' ', Style: style}
			}
		}
	}
	for index := range frame.Cells {
		switch frame.Cells[index].Rune {
		case '▲', '◆':
			frame.Cells[index].Style = StyleFlame
		case '╳':
			frame.Cells[index].Style = StylePatch
		case '☺':
			frame.Cells[index].Style = StylePilot
		case '▙', '▌', '▛', '▘':
			frame.Cells[index].Style = StyleBalloonShade
		}
	}
	return frame
}

func cloudFrame(style, highlight, shadow int, rows ...string) sprite.Frame {
	frame := sprite.FrameFromStrings(rows...)
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune == 0 {
				continue
			}
			frame.Cells[index].Style = style
			if y == 0 {
				frame.Cells[index].Style = highlight
			} else if y == frame.Height-1 {
				frame.Cells[index].Style = shadow
			}
		}
	}
	return frame
}

func (m *Model) addClouds() {
	for _, initial := range []struct {
		depth     cloudDepth
		xFraction float64
	}{
		{farCloud, 0.12},
		{farCloud, 0.58},
		{farCloud, 1.08},
		{nearCloud, 0.82},
		{nearCloud, 1.35},
	} {
		m.clouds = append(m.clouds, m.newCloud(initial.depth, float64(m.Width)*initial.xFraction))
	}
}

func (m *Model) newCloud(depth cloudDepth, x float64) cloudTrack {
	style, layer := StyleFarCloud, farCloudLayer
	speed := -(1.1 + m.random.Float64()*0.7)
	if depth == nearCloud {
		style, layer = StyleNearCloud, nearCloudLayer
		speed = -(3.6 + m.random.Float64()*1.8)
	}
	frame := newCloudFrame(style, m.random.Intn(3))
	yFraction := 0.08 + m.random.Float64()*0.78
	m.cloudID++
	cloud := sprite.New(fmt.Sprintf("cloud-%d", m.cloudID))
	if err := cloud.AddFrame("cloud", frame); err != nil {
		panic(err)
	}
	if err := cloud.SetFrame("cloud"); err != nil {
		panic(err)
	}
	cloud.SetLayer(layer)
	cloud.SetCollidable(false)
	cloud.SetPosition(x, m.cloudY(yFraction, frame.Height))
	cloud.SetVelocity(speed, 0)
	m.scene.Add(cloud)
	return cloudTrack{sprite: cloud, depth: depth, yFraction: yFraction}
}

func newCloudFrame(style, variation int) sprite.Frame {
	if style == StyleFarCloud {
		switch variation {
		case 0:
			return cloudFrame(style, StyleFarCloudHighlight, StyleFarCloudShadow,
				"   ╭──╮   ",
				"╭──╯  ╰──╮",
				"╰────────╯",
			)
		case 1:
			return cloudFrame(style, StyleFarCloudHighlight, StyleFarCloudShadow,
				"    ╭────╮    ",
				"╭───╯    ╰───╮",
				"╰────────────╯",
			)
		default:
			return cloudFrame(style, StyleFarCloudHighlight, StyleFarCloudShadow,
				"  ╭──╮    ╭╮  ",
				"╭─╯  ╰────╯╰─╮",
				"╰────────────╯",
			)
		}
	}

	switch variation {
	case 0:
		return cloudFrame(style, StyleNearCloudHighlight, StyleNearCloudShadow,
			"   ▄▄   ",
			" ▄████▄ ",
			"▀██████▀",
		)
	case 1:
		return cloudFrame(style, StyleNearCloudHighlight, StyleNearCloudShadow,
			"    ▄▄▄     ",
			" ▄▄█████▄▄  ",
			"▀██████████▀",
		)
	default:
		return cloudFrame(style, StyleNearCloudHighlight, StyleNearCloudShadow,
			"  ▄▄      ▄  ",
			"▄████▄▄▄▄██▄ ",
			"▀███████████▀",
		)
	}
}

func (m *Model) cloudY(fraction float64, cloudHeight int) float64 {
	available := max(m.Height-2-cloudHeight, 0)
	return float64(min(1+int(float64(available)*fraction), max(m.Height-cloudHeight, 0)))
}

func (m *Model) spawnStorm(x float64) {
	if m.storm != nil {
		m.scene.Remove(m.storm.cloud)
		m.scene.Remove(m.storm.bolt)
	}
	yFraction := 0.02 + m.random.Float64()*0.18
	y := m.cloudY(yFraction, 5)
	m.storm = newStormHazard(x, y, yFraction)
	m.scene.Add(m.storm.cloud)
	m.scene.Add(m.storm.bolt)
}

func (m *Model) addTerrain() {
	balloonFrame, _ := m.balloon.ActiveFrame()
	groundMax := min(5, max((m.Height-balloonFrame.Height)/2, 1))
	groundMin := min(2, groundMax)
	m.terrain = []*terrainLayer{
		newTerrainLayer("far-terrain", 11, -0.55, farTerrainLayer, 1, 9, StyleFarTerrain, StyleFarTerrain, false),
		newTerrainLayer("ground", 37, -worldScrollSpeed, groundLayer, groundMin, groundMax, StyleGroundSurface, StyleGroundFill, true),
	}
	for _, terrain := range m.terrain {
		terrain.initialize(m.scene, m.Width, m.Height)
	}
}

func (m *Model) addStageTrees() {
	for index, distance := range m.stage.TreeDistances {
		tree := newTree(index, distance)
		m.positionTree(&tree)
		m.trees = append(m.trees, tree)
		m.scene.Add(tree.sprite)
	}
}

func (m *Model) addStageFerns() {
	for index, distance := range m.stage.FernDistances {
		fern := newFern(index, distance)
		m.positionFern(&fern)
		m.ferns = append(m.ferns, fern)
		m.scene.Add(fern.sprite)
	}
}

func (m *Model) addStageGrazers() {
	for index, distance := range m.stage.GrazerDistances {
		grazer := newGrazer(index, distance)
		m.positionGrazer(&grazer)
		m.grazers = append(m.grazers, grazer)
		m.scene.Add(grazer.sprite)
	}
	m.anchorGrazers()
}

func (m *Model) positionGrazer(grazer *grazerTrack) {
	balloonX, _ := m.balloon.Position()
	x := balloonX + grazer.distance - m.distance
	grazer.sprite.SetPosition(x, float64(max(m.Height-grazerHeight-1, 0)))
}

func (m *Model) anchorGrazers() {
	if len(m.terrain) == 0 {
		return
	}
	ground := m.terrain[len(m.terrain)-1]
	for index := range m.grazers {
		frame, visible := m.grazers[index].sprite.ActiveFrame()
		if !visible {
			continue
		}
		x, _ := m.grazers[index].sprite.Position()
		surfaceY, found := 0, false
		for _, footX := range grazerFootColumns {
			candidate, candidateFound := ground.surfaceYAt(int(math.Floor(x)) + footX)
			if candidateFound && (!found || candidate > surfaceY) {
				surfaceY, found = candidate, true
			}
		}
		if found {
			m.grazers[index].sprite.SetPosition(x, float64(max(surfaceY-frame.Height+1, 0)))
		}
	}
}

func (m *Model) positionFern(fern *fernTrack) {
	balloonX, _ := m.balloon.Position()
	x := balloonX + fern.distance - m.distance
	y := float64(max(m.Height-fernHeight-1, 0))
	fern.sprite.SetPosition(x, y)
}

func (m *Model) positionTree(tree *treeTrack) {
	balloonX, _ := m.balloon.Position()
	x := balloonX + tree.distance - m.distance
	y := float64(max(m.Height-treeHeight-1, 0))
	tree.sprite.SetPosition(x, y)
}

func (m *Model) spawnPterosaur(x float64) {
	if m.pterosaur != nil {
		m.scene.Remove(m.pterosaur.sprite)
	}
	yFraction := 0.1 + m.random.Float64()*0.8
	m.pterosaur = newPterosaurHazard(x, m.pterosaurY(yFraction), yFraction)
	m.scene.Add(m.pterosaur.sprite)
}

func (m *Model) spawnRedEye(x float64) {
	if m.redEye != nil {
		m.scene.Remove(m.redEye.sprite)
	}
	yFraction := 0.1 + m.random.Float64()*0.8
	m.redEye = newProfiledPterosaurHazard(redEyePterosaurProfile, x, m.pterosaurY(yFraction), yFraction)
	m.scene.Add(m.redEye.sprite)
}

func (m *Model) pterosaurY(fraction float64) float64 {
	groundHeight := 0
	if len(m.terrain) != 0 {
		groundHeight = m.terrain[len(m.terrain)-1].maxHeight
	}
	minY := 2
	maxY := max(minY, m.Height-groundHeight-pterosaurHeight-1)
	return float64(minY) + float64(maxY-minY)*fraction
}

func (m *Model) placeBalloonOnGround() {
	if len(m.terrain) == 0 {
		return
	}
	frame, visible := m.balloon.ActiveFrame()
	if !visible {
		return
	}
	x, _ := m.balloon.Position()
	surfaceY, found := m.terrain[len(m.terrain)-1].surfaceYAt(int(math.Floor(x)) + frame.Width/2)
	if !found {
		return
	}
	m.balloon.SetPosition(x, float64(max(surfaceY-frame.Height, 0)))
	m.balloon.SetVelocity(0, 0)
	m.vy = 0
}

// Resize changes the play field and keeps the balloon on screen.
func (m *Model) Resize(width, height int) {
	firstSize := m.Width == 0 || m.Height == 0
	m.Width = max(width, 1)
	m.Height = max(height, 1)
	if firstSize {
		frame, _ := m.balloon.ActiveFrame()
		m.balloon.SetPosition(float64(max((m.Width-frame.Width)/2, 0)), float64(max((m.Height-frame.Height)/2, 0)))
	}
	m.clampBalloon()
	for _, cloud := range m.clouds {
		frame, ok := cloud.sprite.ActiveFrame()
		if !ok {
			continue
		}
		x, _ := cloud.sprite.Position()
		cloud.sprite.SetPosition(x, m.cloudY(cloud.yFraction, frame.Height))
	}
	if m.storm != nil {
		x, _ := m.storm.cloud.Position()
		m.storm.setPosition(x, m.cloudY(m.storm.yFraction, 5))
	}
	if m.pterosaur != nil {
		m.pterosaur.setBaseY(m.pterosaurY(m.pterosaur.yFraction))
	}
	if m.redEye != nil {
		m.redEye.setBaseY(m.pterosaurY(m.redEye.yFraction))
	}
	if m.finish != nil {
		m.positionFinishFlag()
	}
	for index := range m.trees {
		m.positionTree(&m.trees[index])
	}
	for _, terrain := range m.terrain {
		terrain.resize(m.Height)
	}
	for index := range m.ferns {
		m.positionFern(&m.ferns[index])
	}
	for index := range m.grazers {
		m.positionGrazer(&m.grazers[index])
	}
	m.anchorGrazers()
	if m.grounded {
		m.placeBalloonOnGround()
	}
}

func (m *Model) BalloonPosition() sprite.Point { return m.balloon.CellPosition() }

// BalloonY returns the balloon's precise vertical position in terminal cells.
func (m *Model) BalloonY() float64 {
	_, y := m.balloon.Position()
	return y
}

func (m *Model) VerticalVelocity() float64 { return m.vy }

func (m *Model) Heat() float64 { return m.heat }

func (m *Model) BurnerOn() bool { return m.burner }

func (m *Model) Popped() bool { return m.popped }

// Grounded reports whether the balloon is resting on its launch field.
func (m *Model) Grounded() bool { return m.grounded }

// Completed reports whether the balloon has crossed the stage goal.
func (m *Model) Completed() bool { return m.completed }

// Distance returns foreground distance traveled in terminal-cell units.
func (m *Model) Distance() float64 { return m.distance }

// Stage returns the current immutable stage definition.
func (m *Model) Stage() Stage { return StageAt(m.stageIndex) }

// StageIndex returns the zero-based position in the journey.
func (m *Model) StageIndex() int { return m.stageIndex }

// ContinueToStage creates the next authored stage without interrupting the
// balloon's flight. Stage-local world state is rebuilt while altitude, heat,
// vertical velocity, and burner state carry across the finish flag.
func (m *Model) ContinueToStage(stageIndex int) *Model {
	next := NewStageModel(m.Width, m.Height, stageIndex)
	x, y := m.balloon.Position()
	next.grounded = false
	next.heat = m.heat
	next.vy = m.vy
	next.balloon.SetPosition(x, y)
	next.balloon.SetVelocity(0, next.vy)
	next.SetBurner(m.burner)
	next.clampBalloon()
	return next
}

// SetBurner changes the sole flight control. Heat responds over time in Update.
func (m *Model) SetBurner(on bool) {
	if m.burner == on {
		return
	}
	m.burner = on
	frame := "off"
	if on {
		frame = "burn-1"
	}
	_ = m.balloon.SetFrame(frame)
}

// Sprites returns the layer-ordered scene snapshot used by renderers.
func (m *Model) Sprites() []*sprite.Sprite { return m.scene.Sprites() }

// TogglePause pauses or resumes simulation time.
func (m *Model) TogglePause() { m.Paused = !m.Paused }

// Update advances the deterministic simulation by one fixed step and returns
// any collisions detected after movement.
func (m *Model) Update(step time.Duration) []sprite.Collision {
	m.events = m.events[:0]
	if m.Paused || step <= 0 {
		return nil
	}
	if m.completed {
		return nil
	}
	if m.popped {
		if m.stopped {
			return nil
		}
		remaining := balloonPopDuration - m.popTime
		advance := min(step, remaining)
		m.Elapsed += advance
		m.popTime += advance
		m.scene.Update(advance)
		m.anchorGrazers()
		if m.popTime >= balloonPopDuration {
			m.stopped = true
		}
		return nil
	}
	dt := step.Seconds()
	m.updateHeat(dt)
	if m.grounded {
		if maximumLift*m.heat <= gravity {
			m.balloon.Update(step)
			return nil
		}
		m.grounded = false
		x, y := m.balloon.Position()
		m.balloon.SetPosition(x, max(y-1, 0))
		m.vy = -0.5
	}
	m.Elapsed += step
	m.distance = min(m.distance+worldScrollSpeed*step.Seconds(), m.stage.GoalDistance)
	if m.storm.updatePhase(step) {
		x, _ := m.storm.bolt.Position()
		if frame, visible := m.storm.bolt.ActiveFrame(); visible && x < float64(m.Width) && x+float64(frame.Width) > 0 {
			m.events = append(m.events, EventLightningStruck)
		}
	}
	for _, hazard := range []*pterosaurHazard{m.pterosaur, m.redEye} {
		if hazard != nil && hazard.updatePath(step) {
			x, _ := hazard.sprite.Position()
			frame, _ := hazard.sprite.ActiveFrame()
			if x < float64(m.Width) && x+float64(frame.Width) > 0 {
				m.events = append(m.events, EventPterosaurLunged)
			}
		}
	}
	m.vy = min(max(m.vy+(gravity-maximumLift*m.heat)*dt, -maximumRiseSpeed), maximumFallSpeed)
	m.balloon.SetVelocity(0, m.vy)
	collisions := m.scene.Update(step)
	m.anchorGrazers()
	for _, collision := range collisions {
		other, hitsBalloon := collisionWith(m.balloon, collision)
		if hitsBalloon && (other == m.storm.bolt || m.isPterosaurHazard(other) || m.isGroundHazard(other)) {
			m.popBalloon()
			break
		}
	}
	if !m.popped && m.distance >= m.stage.GoalDistance {
		m.completeStage()
		return collisions
	}
	for index, cloud := range m.clouds {
		frame, ok := cloud.sprite.ActiveFrame()
		if !ok {
			continue
		}
		x, _ := cloud.sprite.Position()
		if x+float64(frame.Width) < 0 {
			m.scene.Remove(cloud.sprite)
			gap := float64(10 + m.random.Intn(40))
			m.clouds[index] = m.newCloud(cloud.depth, float64(m.Width)+gap)
		}
	}
	stormX, _ := m.storm.cloud.Position()
	if stormX+stormCloudWidth < 0 {
		gap := 30 + m.random.Intn(40)
		m.spawnStorm(float64(m.Width + gap))
	}
	pterosaurX, _ := m.pterosaur.sprite.Position()
	pterosaurFrame, _ := m.pterosaur.sprite.ActiveFrame()
	if pterosaurX+float64(pterosaurFrame.Width) < 0 {
		gap := 30 + m.random.Intn(60)
		m.spawnPterosaur(float64(m.Width + gap))
	}
	if m.redEye != nil {
		redEyeX, _ := m.redEye.sprite.Position()
		redEyeFrame, _ := m.redEye.sprite.ActiveFrame()
		if redEyeX+float64(redEyeFrame.Width) < 0 {
			gap := 40 + m.random.Intn(70)
			m.spawnRedEye(float64(m.Width + gap))
		}
	}
	for _, terrain := range m.terrain {
		terrain.replenish(m.scene, m.Width, m.Height)
	}
	m.clampBalloon()
	return collisions
}

func (m *Model) updateHeat(dt float64) {
	if m.burner {
		m.heat += (1 - m.heat) * (1 - math.Exp(-dt/heatResponse.Seconds()))
	} else {
		m.heat *= math.Exp(-dt / coolingResponse.Seconds())
	}
}

func collisionWith(target *sprite.Sprite, collision sprite.Collision) (*sprite.Sprite, bool) {
	if collision.A == target {
		return collision.B, true
	}
	if collision.B == target {
		return collision.A, true
	}
	return nil, false
}

func (m *Model) isPterosaurHazard(candidate *sprite.Sprite) bool {
	return candidate == m.pterosaur.sprite || m.redEye != nil && candidate == m.redEye.sprite
}

func (m *Model) isGroundHazard(candidate *sprite.Sprite) bool {
	for _, tree := range m.trees {
		if tree.sprite == candidate {
			return true
		}
	}
	for _, grazer := range m.grazers {
		if grazer.sprite == candidate {
			return true
		}
	}
	for _, terrain := range m.terrain {
		if terrain.contains(candidate) {
			return true
		}
	}
	return false
}

func (m *Model) popBalloon() {
	if m.popped {
		return
	}
	m.popped = true
	m.events = append(m.events, EventBalloonPopped)
	m.burner = false
	m.vy = 0
	m.balloon.SetVelocity(0, 0)
	m.balloon.SetLayer(popLayer)
	m.storm.bolt.Hide()
	_ = m.balloon.SetFrame("pop-1")
}

func (m *Model) completeStage() {
	if m.completed || m.popped {
		return
	}
	m.completed = true
	if m.stageIndex+1 >= StageCount() {
		m.SetBurner(false)
		m.vy = 0
		m.balloon.SetVelocity(0, 0)
		m.finish.SetVelocity(0, 0)
		m.storm.bolt.Hide()
	}
	m.events = append(m.events, EventStageCompleted)
}

func (m *Model) positionFinishFlag() {
	balloonX, _ := m.balloon.Position()
	x := balloonX + max(m.stage.GoalDistance-m.distance, 0)
	y := float64(max(m.Height-finishHeight, 0))
	m.finish.SetPosition(x, y)
	if !m.completed {
		m.finish.SetVelocity(-worldScrollSpeed, 0)
	}
}

func (m *Model) clampBalloon() {
	frame, ok := m.balloon.ActiveFrame()
	if !ok {
		return
	}
	maxX := max(m.Width-frame.Width, 0)
	maxY := max(m.Height-frame.Height, 0)
	x, y := m.balloon.Position()
	x = min(max(x, 0), float64(maxX))
	if y < 0 {
		y = 0
		m.vy = 0
	} else if y > float64(maxY) {
		y = float64(maxY)
		m.vy = 0
	}
	m.balloon.SetPosition(x, y)
	m.balloon.SetVelocity(0, m.vy)
}
