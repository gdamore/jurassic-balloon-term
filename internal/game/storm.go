// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	stormCloudLayer = 12
	lightningLayer  = 18
	stormSpeed      = -2.8
	stormCloudWidth = 15
	lightningWidth  = 5
)

type stormPhase int

const (
	stormGather stormPhase = iota
	stormRain
	stormStrike
	stormCooldown
)

var stormPhaseDurations = [...]time.Duration{
	stormGather:   1400 * time.Millisecond,
	stormRain:     1200 * time.Millisecond,
	stormStrike:   300 * time.Millisecond,
	stormCooldown: 1700 * time.Millisecond,
}

type stormHazard struct {
	cloud        *sprite.Sprite
	bolt         *sprite.Sprite
	phase        stormPhase
	phaseElapsed time.Duration
	yFraction    float64
}

func newStormHazard(x, y float64, yFraction float64) *stormHazard {
	cloud := sprite.New("storm-cloud")
	calm := stormCloudFrame(false,
		"               ",
		"    ▄▄▄▄▄      ",
		" ▄█████████▄   ",
		"▀████████████▀ ",
		"               ",
	)
	rain1 := stormCloudFrame(true,
		"               ",
		"    ▄▄▄▄▄      ",
		" ▄█████████▄   ",
		"▀████████████▀ ",
		" ╲  │  ╲  │   ",
	)
	rain1.Duration, rain1.Next = 140*time.Millisecond, "rain-2"
	rain2 := stormCloudFrame(true,
		"               ",
		"    ▄▄▄▄▄      ",
		" ▄█████████▄   ",
		"▀████████████▀ ",
		"  │  ╲  │  ╲  ",
	)
	rain2.Duration, rain2.Next = 140*time.Millisecond, "rain-1"
	for name, frame := range map[string]sprite.Frame{
		"calm": calm, "rain-1": rain1, "rain-2": rain2,
	} {
		if err := cloud.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := cloud.SetFrame("calm"); err != nil {
		panic(err)
	}
	cloud.SetLayer(stormCloudLayer)
	cloud.SetCollidable(false)
	cloud.SetPosition(x, y)
	cloud.SetVelocity(stormSpeed, 0)

	bolt := sprite.New("lightning")
	strike := styledFrame(StyleLightning,
		" ϟ   ",
		"  ╲  ",
		"   ╲ ",
		"  ╱  ",
		" ↯   ",
	)
	if err := bolt.AddFrame("strike", strike); err != nil {
		panic(err)
	}
	bolt.Hide()
	bolt.SetLayer(lightningLayer)
	bolt.SetPosition(x+float64((stormCloudWidth-lightningWidth)/2), y+5)
	bolt.SetVelocity(stormSpeed, 0)

	return &stormHazard{cloud: cloud, bolt: bolt, phase: stormGather, yFraction: yFraction}
}

func stormCloudFrame(raining bool, rows ...string) sprite.Frame {
	frame := sprite.FrameFromStrings(rows...)
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune == 0 {
				continue
			}
			if raining && (frame.Cells[index].Rune == '╲' || frame.Cells[index].Rune == '│') {
				frame.Cells[index].Style = StyleRain
			} else if y == 1 {
				frame.Cells[index].Style = StyleStormCloudHighlight
			} else if y == 3 {
				frame.Cells[index].Style = StyleStormCloudShadow
			} else {
				frame.Cells[index].Style = StyleStormCloud
			}
		}
	}
	return frame
}

func styledFrame(style int, rows ...string) sprite.Frame {
	frame := sprite.FrameFromStrings(rows...)
	for index := range frame.Cells {
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = style
		}
	}
	return frame
}

func (s *stormHazard) updatePhase(step time.Duration) bool {
	struck := false
	s.phaseElapsed += step
	for range len(stormPhaseDurations) {
		duration := stormPhaseDurations[s.phase]
		if s.phaseElapsed < duration {
			return struck
		}
		s.phaseElapsed -= duration
		s.phase = (s.phase + 1) % stormPhase(len(stormPhaseDurations))
		struck = struck || s.phase == stormStrike
		s.applyPhase()
	}
	return struck
}

func (s *stormHazard) applyPhase() {
	switch s.phase {
	case stormGather, stormCooldown:
		_ = s.cloud.SetFrame("calm")
		s.bolt.Hide()
	case stormRain:
		_ = s.cloud.SetFrame("rain-1")
		s.bolt.Hide()
	case stormStrike:
		_ = s.cloud.SetFrame("rain-1")
		_ = s.bolt.SetFrame("strike")
	}
}

func (s *stormHazard) setPosition(x, y float64) {
	s.cloud.SetPosition(x, y)
	s.bolt.SetPosition(x+float64((stormCloudWidth-lightningWidth)/2), y+5)
}
