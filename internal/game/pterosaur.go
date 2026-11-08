// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	pterosaurLayer        = 14
	pterosaurSpeed        = -5.2
	pterosaurHeight       = 5
	pterosaurGlideTime    = 1200 * time.Millisecond
	pterosaurWingUpTime   = 140 * time.Millisecond
	pterosaurWingDownTime = 180 * time.Millisecond
	pterosaurBoostSpeed   = -18.0
	pterosaurBoostHold    = 150 * time.Millisecond
	pterosaurBoostDecay   = 700 * time.Millisecond
	pterosaurBoostTime    = pterosaurBoostHold + pterosaurBoostDecay
	pterosaurBobPeriod    = 2600 * time.Millisecond
	pterosaurBobHeight    = 1.0
)

type pterosaurHazard struct {
	sprite    *sprite.Sprite
	yFraction float64
	baseY     float64
	phase     time.Duration
	lastFrame string
	boostTime time.Duration
	lungeY    float64
}

func newPterosaurHazard(x, baseY, yFraction float64) *pterosaurHazard {
	pterosaur := sprite.New("pterosaur")
	glide := pterosaurFrame(
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▆▄▂   ",
		"               ",
		"               ",
		"               ",
	)
	paintPterosaur(&glide, StylePterosaurCrest, 0, 3, 5)
	paintPterosaur(&glide, StylePterosaurShadow, 1, 7, 8)
	paintPterosaur(&glide, StylePterosaurWing, 1, 8, 9)
	paintPterosaur(&glide, StylePterosaurWingHighlight, 1, 10, 10)
	glide.Duration, glide.Next = pterosaurGlideTime, "up"
	up := pterosaurFrame(
		"       ▗▆▄▂   ",
		"   ▗▄▟▛▘      ",
		"◂▰●██▀▀▘       ",
		"               ",
		"               ",
	)
	paintPterosaur(&up, StylePterosaurWingHighlight, 0, 7, 10)
	paintPterosaur(&up, StylePterosaurCrest, 1, 3, 4)
	paintPterosaur(&up, StylePterosaurWing, 1, 5, 8)
	paintPterosaur(&up, StylePterosaurShadow, 2, 7, 7)
	up.Duration, up.Next = pterosaurWingUpTime, "down"
	down := pterosaurFrame(
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▅▂    ",
		"      ▜▆▄▖    ",
		"        ▝▄▂   ",
		"               ",
	)
	paintPterosaur(&down, StylePterosaurCrest, 0, 3, 5)
	paintPterosaur(&down, StylePterosaurShadow, 1, 7, 8)
	paintPterosaur(&down, StylePterosaurWingHighlight, 1, 9, 10)
	paintPterosaur(&down, StylePterosaurWing, 2, 7, 10)
	paintPterosaur(&down, StylePterosaurWingHighlight, 3, 8, 10)
	down.Duration, down.Next = pterosaurWingDownTime, "glide"
	for name, frame := range map[string]sprite.Frame{
		"glide": glide, "up": up, "down": down,
	} {
		if err := pterosaur.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := pterosaur.SetFrame("glide"); err != nil {
		panic(err)
	}
	pterosaur.SetLayer(pterosaurLayer)
	pterosaur.SetPosition(x, baseY)
	pterosaur.SetVelocity(pterosaurSpeed, 0)
	return &pterosaurHazard{
		sprite: pterosaur, yFraction: yFraction, baseY: baseY, lastFrame: "glide",
	}
}

func pterosaurFrame(rows ...string) sprite.Frame {
	frame := styledFrame(StylePterosaur, rows...)
	for index := range frame.Cells {
		switch frame.Cells[index].Rune {
		case '●':
			frame.Cells[index].Style = StylePterosaurEye
		case '◂', '▰':
			frame.Cells[index].Style = StylePterosaurBeak
		}
	}
	return frame
}

func paintPterosaur(frame *sprite.Frame, style, y, x0, x1 int) {
	for x := x0; x <= x1; x++ {
		index := y*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = style
		}
	}
}

func (p *pterosaurHazard) updatePath(step time.Duration) bool {
	lunged := false
	frame := p.sprite.FrameName()
	if p.lastFrame == "down" && frame == "glide" {
		p.boostTime = pterosaurBoostTime
		_, p.lungeY = p.sprite.Position()
		lunged = true
	}
	p.lastFrame = frame
	if p.boostTime > 0 {
		p.boostTime = max(p.boostTime-step, 0)
		surge := min(float64(p.boostTime)/float64(pterosaurBoostDecay), 1)
		velocity := pterosaurSpeed + (pterosaurBoostSpeed-pterosaurSpeed)*surge
		p.sprite.SetVelocity(velocity, 0)
		x, _ := p.sprite.Position()
		p.sprite.SetPosition(x, p.lungeY)
		return lunged
	}
	p.sprite.SetVelocity(pterosaurSpeed, 0)

	p.phase = (p.phase + step) % pterosaurBobPeriod
	x, _ := p.sprite.Position()
	angle := 2 * math.Pi * float64(p.phase) / float64(pterosaurBobPeriod)
	p.sprite.SetPosition(x, p.baseY+math.Sin(angle)*pterosaurBobHeight)
	return lunged
}

func (p *pterosaurHazard) setBaseY(baseY float64) {
	p.lungeY += baseY - p.baseY
	p.baseY = baseY
	p.updatePath(0)
}
