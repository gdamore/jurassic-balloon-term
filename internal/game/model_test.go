// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"testing"
	"time"
)

func makeAirborne(m *Model) {
	m.grounded = false
	x, _ := m.balloon.Position()
	frame, _ := m.balloon.ActiveFrame()
	m.balloon.SetPosition(x, float64(max((m.Height-frame.Height)/2, 0)))
	m.balloon.SetVelocity(0, 0)
	m.vy = 0
}

func TestUpdateHonorsPause(t *testing.T) {
	m := NewModel(20, 10)
	makeAirborne(m)
	m.Update(10 * time.Millisecond)
	m.TogglePause()
	m.Update(10 * time.Millisecond)

	if m.Elapsed != 10*time.Millisecond {
		t.Fatalf("elapsed = %v, want 10ms", m.Elapsed)
	}
}

func TestResizeClampsPosition(t *testing.T) {
	m := NewModel(80, 24)
	m.balloon.SetPosition(100, 100)
	m.Resize(10, 5)

	if position := m.BalloonPosition(); position.X != 1 || position.Y != 0 {
		t.Fatalf("position after resize = %v, want (1, 0)", position)
	}
}

func TestBalloonWaitsOnGroundUntilBurnerProducesLift(t *testing.T) {
	m := NewModel(80, 30)
	if !m.Grounded() {
		t.Fatal("new balloon is not grounded")
	}
	frame, _ := m.balloon.ActiveFrame()
	position := m.BalloonPosition()
	surfaceY, found := m.terrain[len(m.terrain)-1].surfaceYAt(position.X + frame.Width/2)
	if !found || position.Y+frame.Height != surfaceY {
		t.Fatalf("balloon bottom = %d, ground surface = %d, found = %t", position.Y+frame.Height, surfaceY, found)
	}

	startY := m.BalloonY()
	m.Update(time.Second)
	if m.BalloonY() != startY || m.Distance() != 0 || m.Elapsed != 0 || m.Popped() {
		t.Fatalf("idle launch state moved: y=%v distance=%v elapsed=%v popped=%t", m.BalloonY(), m.Distance(), m.Elapsed, m.Popped())
	}

	m.SetBurner(true)
	for range 30 {
		m.Update(33 * time.Millisecond)
		if !m.Grounded() {
			break
		}
	}
	if m.Grounded() || m.VerticalVelocity() >= 0 {
		t.Fatalf("burner did not launch balloon: grounded=%t vy=%v heat=%v", m.Grounded(), m.VerticalVelocity(), m.Heat())
	}
}

func TestInvincibilityIgnoresDamageButKeepsCollisionDetection(t *testing.T) {
	m := NewStageModel(80, 30, 0)
	makeAirborne(m)
	m.SetInvincible(true)
	x, y := m.balloon.Position()
	m.trees[0].sprite.SetPosition(x, y)

	collisions := m.Update(time.Millisecond)
	if len(collisions) == 0 {
		t.Fatal("invincible contact was not detected")
	}
	if m.Popped() {
		t.Fatal("invincible contact popped the balloon")
	}

	m.SetInvincible(false)
	m.Update(time.Millisecond)
	if !m.Popped() {
		t.Fatal("damage remained disabled after leaving invincible mode")
	}
}

func TestUpdateAnimatesBalloon(t *testing.T) {
	m := NewModel(20, 10)
	if got := m.balloon.FrameName(); got != "off" {
		t.Fatalf("initial frame = %q, want off", got)
	}
	m.Update(300 * time.Millisecond)
	if got := m.balloon.FrameName(); got != "off" {
		t.Fatalf("frame after update = %q, want off", got)
	}

	m.SetBurner(true)
	if got := m.balloon.FrameName(); got != "burn-1" {
		t.Fatalf("burner frame = %q, want burn-1", got)
	}
	m.Update(100 * time.Millisecond)
	if got := m.balloon.FrameName(); got != "burn-2" {
		t.Fatalf("burner frame after update = %q, want burn-2", got)
	}
	m.Update(50 * time.Millisecond)
	if got := m.balloon.FrameName(); got != "burn-3" {
		t.Fatalf("burner fade frame = %q, want burn-3", got)
	}
	frame, _ := m.balloon.ActiveFrame()
	if flame := frame.CellAt(4, 5); flame.Rune != '▼' || flame.Style != StyleFlameEmber {
		t.Fatalf("burner ember = %+v, want red-orange down flame", flame)
	}
}

