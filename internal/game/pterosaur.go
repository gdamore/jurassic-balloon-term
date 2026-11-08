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
	redEyePterosaurSpeed  = 2 * pterosaurSpeed
	redEyeBoostSpeed      = -26.0
)

type pterosaurPalette struct {
	body, crest, eye, wing, wingHighlight, shadow, beak int
}

type pterosaurProfile struct {
	name        string
	cruiseSpeed float64
	boostSpeed  float64
	palette     pterosaurPalette
}

var ordinaryPterosaurProfile = pterosaurProfile{
	name: "pterosaur", cruiseSpeed: pterosaurSpeed, boostSpeed: pterosaurBoostSpeed,
	palette: pterosaurPalette{
		body: StylePterosaur, crest: StylePterosaurCrest, eye: StylePterosaurEye,
		wing: StylePterosaurWing, wingHighlight: StylePterosaurWingHighlight,
		shadow: StylePterosaurShadow, beak: StylePterosaurBeak,
	},
}

var redEyePterosaurProfile = pterosaurProfile{
	name: "red-eye-pterosaur", cruiseSpeed: redEyePterosaurSpeed, boostSpeed: redEyeBoostSpeed,
	palette: pterosaurPalette{
		body: StyleRedEyePterosaur, crest: StyleRedEyePterosaurCrest, eye: StyleRedEyePterosaurEye,
		wing: StyleRedEyePterosaurWing, wingHighlight: StyleRedEyePterosaurWingHighlight,
		shadow: StyleRedEyePterosaurShadow, beak: StyleRedEyePterosaurBeak,
	},
}

type pterosaurHazard struct {
	sprite    *sprite.Sprite
	yFraction float64
	baseY     float64
	phase     time.Duration
	lastFrame string
	boostTime time.Duration
	lungeY    float64
	profile   pterosaurProfile
}

func newPterosaurHazard(x, baseY, yFraction float64) *pterosaurHazard {
	return newProfiledPterosaurHazard(ordinaryPterosaurProfile, x, baseY, yFraction)
}

func newProfiledPterosaurHazard(profile pterosaurProfile, x, baseY, yFraction float64) *pterosaurHazard {
	pterosaur := sprite.New(profile.name)
	glide := pterosaurFrame(profile.palette,
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▆▄▂   ",
		"               ",
		"               ",
		"               ",
	)
	paintPterosaur(&glide, profile.palette.crest, 0, 3, 5)
	paintPterosaur(&glide, profile.palette.shadow, 1, 7, 8)
	paintPterosaur(&glide, profile.palette.wing, 1, 8, 9)
	paintPterosaur(&glide, profile.palette.wingHighlight, 1, 10, 10)
	glide.Duration, glide.Next = pterosaurGlideTime, "up"
	up := pterosaurFrame(profile.palette,
		"       ▗▆▄▂   ",
		"   ▗▄▟▛▘      ",
		"◂▰●██▀▀▘       ",
		"               ",
		"               ",
	)
	paintPterosaur(&up, profile.palette.wingHighlight, 0, 7, 10)
	paintPterosaur(&up, profile.palette.crest, 1, 3, 4)
	paintPterosaur(&up, profile.palette.wing, 1, 5, 8)
	paintPterosaur(&up, profile.palette.shadow, 2, 7, 7)
	up.Duration, up.Next = pterosaurWingUpTime, "down"
	down := pterosaurFrame(profile.palette,
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▅▂    ",
		"      ▜▆▄▖    ",
		"        ▝▄▂   ",
		"               ",
	)
	paintPterosaur(&down, profile.palette.crest, 0, 3, 5)
	paintPterosaur(&down, profile.palette.shadow, 1, 7, 8)
	paintPterosaur(&down, profile.palette.wingHighlight, 1, 9, 10)
	paintPterosaur(&down, profile.palette.wing, 2, 7, 10)
	paintPterosaur(&down, profile.palette.wingHighlight, 3, 8, 10)
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
	pterosaur.SetVelocity(profile.cruiseSpeed, 0)
	return &pterosaurHazard{
		sprite: pterosaur, yFraction: yFraction, baseY: baseY, lastFrame: "glide", profile: profile,
	}
}

func pterosaurFrame(palette pterosaurPalette, rows ...string) sprite.Frame {
	frame := styledFrame(palette.body, rows...)
	for index := range frame.Cells {
		switch frame.Cells[index].Rune {
		case '●':
			frame.Cells[index].Style = palette.eye
		case '◂', '▰':
			frame.Cells[index].Style = palette.beak
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
		velocity := p.profile.cruiseSpeed + (p.profile.boostSpeed-p.profile.cruiseSpeed)*surge
		p.sprite.SetVelocity(velocity, 0)
		x, _ := p.sprite.Position()
		p.sprite.SetPosition(x, p.lungeY)
		return lunged
	}
	p.sprite.SetVelocity(p.profile.cruiseSpeed, 0)

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
