// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	sauropodLayer   = 9
	sauropodHeight  = 16
	sauropodBodyY   = 10
	sauropodHerdGap = 50
)

var sauropodFootColumns = [...]int{3, 4, 15, 16}

type sauropodStyles struct {
	body, highlight, shadow int
	neck, head, eye, leg    int
}

var sauropodPalettes = [...]sauropodStyles{
	{StyleSauropodBody, StyleSauropodHighlight, StyleSauropodShadow, StyleSauropodNeck, StyleSauropodHead, StyleSauropodEye, StyleSauropodLeg},
	{StyleSauropodBodyOchre, StyleSauropodHighlightOchre, StyleSauropodShadowOchre, StyleSauropodNeckOchre, StyleSauropodHeadOchre, StyleSauropodEyeOchre, StyleSauropodLegOchre},
	{StyleSauropodBodySlate, StyleSauropodHighlightSlate, StyleSauropodShadowSlate, StyleSauropodNeckSlate, StyleSauropodHeadSlate, StyleSauropodEyeSlate, StyleSauropodLegSlate},
}

type sauropodTrack struct {
	sprite   *sprite.Sprite
	distance float64
	called   bool
	vocal    bool
	anchored bool
	motion   sauropodMotion
}

type sauropodMotion int

const (
	sauropodVertical sauropodMotion = iota
	sauropodSwivel
	sauropodDiagonal
)

func newSauropod(index int, distance float64) sauropodTrack {
	longNeck := sprite.New(fmt.Sprintf("long-neck-%d", index+1))
	motion := sauropodVertical
	frames, initial := verticalSauropodFrames(index)
	switch index % 5 {
	case 3:
		motion = sauropodSwivel
		frames, initial = swivelSauropodFrames(index)
	case 4:
		motion = sauropodDiagonal
		frames, initial = diagonalSauropodFrames(index)
	}
	for name, frame := range frames {
		if err := longNeck.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := longNeck.SetFrame(initial); err != nil {
		panic(err)
	}
	longNeck.SetLayer(sauropodLayer)
	longNeck.SetVelocity(-worldScrollSpeed, 0)
	return sauropodTrack{sprite: longNeck, distance: distance, motion: motion}
}

func verticalSauropodFrames(index int) (map[string]sprite.Frame, string) {
	low := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"▗▄▄                   ",
		"██●▙▄▄▄               ",
		"▀██████▙▄▄            ",
	)
	rising := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		" ▗▄▄                  ",
		"▐██●▌                 ",
		" ▝██▙                 ",
		"    ▜▙                ",
		"     ▐▙               ",
		"      ▜▙              ",
	)
	high := sauropodFrame(index,
		"   ▗▄▄                ",
		"  ▐██●▌               ",
		"   ▝██▘               ",
		"     ▐▌               ",
		"     ▐▌               ",
		"      █               ",
		"      ▐▌              ",
		"       █              ",
		"       ▐▌             ",
		"        █             ",
	)
	low.Duration, low.Next = 1400*time.Millisecond+time.Duration(index%3)*180*time.Millisecond, "rising"
	rising.Duration, rising.Next = 320*time.Millisecond, "high"
	high.Duration, high.Next = 1200*time.Millisecond+time.Duration((index+1)%3)*160*time.Millisecond, "lowering"
	lowering := rising
	lowering.Duration, lowering.Next = 320*time.Millisecond, "low"
	frames := map[string]sprite.Frame{
		"low": low, "rising": rising, "high": high, "lowering": lowering,
	}
	poses := [...]string{"low", "rising", "high", "lowering"}
	return frames, poses[index%len(poses)]
}

func swivelSauropodFrames(index int) (map[string]sprite.Frame, string) {
	left := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		" ▗▄▄                  ",
		"▐██●▌                 ",
		" ▝██▙                 ",
		"    ▜▙                ",
		"     ▐▙               ",
		"      ▜▙              ",
	)
	center := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"    ▗▄▄               ",
		"   ▐██●▌              ",
		"    ▝██▘              ",
		"      ▜▙              ",
		"       █              ",
		"       ▜▙             ",
	)
	right := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"       ▗▄▄            ",
		"      ▐█●██▌          ",
		"       ▝██▘           ",
		"       ▟▘             ",
		"      ▟▘              ",
		"      ▜▙              ",
	)
	left.Duration, left.Next = 2600*time.Millisecond, "turn-right"
	center.Duration, center.Next = 850*time.Millisecond, "right"
	right.Duration, right.Next = 2600*time.Millisecond, "turn-left"
	turnLeft := center
	turnLeft.Next = "left"
	initial := "left"
	if (index/5)%2 != 0 {
		initial = "right"
	}
	return map[string]sprite.Frame{
		"left": left, "turn-right": center, "right": right, "turn-left": turnLeft,
	}, initial
}