func TestBalloonEnvelopeHasOpaqueFill(t *testing.T) {
	m := NewModel(40, 20)
	frame, _ := m.balloon.ActiveFrame()
	cell := frame.CellAt(4, 1)
	if cell.Rune != ' ' || cell.Style != StyleBalloonFill {
		t.Fatalf("envelope interior = %+v, want opaque fill style", cell)
	}
	bounds, _ := m.balloon.Bounds()
	if !m.balloon.CollidesAt(bounds.Min.X+4, bounds.Min.Y+1) {
		t.Fatal("opaque envelope interior is absent from balloon mask")
	}
	if cell := frame.CellAt(0, 0); cell.Rune != 0 {
		t.Fatalf("outside of envelope is opaque: %+v", cell)
	}
}

func TestBalloonEnvelopeUsesCurvedShadedCells(t *testing.T) {
	m := NewModel(40, 20)
	frame, _ := m.balloon.ActiveFrame()
	if corner := frame.CellAt(2, 0); corner.Rune != '▗' || corner.Style != StyleBalloon {
		t.Fatalf("upper-left corner = %+v, want curved outline", corner)
	}
	if highlight := frame.CellAt(2, 1); highlight.Rune != ' ' || highlight.Style != StyleBalloonHighlightFill {
		t.Fatalf("envelope highlight = %+v", highlight)
	}
	if shadow := frame.CellAt(6, 1); shadow.Rune != ' ' || shadow.Style != StyleBalloonShadowFill {
		t.Fatalf("envelope shadow = %+v", shadow)
	}
	if edge := frame.CellAt(7, 1); edge.Rune != '▙' || edge.Style != StyleBalloonShade {
		t.Fatalf("shaded edge = %+v", edge)
	}
}

func TestCloudsUseDepthAppropriateShading(t *testing.T) {
	far := newCloudFrame(StyleFarCloud, 0)
	if top := far.CellAt(3, 0); top.Style != StyleFarCloudHighlight {
		t.Fatalf("far cloud crown style = %d, want %d", top.Style, StyleFarCloudHighlight)
	}
	if bottom := far.CellAt(0, 2); bottom.Style != StyleFarCloudShadow {
		t.Fatalf("far cloud underside style = %d, want %d", bottom.Style, StyleFarCloudShadow)
	}

	near := newCloudFrame(StyleNearCloud, 0)
	if top := near.CellAt(3, 0); top.Style != StyleNearCloudHighlight {
		t.Fatalf("near cloud crown style = %d, want %d", top.Style, StyleNearCloudHighlight)
	}
	if bottom := near.CellAt(0, 2); bottom.Style != StyleNearCloudShadow {
		t.Fatalf("near cloud underside style = %d, want %d", bottom.Style, StyleNearCloudShadow)
	}
}

func TestCloudDepthAndMotion(t *testing.T) {
	m := NewModel(100, 30)
	makeAirborne(m)
	if len(m.clouds) != 5 {
		t.Fatalf("cloud count = %d, want 5", len(m.clouds))
	}

	farX, _ := m.clouds[0].sprite.Position()
	nearX, _ := m.clouds[3].sprite.Position()
	m.Update(time.Second)
	farAfter, _ := m.clouds[0].sprite.Position()
	nearAfter, _ := m.clouds[3].sprite.Position()
	if farX-farAfter >= nearX-nearAfter {
		t.Fatalf("far cloud moved %v cells, near cloud moved %v; near should move faster", farX-farAfter, nearX-nearAfter)
	}
	if m.clouds[0].sprite.Layer() >= m.balloon.Layer() || m.clouds[3].sprite.Layer() <= m.balloon.Layer() {
		t.Fatal("cloud layers do not straddle balloon layer")
	}
	for _, cloud := range m.clouds {
		if cloud.sprite.Collidable() {
			t.Fatalf("cloud %q is collidable", cloud.sprite.Name())
		}
	}
}

func TestFarCloudOutlinesCloseAtBothEdges(t *testing.T) {
	for variation := range 3 {
		frame := newCloudFrame(StyleFarCloud, variation)
		for y := 1; y < frame.Height; y++ {
			left := frame.CellAt(0, y).Rune
			right := frame.CellAt(frame.Width-1, y).Rune
			if left == 0 || right == 0 {
				t.Errorf("variation %d row %d has open edge: left=%q right=%q", variation, y, left, right)
			}
		}
	}
}

func TestCloudLeavingScreenIsReplaced(t *testing.T) {
	m := NewModel(80, 24)
	makeAirborne(m)
	old := m.clouds[0]
	frame, _ := old.sprite.ActiveFrame()
	old.sprite.SetPosition(float64(-frame.Width-1), 4)

	m.Update(33 * time.Millisecond)
	replacement := m.clouds[0]
	if replacement.sprite == old.sprite {
		t.Fatal("departed cloud was repositioned instead of replaced")
	}
	if replacement.depth != old.depth {
		t.Fatalf("replacement depth = %v, want %v", replacement.depth, old.depth)
	}
	x, _ := replacement.sprite.Position()
	if x <= float64(m.Width) {
		t.Fatalf("replacement cloud x = %v, want beyond right edge %d", x, m.Width)
	}
	for _, item := range m.Sprites() {
		if item == old.sprite {
			t.Fatal("departed cloud remains in scene")
		}
	}
}

