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
	redEyePterosaurHeight = 6
	pterosaurGlideTime    = 1200 * time.Millisecond
	pterosaurWingUpTime   = 140 * time.Millisecond
	pterosaurWingDownTime = 180 * time.Millisecond
	pterosaurCallTime     = 450 * time.Millisecond
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
	large       bool
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
	name: "red-eye-pterosaur", cruiseSpeed: redEyePterosaurSpeed, boostSpeed: redEyeBoostSpeed, large: true,
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
	return newProfiledPterosaurHazardWithDelay(profile, x, baseY, yFraction, 0)
}

func newProfiledPterosaurHazardWithDelay(profile pterosaurProfile, x, baseY, yFraction float64, delay time.Duration) *pterosaurHazard {
	pterosaur := sprite.New(profile.name)
	frames := ordinaryPterosaurFrames(profile.palette)
	if profile.large {
		frames = largePterosaurFrames(profile.palette)
	}
	initialFrame := "glide"
	if delay > 0 {
		delayed := frames["glide"]
		delayed.Duration, delayed.Next = delay, "glide"
		frames["delay"] = delayed
		initialFrame = "delay"
	}
	for name, frame := range frames {
		if err := pterosaur.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := pterosaur.SetFrame(initialFrame); err != nil {
		panic(err)
	}
	pterosaur.SetLayer(pterosaurLayer)
	pterosaur.SetPosition(x, baseY)
	pterosaur.SetVelocity(profile.cruiseSpeed, 0)
	return &pterosaurHazard{
		sprite: pterosaur, yFraction: yFraction, baseY: baseY, lastFrame: initialFrame, profile: profile,
	}
}

func ordinaryPterosaurFrames(palette pterosaurPalette) map[string]sprite.Frame {
	glide := pterosaurFrame(palette,
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▆▄▂   ",
		"               ",
		"               ",
		"               ",
	)
	paintPterosaur(&glide, palette.crest, 0, 3, 5)
	paintPterosaur(&glide, palette.shadow, 1, 7, 8)
	paintPterosaur(&glide, palette.wing, 1, 8, 9)
	paintPterosaur(&glide, palette.wingHighlight, 1, 10, 10)
	glide.Duration, glide.Next = pterosaurGlideTime, "up"
	up := pterosaurFrame(palette,
		"       ▗▆▄▂   ",
		"   ▗▄▟▛▘      ",
		"◂▰●██▀▀▘       ",
		"               ",
		"               ",
	)
	paintPterosaur(&up, palette.wingHighlight, 0, 7, 10)
	paintPterosaur(&up, palette.crest, 1, 3, 4)
	paintPterosaur(&up, palette.wing, 1, 5, 8)
	paintPterosaur(&up, palette.shadow, 2, 7, 7)
	up.Duration, up.Next = pterosaurWingUpTime, "down"
	down := pterosaurFrame(palette,
		"   ▗▄▖         ",
		"◂▰●██▀▀▜▅▂    ",
		"      ▜▆▄▖    ",
		"        ▝▄▂   ",
		"               ",
	)
	paintPterosaur(&down, palette.crest, 0, 3, 5)
	paintPterosaur(&down, palette.shadow, 1, 7, 8)
	paintPterosaur(&down, palette.wingHighlight, 1, 9, 10)
	paintPterosaur(&down, palette.wing, 2, 7, 10)
	paintPterosaur(&down, palette.wingHighlight, 3, 8, 10)
	down.Duration, down.Next = pterosaurWingDownTime, "call"
	call := pterosaurFrame(palette,
		"◂▰▘▗▄▖         ",
		"  ●██▀▀▜▆▄▂    ",
		"◂▰▖            ",
		"               ",
		"               ",
	)
	paintPterosaur(&call, palette.crest, 0, 3, 5)
	paintPterosaur(&call, palette.shadow, 1, 7, 8)
	paintPterosaur(&call, palette.wing, 1, 8, 9)
	paintPterosaur(&call, palette.wingHighlight, 1, 10, 10)
	call.Duration, call.Next = pterosaurCallTime, "glide"
	return map[string]sprite.Frame{"glide": glide, "up": up, "down": down, "call": call}
}

func largePterosaurFrames(palette pterosaurPalette) map[string]sprite.Frame {
	glide := pterosaurFrame(palette,
		"    ▗▄▄▖          ",
		"◂▰●████▀▀▜▆▄▂     ",
		"                  ",
		"                  ",
		"                  ",
		"                  ",
	)
	paintPterosaur(&glide, palette.crest, 0, 4, 7)
	paintPterosaur(&glide, palette.shadow, 1, 8, 9)
	paintPterosaur(&glide, palette.wing, 1, 9, 11)
	paintPterosaur(&glide, palette.wingHighlight, 1, 12, 12)
	glide.Duration, glide.Next = pterosaurGlideTime, "up"
	up := pterosaurFrame(palette,
		"        ▗▆▄▄▂     ",
		"    ▗▄▟█▛▘        ",
		"◂▰●████▀▀▘        ",
		"                  ",
		"                  ",
		"                  ",
	)
	paintPterosaur(&up, palette.wingHighlight, 0, 8, 12)
	paintPterosaur(&up, palette.crest, 1, 4, 6)
	paintPterosaur(&up, palette.wing, 1, 7, 10)
	paintPterosaur(&up, palette.shadow, 2, 8, 9)
	up.Duration, up.Next = pterosaurWingUpTime, "down"
	down := pterosaurFrame(palette,
		"    ▗▄▄▖          ",
		"◂▰●████▀▀▜▅▂      ",
		"       ▜▆▄▄▖      ",
		"         ▝▄▄▂     ",
		"                  ",
		"                  ",
	)
	paintPterosaur(&down, palette.crest, 0, 4, 7)
	paintPterosaur(&down, palette.shadow, 1, 8, 9)
	paintPterosaur(&down, palette.wingHighlight, 1, 10, 11)
	paintPterosaur(&down, palette.wing, 2, 7, 11)
	paintPterosaur(&down, palette.wingHighlight, 3, 9, 12)
	down.Duration, down.Next = pterosaurWingDownTime, "call"
	call := pterosaurFrame(palette,
		"◂▰▘ ▗▄▄▖          ",
		"  ●████▀▀▜▆▄▂     ",
		"◂▰▖               ",
		"                  ",
		"                  ",
		"                  ",
	)
	paintPterosaur(&call, palette.crest, 0, 4, 7)
	paintPterosaur(&call, palette.shadow, 1, 8, 9)
	paintPterosaur(&call, palette.wing, 1, 9, 11)
	paintPterosaur(&call, palette.wingHighlight, 1, 12, 12)
	call.Duration, call.Next = pterosaurCallTime, "glide"
	return map[string]sprite.Frame{"glide": glide, "up": up, "down": down, "call": call}
}

func pterosaurFrame(palette pterosaurPalette, rows ...string) sprite.Frame {
	frame := styledFrame(palette.body, rows...)
	for index := range frame.Cells {
		switch frame.Cells[index].Rune {
		case '●':
			frame.Cells[index].Style = palette.eye
		case '◂', '▰', '▘', '▖':
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
	if p.lastFrame == "down" && frame == "call" {
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

func hazardHeight(hazard *pterosaurHazard) int {
	return profileHeight(hazard.profile)
}

func profileHeight(profile pterosaurProfile) int {
	if profile.large {
		return redEyePterosaurHeight
	}
	return pterosaurHeight
}
