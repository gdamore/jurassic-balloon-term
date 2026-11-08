// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math/rand"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	ceilingCloudLayer = 13
	ceilingBoltLayer  = 18
	ceilingCloudRows  = 2
	ceilingWarnTime   = 420 * time.Millisecond
	ceilingStrikeTime = 180 * time.Millisecond
	ceilingRestTime   = 650 * time.Millisecond
)

type ceilingPhase uint8

const (
	ceilingGather ceilingPhase = iota
	ceilingWarn
	ceilingStrike
	ceilingRest
)

type ceilingStorm struct {
	cloud     *sprite.Sprite
	bolt      *sprite.Sprite
	width     int
	phase     ceilingPhase
	remaining time.Duration
	random    *rand.Rand
	offsetX   float64
	boltX     float64
}

func newCeilingStorm(width int, seed int64) *ceilingStorm {
	width = max(width, 1)
	random := rand.New(rand.NewSource(seed))
	cloud := sprite.New("thunderhead-ceiling")
	if err := cloud.AddFrame("cloud", ceilingCloudFrame(width, random)); err != nil {
		panic(err)
	}
	if err := cloud.SetFrame("cloud"); err != nil {
		panic(err)
	}
	cloud.SetLayer(ceilingCloudLayer)

	bolt := sprite.New("ceiling-lightning")
	warning := styledFrame(StyleRain, "┊", " ")
	strike := styledFrame(StyleLightning, "ϟ", "↯")
	if err := bolt.AddFrame("warning", warning); err != nil {
		panic(err)
	}
	if err := bolt.AddFrame("strike", strike); err != nil {
		panic(err)
	}
	bolt.SetLayer(ceilingBoltLayer)
	bolt.SetCollidable(false)
	bolt.Hide()
	bolt.SetPosition(0, ceilingCloudRows)

	storm := &ceilingStorm{cloud: cloud, bolt: bolt, width: width, random: random}
	storm.remaining = storm.nextGatherTime()
	return storm
}

func ceilingCloudFrame(width int, random *rand.Rand) sprite.Frame {
	frame := sprite.Frame{Width: width, Height: ceilingCloudRows, Cells: make([]sprite.Cell, width*ceilingCloudRows)}
	underside := [...]rune{'▀', '█', '▛', '▜'}
	for x := range width {
		style := StyleStormCloud
		if random.Intn(5) == 0 {
			style = StyleStormCloudHighlight
		}
		frame.Cells[x] = sprite.Cell{Rune: '█', Style: style}
		frame.Cells[width+x] = sprite.Cell{Rune: underside[random.Intn(len(underside))], Style: StyleStormCloudShadow}
	}
	return frame
}

func (s *ceilingStorm) nextGatherTime() time.Duration {
	return time.Duration(700+s.random.Intn(701)) * time.Millisecond
}

func (s *ceilingStorm) update(step time.Duration) bool {
	struck := false
	s.remaining -= step
	for transitions := 0; s.remaining <= 0 && transitions < 8; transitions++ {
		overrun := -s.remaining
		s.phase = (s.phase + 1) % 4
		s.enterPhase()
		struck = struck || s.phase == ceilingStrike
		s.remaining -= overrun
	}
	return struck
}

func (s *ceilingStorm) enterPhase() {
	switch s.phase {
	case ceilingGather:
		s.bolt.Hide()
		s.bolt.SetCollidable(false)
		s.remaining = s.nextGatherTime()
	case ceilingWarn:
		s.boltX = float64(s.random.Intn(s.width))
		s.positionBolt()
		_ = s.bolt.SetFrame("warning")
		s.bolt.SetCollidable(false)
		s.remaining = ceilingWarnTime
	case ceilingStrike:
		_ = s.bolt.SetFrame("strike")
		s.bolt.SetCollidable(true)
		s.remaining = ceilingStrikeTime
	case ceilingRest:
		s.bolt.Hide()
		s.bolt.SetCollidable(false)
		s.remaining = ceilingRestTime
	}
}

func (s *ceilingStorm) setOffset(x float64) {
	s.offsetX = x
	s.cloud.SetPosition(x, 0)
	s.positionBolt()
}

func (s *ceilingStorm) positionBolt() {
	s.bolt.SetPosition(s.offsetX+s.boltX, ceilingCloudRows)
}