func TestGravityAndBurnerControlVerticalMotion(t *testing.T) {
	m := NewModel(80, 40)
	makeAirborne(m)
	start := m.BalloonY()
	for range 10 {
		m.Update(33 * time.Millisecond)
	}
	if m.BalloonY() <= start || m.VerticalVelocity() <= 0 {
		t.Fatalf("gravity did not make balloon fall: y=%v start=%v vy=%v", m.BalloonY(), start, m.VerticalVelocity())
	}

	m.SetBurner(true)
	for range 60 {
		m.Update(33 * time.Millisecond)
	}
	if m.VerticalVelocity() >= 0 {
		t.Fatalf("burner did not reverse descent: vy=%v heat=%v", m.VerticalVelocity(), m.Heat())
	}
	if m.Heat() <= 0.9 {
		t.Fatalf("heat after sustained burn = %v, want > 0.9", m.Heat())
	}
}

func TestBurnerHeatDecaysAfterRelease(t *testing.T) {
	m := NewModel(80, 40)
	m.SetBurner(true)
	for range 30 {
		m.Update(33 * time.Millisecond)
	}
	heated := m.Heat()
	m.SetBurner(false)
	for range 30 {
		m.Update(33 * time.Millisecond)
	}
	if m.Heat() >= heated {
		t.Fatalf("heat after release = %v, want less than %v", m.Heat(), heated)
	}
}

func TestVerticalSpeedIsCapped(t *testing.T) {
	m := NewModel(80, 1000)
	m.scene.Remove(m.pterosaur.sprite)
	m.pterosaur = nil
	m.stage.PterosaurGroups = nil
	makeAirborne(m)
	for range 1000 {
		m.Update(33 * time.Millisecond)
	}
	if got := m.VerticalVelocity(); got != maximumFallSpeed {
		t.Fatalf("fall speed = %v, want %v", got, maximumFallSpeed)
	}

	m.SetBurner(true)
	for range 1000 {
		m.Update(33 * time.Millisecond)
	}
	if got := m.VerticalVelocity(); got != -maximumRiseSpeed {
		t.Fatalf("rise speed = %v, want %v", got, -maximumRiseSpeed)
	}
}

func TestVerticalVelocityRatioUsesAsymmetricCaps(t *testing.T) {
	m := NewModel(80, 30)
	m.vy = -maximumRiseSpeed
	if got := m.VerticalVelocityRatio(); got != -1 {
		t.Fatalf("maximum rise ratio = %v, want -1", got)
	}
	m.vy = maximumFallSpeed
	if got := m.VerticalVelocityRatio(); got != 1 {
		t.Fatalf("maximum fall ratio = %v, want 1", got)
	}
}

func TestPopAnimationCompleteOnlyAfterAnimation(t *testing.T) {
	m := NewModel(80, 30)
	m.popBalloon()
	if m.PopAnimationComplete() {
		t.Fatal("pop animation reported complete immediately")
	}
	m.Update(balloonPopDuration)
	if !m.PopAnimationComplete() {
		t.Fatal("pop animation did not report complete after its duration")
	}
}

func TestBalloonPopUsesFastLayeredExplosion(t *testing.T) {
	balloon := newBalloon()
	names := []string{"pop-1", "pop-2", "pop-3", "pop-4", "pop-5", "pop-6"}
	total := time.Duration(0)
	styles := map[int]bool{}
	for _, name := range names {
		if err := balloon.SetFrame(name); err != nil {
			t.Fatal(err)
		}
		frame, visible := balloon.ActiveFrame()
		if !visible || frame.Width != 9 || frame.Height != 7 {
			t.Fatalf("%s dimensions = %dx%d visible:%t, want 9x7", name, frame.Width, frame.Height, visible)
		}
		total += frame.Duration
		for _, cell := range frame.Cells {
			if cell.Rune != 0 && !cell.NoCollision {
				t.Fatalf("%s has collidable explosion cell %+v", name, cell)
			}
			styles[cell.Style] = true
		}
	}
	if total != balloonPopDuration {
		t.Fatalf("pop frame duration = %v, want %v", total, balloonPopDuration)
	}
	for _, style := range []int{StylePopFlash, StylePopHot, StylePop, StylePopEmber, StylePopDebris, StylePopBasket, StylePilot} {
		if !styles[style] {
			t.Errorf("pop animation does not use style %d", style)
		}
	}
}