func diagonalSauropodFrames(index int) (map[string]sprite.Frame, string) {
	browse := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"▗▄▄                   ",
		"██●▙▄▄▄               ",
		"▀██████▙▄▄            ",
	)
	lifting := sauropodFrame(index,
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"                      ",
		"  ▗▄▄                 ",
		" ▐██●▌                ",
		"  ▝██▙▄▄              ",
		"      ▜██▄            ",
		"         ▜▙           ",
	)
	reach := sauropodFrame(index,
		"                      ",
		"                      ",
		"        ▗▄▄           ",
		"       ▐██●▌          ",
		"        ▝██▙          ",
		"           ▜▙         ",
		"            █         ",
		"           ▟▘         ",
		"         ▟▘           ",
		"        ▟▘            ",
	)
	browse.Duration, browse.Next = 2400*time.Millisecond, "lifting"
	lifting.Duration, lifting.Next = 700*time.Millisecond, "reach"
	reach.Duration, reach.Next = 2200*time.Millisecond, "lowering"
	lowering := lifting
	lowering.Duration, lowering.Next = 700*time.Millisecond, "browse"
	initial := "browse"
	if (index/5)%2 != 0 {
		initial = "reach"
	}
	return map[string]sprite.Frame{
		"browse": browse, "lifting": lifting, "reach": reach, "lowering": lowering,
	}, initial
}

func sauropodFrame(variant int, upper ...string) sprite.Frame {
	rows := append(append([]string(nil), upper...), sauropodBody(variant)...)
	styles := sauropodPalettes[variant%len(sauropodPalettes)]
	frame := styledFrame(styles.body, rows...)
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune == 0 {
				continue
			}
			switch {
			case frame.Cells[index].Rune == '●':
				frame.Cells[index].Style = styles.eye
			case y < sauropodBodyY:
				frame.Cells[index].Style = styles.neck
			case y == sauropodBodyY:
				frame.Cells[index].Style = styles.highlight
			case y >= frame.Height-2:
				frame.Cells[index].Style = styles.leg
			case x >= 9:
				frame.Cells[index].Style = styles.shadow
			}
		}
	}
	paintSauropodHead(&frame, styles.head)
	return frame
}

func sauropodBody(variant int) []string {
	switch variant % 5 {
	case 1:
		return []string{
			"       ▗▄██▄▖         ",
			"    ▗▄████████▄▖      ",
			" ▗▄███████████████▙▄▄ ",
			"▐███████████████████▌ ",
			" ▐██▌          ▐██▌  ",
			" ▝██            ██▘  ",
		}
	case 2:
		return []string{
			"        ▗▄▄▄▖         ",
			"   ▗▄▄████████▄▖      ",
			" ▗████████████████▄▖  ",
			"▐██████████████████▙▄▄",
			" ▐██▌          ▐██▌  ",
			" ▝██            ██▘  ",
		}
	case 3:
		return []string{
			"      ▗▄▄██▄▄▖       ",
			"   ▗▄███████████▄▖   ",
			" ▗█████████████████▙▄",
			"▐███████████████████▌",
			" ▐██▌          ▐██▌ ",
			" ▝██            ██▘ ",
		}
	case 4:
		return []string{
			"        ▗▄██▄▖        ",
			"    ▗▄██████████▄▄▖   ",
			" ▗▄████████████████▙▄▄",
			"▐███████████████████▌ ",
			" ▐██▌          ▐██▌  ",
			" ▝██            ██▘  ",
		}
	default:
		return []string{
			"       ▗▄▄██▄▖        ",
			"   ▗▄██████████▄▖     ",
			" ▗████████████████▙▄▄ ",
			"▐███████████████████▌ ",
			" ▐██▌          ▐██▌  ",
			" ▝██            ██▘  ",
		}
	}
}

func paintSauropodHead(frame *sprite.Frame, headStyle int) {
	eyeX, eyeY := -1, -1
	for y := 0; y < frame.Height && eyeX < 0; y++ {
		for x := 0; x < frame.Width; x++ {
			if frame.Cells[y*frame.Width+x].Rune == '●' {
				eyeX, eyeY = x, y
				break
			}
		}
	}
	if eyeX < 0 {
		return
	}
	for y := max(eyeY-1, 0); y <= min(eyeY+1, frame.Height-1); y++ {
		for x := max(eyeX-3, 0); x <= min(eyeX+2, frame.Width-1); x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune != 0 && frame.Cells[index].Rune != '●' {
				frame.Cells[index].Style = headStyle
			}
		}
	}
}
