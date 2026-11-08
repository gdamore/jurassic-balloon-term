// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	groundStormCloudWidth  = 23
	groundStormCloudHeight = 5
	groundStormBoltWidth   = 5
	groundStormImpactWidth = 7
)

type groundStormPhase uint8

const (
	groundStormGather groundStormPhase = iota
	groundStormWarn
	groundStormStrike
	groundStormRest
)

type groundStormHazard struct {
	cloud      *sprite.Sprite
	bolt       *sprite.Sprite
	impact     *sprite.Sprite
	phase      groundStormPhase
	yFraction  float64
	boltHeight int
}

func newGroundStormHazard(x, y, yFraction float64, screenHeight int) *groundStormHazard {
	cloud := sprite.New("ground-strike-thunderhead")
	calm := groundStormCloudFrame(false)
	warning := groundStormCloudFrame(true)
	warning.Duration, warning.Next = 140*time.Millisecond, "warning-2"
	warning2 := groundStormCloudFrame(true)
	warning2.Duration, warning2.Next = 140*time.Millisecond, "warning-1"
	for index := range warning2.Cells {
		if warning2.Cells[index].Style == StyleRain && warning2.Cells[index].Rune == '╲' {
			warning2.Cells[index].Rune = '│'
		} else if warning2.Cells[index].Style == StyleRain && warning2.Cells[index].Rune == '│' {
			warning2.Cells[index].Rune = '╲'
		}
	}
	for name, frame := range map[string]sprite.Frame{"calm": calm, "warning-1": warning, "warning-2": warning2} {
		if err := cloud.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := cloud.SetFrame("calm"); err != nil {
		panic(err)
	}
	cloud.SetLayer(stormCloudLayer)
	cloud.SetCollidable(false)

	boltHeight := max(screenHeight-int(y)-groundStormCloudHeight, 1)
	bolt := sprite.New("ground-strike-lightning")
	setGroundStormBoltFrames(bolt, boltHeight, boltHeight)
	bolt.SetLayer(lightningLayer)
	bolt.SetCollidable(false)
	bolt.Hide()

	impact := sprite.New("ground-strike-impact")
	flash1, flash2 := groundStormImpactFrames()
	if err := impact.AddFrame("flash-1", flash1); err != nil {
		panic(err)
	}
	if err := impact.AddFrame("flash-2", flash2); err != nil {
		panic(err)
	}
	impact.SetLayer(lightningLayer + 1)
	impact.SetCollidable(false)
	impact.Hide()

	hazard := &groundStormHazard{cloud: cloud, bolt: bolt, impact: impact, yFraction: yFraction, boltHeight: boltHeight}
	hazard.setPosition(x, y)
	hazard.impact.SetPosition(x+float64((groundStormCloudWidth-groundStormImpactWidth)/2), float64(max(screenHeight-2, 0)))
	return hazard
}

func groundStormCloudFrame(raining bool) sprite.Frame {
	rows := []string{
		"        ▄▄▄▄▄          ",
		"    ▄███████████▄      ",
		" ▄█████████████████▄   ",
		"▀█████████████████████▀",
		"                       ",
	}
	if raining {
		rows[4] = "  ╲ │ ╲ │ ╲ │ ╲ │    "
	}
	return stormCloudFrame(raining, rows...)
}

func groundStormBoltFrame(height, visibleHeight, style int) sprite.Frame {
	frame := sprite.Frame{Width: groundStormBoltWidth, Height: height, Cells: make([]sprite.Cell, groundStormBoltWidth*height)}
	for y := range min(visibleHeight, height) {
		x, glyph := 2, 'ϟ'
		switch y % 4 {
		case 1:
			x, glyph = 2, '█'
		case 2:
			x, glyph = 1, '╱'
		case 3:
			x, glyph = 1, '↯'
		}
		noCollision := style != StyleLightning
		frame.Cells[y*groundStormBoltWidth+x] = sprite.Cell{Rune: glyph, Style: style, NoCollision: noCollision}
		companionX := min(x+1, groundStormBoltWidth-1)
		companion := '▐'
		if y%4 >= 2 {
			companionX = max(x-1, 0)
			companion = '▌'
		}
		frame.Cells[y*groundStormBoltWidth+companionX] = sprite.Cell{Rune: companion, Style: style, NoCollision: noCollision}
	}
	return frame
}

func setGroundStormBoltFrames(bolt *sprite.Sprite, height, visibleHeight int) {
	visibleHeight = min(max(visibleHeight, 1), height)
	warning := groundStormBoltFrame(height, visibleHeight, StyleLightningWarning)
	bright := groundStormBoltFrame(height, visibleHeight, StyleLightning)
	bright.Duration, bright.Next = 133*time.Millisecond, "strike-fade"
	fade := groundStormBoltFrame(height, visibleHeight, StyleLightningFade)
	fade.Duration, fade.Next = 25*time.Millisecond, "strike-dim"
	dim := groundStormBoltFrame(height, visibleHeight, StyleLightningDim)
	dim.Duration = 60 * time.Millisecond
	for name, frame := range map[string]sprite.Frame{
		"warning": warning, "strike-bright": bright, "strike-fade": fade, "strike-dim": dim,
	} {
		if err := bolt.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
}

func groundStormImpactFrames() (sprite.Frame, sprite.Frame) {
	first := styledFrame(StyleGroundStrike,
		"  ╲│╱  ",
		" ▂▆✹▆▂ ",
	)
	first.Duration, first.Next = 60*time.Millisecond, "flash-2"
	second := styledFrame(StyleGroundStrike,
		" ╲ ✦ ╱ ",
		"▂▅▇█▇▅▂",
	)
	second.Duration = 80 * time.Millisecond
	return first, second
}

// updateForApproach advances the encounter according to its horizontal
// relationship with the balloon. This makes every pass an altitude gate:
// the channel gives a short preview, then flashes as it enters the balloon's
// lane instead of striking at an unrelated timer phase.
func (s *groundStormHazard) updateForApproach(balloonMaxX float64) bool {
	boltX, _ := s.bolt.Position()
	previous := s.phase
	switch s.phase {
	case groundStormGather:
		if boltX <= balloonMaxX+2 {
			s.phase = groundStormWarn
		}
	case groundStormWarn:
		if boltX <= balloonMaxX {
			s.phase = groundStormStrike
		}
	case groundStormStrike:
		if boltX < balloonMaxX-1.5 {
			s.phase = groundStormRest
		}
	}
	if s.phase != previous {
		s.applyPhase()
	}
	return previous != groundStormStrike && s.phase == groundStormStrike
}

func (s *groundStormHazard) applyPhase() {
	switch s.phase {
	case groundStormGather, groundStormRest:
		_ = s.cloud.SetFrame("calm")
		s.bolt.Hide()
		s.bolt.SetCollidable(false)
		s.impact.Hide()
	case groundStormWarn:
		_ = s.cloud.SetFrame("warning-1")
		_ = s.bolt.SetFrame("warning")
		s.bolt.SetCollidable(false)
	case groundStormStrike:
		_ = s.cloud.SetFrame("warning-1")
		_ = s.bolt.SetFrame("strike-bright")
		s.bolt.SetCollidable(true)
		_ = s.impact.SetFrame("flash-1")
	}
}

func (s *groundStormHazard) setTerrainY(surfaceY int) {
	_, boltY := s.bolt.Position()
	visibleHeight := surfaceY - int(math.Floor(boltY)) + 1
	setGroundStormBoltFrames(s.bolt, s.boltHeight, visibleHeight)
	impactX, _ := s.impact.Position()
	s.impact.SetPosition(impactX, float64(max(surfaceY-1, 0)))
}

func (s *groundStormHazard) setPosition(x, y float64) {
	s.cloud.SetPosition(x, y)
	s.bolt.SetPosition(x+float64((groundStormCloudWidth-groundStormBoltWidth)/2), y+groundStormCloudHeight)
	_, impactY := s.impact.Position()
	s.impact.SetPosition(x+float64((groundStormCloudWidth-groundStormImpactWidth)/2), impactY)
}

func (s *groundStormHazard) setSpeed(speed float64) {
	s.cloud.SetVelocity(speed, 0)
	s.bolt.SetVelocity(speed, 0)
	s.impact.SetVelocity(speed, 0)
}

func (s *groundStormHazard) reset(x, y float64) {
	s.phase = groundStormGather
	s.applyPhase()
	s.setPosition(x, y)
}
