// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/jurassic-balloon-term/internal/game"
	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func (a *App) draw() {
	width, height := a.screen.Size()
	for y := 0; y < height; y++ {
		a.screen.FillArea(0, y, width, 1, ' ', tcell.StyleDefault.Background(skyColor(y, height)))
	}

	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorLightYellow).Background(skyColor(0, height)).Bold(true)
	infoStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(skyColor(height-1, height))
	flightStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(skyColor(height-3, height))
	balloonStyle := tcell.StyleDefault.Foreground(tcell.ColorLightYellow)
	envelopeFill := color.NewRGBColor(91, 72, 55)
	balloonFillStyle := tcell.StyleDefault.Background(envelopeFill)
	flameStyle := tcell.StyleDefault.Foreground(tcell.ColorOrange).Bold(true)
	patchStyle := tcell.StyleDefault.Foreground(tcell.ColorMaroon).Background(envelopeFill)
	pilotStyle := tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Bold(true)
	farCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(154, 186, 199))
	nearCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(236, 245, 244)).Bold(true)
	stormCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(58, 62, 76)).Bold(true)
	rainStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(135, 203, 224))
	lightningStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	popStyle := tcell.StyleDefault.Foreground(tcell.ColorOrangeRed).Bold(true)
	farTerrainStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(43, 85, 96))
	groundSurfaceStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(91, 133, 72))
	groundFillStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(82, 61, 43))
	pterosaurBody := color.NewRGBColor(190, 78, 57)
	pterosaurStyle := tcell.StyleDefault.Foreground(pterosaurBody).Bold(true)
	pterosaurCrestStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(218, 101, 67)).Bold(true)
	pterosaurEyeStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(pterosaurBody).Bold(true)
	pterosaurWingStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(202, 123, 112)).Bold(true)
	pterosaurWingHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(230, 160, 120)).Bold(true)
	pterosaurShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(115, 63, 78)).Bold(true)
	pterosaurBeakStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(220, 103, 68)).Bold(true)

	status := "hold space/up/W to burn  •  p pauses  •  q/esc quits"
	if a.model.Paused {
		status = "PAUSED  •  p resumes  •  q/esc quits"
	} else if a.model.Popped() {
		status = "POPPED!  •  space/R/enter retries  •  q/esc quits"
	}

	drawSprites(a.screen, a.model.Sprites(), map[int]tcell.Style{
		game.StyleBalloon:                balloonStyle,
		game.StyleBalloonFill:            balloonFillStyle,
		game.StyleFlame:                  flameStyle,
		game.StylePatch:                  patchStyle,
		game.StylePilot:                  pilotStyle,
		game.StyleFarCloud:               farCloudStyle,
		game.StyleNearCloud:              nearCloudStyle,
		game.StyleStormCloud:             stormCloudStyle,
		game.StyleRain:                   rainStyle,
		game.StyleLightning:              lightningStyle,
		game.StylePop:                    popStyle,
		game.StyleFarTerrain:             farTerrainStyle,
		game.StyleGroundSurface:          groundSurfaceStyle,
		game.StyleGroundFill:             groundFillStyle,
		game.StylePterosaur:              pterosaurStyle,
		game.StylePterosaurCrest:         pterosaurCrestStyle,
		game.StylePterosaurEye:           pterosaurEyeStyle,
		game.StylePterosaurWing:          pterosaurWingStyle,
		game.StylePterosaurWingHighlight: pterosaurWingHighlightStyle,
		game.StylePterosaurShadow:        pterosaurShadowStyle,
		game.StylePterosaurBeak:          pterosaurBeakStyle,
	})

	putCentered(a.screen, 0, width, "JURASSIC BALLOON TERM", titleStyle)
	putCentered(a.screen, height-1, width, status, infoStyle)
	if !a.model.Popped() && width >= 34 && height >= 12 {
		flight := fmt.Sprintf("heat %3.0f%%  •  vertical speed %+0.1f", a.model.Heat()*100, a.model.VerticalVelocity())
		putCentered(a.screen, height-3, width, flight, flightStyle)
	}
	a.screen.Show()
}

func drawSprites(screen tcell.Screen, sprites []*sprite.Sprite, palette map[int]tcell.Style) {
	_, height := screen.Size()
	for _, item := range sprites {
		frame, visible := item.ActiveFrame()
		bounds, bounded := item.Bounds()
		if !visible || !bounded {
			continue
		}
		for y := 0; y < frame.Height; y++ {
			for x := 0; x < frame.Width; x++ {
				cell := frame.CellAt(x, y)
				if cell.Rune == 0 {
					continue
				}
				style, ok := palette[cell.Style]
				if !ok {
					style = tcell.StyleDefault
				}
				if style.GetBackground() == color.Default {
					background := skyColor(bounds.Min.Y+y, height)
					if isPartialBlock(cell.Rune) {
						existing, existingStyle, _ := screen.Get(bounds.Min.X+x, bounds.Min.Y+y)
						background = underlyingColor(existing, existingStyle, background)
					}
					style = style.Background(background)
				}
				screen.SetContent(bounds.Min.X+x, bounds.Min.Y+y, cell.Rune, nil, style)
			}
		}
	}
}

func isPartialBlock(r rune) bool {
	switch r {
	case '▀', '▄', '▌', '▐', '▖', '▗', '▘', '▝', '▚', '▞', '▙', '▛', '▜', '▟':
		return true
	default:
		return false
	}
}

func underlyingColor(existing string, style tcell.Style, fallback color.Color) color.Color {
	runes := []rune(existing)
	if len(runes) == 0 {
		return fallback
	}
	if runes[0] == ' ' {
		if background := style.GetBackground(); background != color.Default {
			return background
		}
		return fallback
	}
	if foreground := style.GetForeground(); foreground != color.Default {
		return foreground
	}
	return fallback
}

func skyColor(y, height int) color.Color {
	if height <= 1 {
		return color.NewRGBColor(20, 47, 78)
	}
	y = min(max(y, 0), height-1)
	ratio := float64(y) / float64(height-1)
	blend := func(top, bottom int32) int32 {
		return top + int32(float64(bottom-top)*ratio)
	}
	return color.NewRGBColor(
		blend(20, 72),
		blend(47, 135),
		blend(78, 158),
	)
}

func putCentered(screen tcell.Screen, y, width int, text string, style tcell.Style) {
	if y < 0 {
		return
	}
	x := max((width-len([]rune(text)))/2, 0)
	putString(screen, x, y, text, style)
}

func putString(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	width, height := screen.Size()
	if y < 0 || y >= height || x >= width {
		return
	}
	if x < 0 {
		runes := []rune(text)
		trim := min(-x, len(runes))
		text = string(runes[trim:])
		x = 0
	}
	if x < width {
		screen.PutStrStyled(x, y, strings.TrimRight(text, "\n"), style)
	}
}
