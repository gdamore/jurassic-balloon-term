// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"math/rand"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	ceilingCloudLayer  = 13
	ceilingBoltLayer   = 18
	ceilingCloudRows   = 4
	ceilingBoltRows    = 5
	ceilingBoltWidth   = 3
	ceilingScrollSpeed = 1.9
	ceilingWarnTime    = 420 * time.Millisecond
	ceilingStrikeTime  = 180 * time.Millisecond
	ceilingRestTime    = 650 * time.Millisecond
)

type ceilingPhase uint8

const (
	ceilingGather ceilingPhase = iota
	ceilingWarn
	ceilingStrike
	ceilingRest
)

type ceilingStorm struct {
	cloud       *sprite.Sprite
	bolt        *sprite.Sprite
	width       int
	phase       ceilingPhase
	remaining   time.Duration
	random      *rand.Rand
	offsetX     float64
	boltX       float64
	cloudFrame  sprite.Frame
	scrollCells float64
}

func newCeilingStorm(width int, seed int64) *ceilingStorm {
	width = max(width, 1)
	random := rand.New(rand.NewSource(seed))
	cloud := sprite.New("thunderhead-ceiling")
	cloudFrame := ceilingCloudFrame(width, random)
	if err := cloud.AddFrame("cloud", cloudFrame); err != nil {
		panic(err)
	}
	if err := cloud.SetFrame("cloud"); err != nil {
		panic(err)
	}
	cloud.SetLayer(ceilingCloudLayer)

	bolt := sprite.New("ceiling-lightning")
	boltRows := []string{" ϟ ", "  ╲", " ╱ ", "↯  ", " ╲ "}
	warning := styledFrame(StyleLightningWarning, boltRows...)
	if err := bolt.AddFrame("warning", warning); err != nil {
		panic(err)
	}
	for name, frame := range fadingLightningFrames(boltRows...) {
		if err := bolt.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	bolt.SetLayer(ceilingBoltLayer)
	bolt.SetCollidable(false)
	bolt.Hide()
	bolt.SetPosition(0, ceilingCloudRows)

	storm := &ceilingStorm{cloud: cloud, bolt: bolt, width: width, random: random, cloudFrame: cloudFrame}
	storm.remaining = storm.nextGatherTime()
	return storm
}

func ceilingCloudFrame(width int, random *rand.Rand) sprite.Frame {
	frame := sprite.Frame{Width: width, Height: ceilingCloudRows, Cells: make([]sprite.Cell, width*ceilingCloudRows)}

	// The middle rows stay continuous so the ceiling remains an unmistakable
	// boundary. Rounded crown clusters and a rolling half-block underside keep
	// it from reading as a flat architectural slab.
	highlightOffset := random.Intn(12)
	shadowOffset := random.Intn(14)
	for x := range width {
		middleStyle := StyleStormCloud
		if ((x + highlightOffset) / 6 % 4) == 0 {
			middleStyle = StyleStormCloudHighlight
		}
		lowerStyle := StyleStormCloud
		if ((x + shadowOffset) / 7 % 3) == 0 {
			lowerStyle = StyleStormCloudShadow
		}
		frame.Cells[x] = sprite.Cell{Rune: '▄', Style: StyleStormCloud}
		frame.Cells[width+x] = sprite.Cell{Rune: '█', Style: middleStyle}
		frame.Cells[2*width+x] = sprite.Cell{Rune: '█', Style: lowerStyle}
	}

	for x := 0; x < width; {
		x += random.Intn(3)
		bumpWidth := 5 + random.Intn(6)
		for local := 0; local < bumpWidth && x+local < width; local++ {
			r := '█'
			if local == 0 || local == bumpWidth-1 {
				r = '▄'
			} else if local == 1 {
				r = '▟'
			} else if local == bumpWidth-2 {
				r = '▙'
			}
			frame.Cells[x+local] = sprite.Cell{Rune: r, Style: StyleStormCloudHighlight}
		}
		x += bumpWidth
	}

	undersideMotifs := [][]rune{
		{'▛', '▀', '▀', '▜', '█', '▀'},
		{'▀', '▛', '█', '▜', '▀', '▀'},
		{'▀', '▜', '█', '▀', '▛', '▀'},
	}
	for x := 0; x < width; {
		motif := undersideMotifs[random.Intn(len(undersideMotifs))]
		for _, r := range motif {
			if x >= width {
				break
			}
			frame.Cells[3*width+x] = sprite.Cell{Rune: r, Style: StyleStormCloudShadow}
			x++
		}
	}
	return frame
}

func (s *ceilingStorm) nextGatherTime() time.Duration {
	return time.Duration(700+s.random.Intn(701)) * time.Millisecond
}

func (s *ceilingStorm) update(step time.Duration) bool {
	s.scrollCloud(step)
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

func (s *ceilingStorm) scrollCloud(step time.Duration) {
	s.scrollCells += ceilingScrollSpeed * step.Seconds()
	shift := int(s.scrollCells)
	if shift == 0 || s.width == 0 {
		return
	}
	s.scrollCells -= float64(shift)
	shift %= s.width
	shifted := sprite.Frame{Width: s.cloudFrame.Width, Height: s.cloudFrame.Height, Cells: make([]sprite.Cell, len(s.cloudFrame.Cells))}
	for y := range s.cloudFrame.Height {
		for x := range s.cloudFrame.Width {
			shifted.Cells[y*shifted.Width+x] = s.cloudFrame.Cells[y*s.cloudFrame.Width+(x+shift)%s.cloudFrame.Width]
		}
	}
	s.cloudFrame = shifted
	if err := s.cloud.AddFrame("cloud", shifted); err != nil {
		panic(err)
	}
}

func (s *ceilingStorm) enterPhase() {
	switch s.phase {
	case ceilingGather:
		s.bolt.Hide()
		s.bolt.SetCollidable(false)
		s.remaining = s.nextGatherTime()
	case ceilingWarn:
		s.boltX = float64(s.random.Intn(max(s.width-ceilingBoltWidth+1, 1)))
		s.positionBolt()
		_ = s.bolt.SetFrame("warning")
		s.bolt.SetCollidable(false)
		s.remaining = ceilingWarnTime
	case ceilingStrike:
		_ = s.bolt.SetFrame("strike-bright")
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
