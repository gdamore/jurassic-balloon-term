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
	if !a.started {
		a.drawStartScreen(width, height)
		a.screen.Show()
		return
	}
	balloonStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(238, 202, 111)).Bold(true)
	envelopeFill := color.NewRGBColor(145, 91, 55)
	balloonFillStyle := tcell.StyleDefault.Background(envelopeFill)
	balloonHighlightFillStyle := tcell.StyleDefault.Background(color.NewRGBColor(181, 119, 67))
	balloonShadowFillStyle := tcell.StyleDefault.Background(color.NewRGBColor(101, 64, 55))
	balloonShadeStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(174, 112, 59)).Bold(true)
	flameStyle := tcell.StyleDefault.Foreground(tcell.ColorOrange).Bold(true)
	patchStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(190, 52, 55)).Background(envelopeFill).Bold(true)
	pilotStyle := tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Bold(true)
	farCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(154, 186, 199))
	farCloudHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(177, 204, 211))
	farCloudShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(111, 155, 174))
	nearCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(218, 235, 235)).Bold(true)
	nearCloudHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(244, 249, 246)).Bold(true)
	nearCloudShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(169, 202, 210)).Bold(true)
	stormCloudStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(58, 62, 76)).Bold(true)
	stormCloudHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(86, 91, 108)).Bold(true)
	stormCloudShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(35, 40, 55)).Bold(true)
	rainStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(135, 203, 224))
	lightningStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 255, 32)).Bold(true)
	lightningWarningStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 255, 255)).Bold(true)
	lightningFadeStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(170, 235, 255)).Bold(true)
	lightningDimStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(65, 145, 255)).Bold(true)
	groundStrikeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorOrangeRed).Bold(true)
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
	redEyeBody := color.NewRGBColor(56, 122, 67)
	redEyePterosaurStyle := tcell.StyleDefault.Foreground(redEyeBody).Bold(true)
	redEyePterosaurCrestStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(101, 158, 75)).Bold(true)
	redEyePterosaurEyeStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 48, 45)).Background(redEyeBody).Bold(true)
	redEyePterosaurWingStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(83, 143, 91)).Bold(true)
	redEyePterosaurWingHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(137, 184, 100)).Bold(true)
	redEyePterosaurShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(38, 78, 55)).Bold(true)
	redEyePterosaurBeakStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(121, 151, 65)).Bold(true)
	finishFlagStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(247, 189, 74)).Bold(true)
	finishPoleStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(224, 231, 214)).Bold(true)
	treeCanopyStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(38, 105, 55)).Bold(true)
	treeHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(76, 145, 71)).Bold(true)
	treeTrunkStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(105, 67, 43)).Bold(true)
	fernLeafStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(45, 119, 70)).Bold(true)
	fernHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(93, 158, 80)).Bold(true)
	fernStemStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(72, 103, 50)).Bold(true)
	grazerBodyStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(126, 112, 67)).Bold(true)
	grazerHighlightStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(159, 139, 78)).Bold(true)
	grazerShadowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(88, 79, 57)).Bold(true)
	grazerHornStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(224, 205, 151)).Bold(true)
	grazerHeadStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(104, 91, 55)).Bold(true)
	grazerBeakStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(205, 184, 127)).Bold(true)
	grazerLegStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(67, 69, 47)).Bold(true)

	status := "hold space/up/W to burn  •  p pauses  •  q/esc quits"
	if a.model.Paused {
		status = "PAUSED  •  p resumes  •  q/esc quits"
	} else if a.model.Popped() {
		status = "POPPED!  •  space/R/enter retries  •  q/esc quits"
	} else if a.model.Completed() {
		if a.model.StageIndex()+1 == game.StageCount() {
			status = "JOURNEY COMPLETE!  •  space/R/enter restarts  •  q/esc quits"
		} else {
			status = "STAGE CLEAR!  •  space/R/enter continues  •  q/esc quits"
		}
	}

	drawSprites(a.screen, a.model.Sprites(), map[int]tcell.Style{
		game.StyleBalloon:                      balloonStyle,
		game.StyleBalloonFill:                  balloonFillStyle,
		game.StyleBalloonHighlightFill:         balloonHighlightFillStyle,
		game.StyleBalloonShadowFill:            balloonShadowFillStyle,
		game.StyleBalloonShade:                 balloonShadeStyle,
		game.StyleFlame:                        flameStyle,
		game.StylePatch:                        patchStyle,
		game.StylePilot:                        pilotStyle,
		game.StyleFarCloud:                     farCloudStyle,
		game.StyleFarCloudHighlight:            farCloudHighlightStyle,
		game.StyleFarCloudShadow:               farCloudShadowStyle,
		game.StyleNearCloud:                    nearCloudStyle,
		game.StyleNearCloudHighlight:           nearCloudHighlightStyle,
		game.StyleNearCloudShadow:              nearCloudShadowStyle,
		game.StyleStormCloud:                   stormCloudStyle,
		game.StyleStormCloudHighlight:          stormCloudHighlightStyle,
		game.StyleStormCloudShadow:             stormCloudShadowStyle,
		game.StyleRain:                         rainStyle,
		game.StyleLightning:                    lightningStyle,
		game.StyleLightningWarning:             lightningWarningStyle,
		game.StyleLightningFade:                lightningFadeStyle,
		game.StyleLightningDim:                 lightningDimStyle,
		game.StyleGroundStrike:                 groundStrikeStyle,
		game.StylePop:                          popStyle,
		game.StyleFarTerrain:                   farTerrainStyle,
		game.StyleGroundSurface:                groundSurfaceStyle,
		game.StyleGroundFill:                   groundFillStyle,
		game.StylePterosaur:                    pterosaurStyle,
		game.StylePterosaurCrest:               pterosaurCrestStyle,
		game.StylePterosaurEye:                 pterosaurEyeStyle,
		game.StylePterosaurWing:                pterosaurWingStyle,
		game.StylePterosaurWingHighlight:       pterosaurWingHighlightStyle,
		game.StylePterosaurShadow:              pterosaurShadowStyle,
		game.StylePterosaurBeak:                pterosaurBeakStyle,
		game.StyleRedEyePterosaur:              redEyePterosaurStyle,
		game.StyleRedEyePterosaurCrest:         redEyePterosaurCrestStyle,
		game.StyleRedEyePterosaurEye:           redEyePterosaurEyeStyle,
		game.StyleRedEyePterosaurWing:          redEyePterosaurWingStyle,
		game.StyleRedEyePterosaurWingHighlight: redEyePterosaurWingHighlightStyle,
		game.StyleRedEyePterosaurShadow:        redEyePterosaurShadowStyle,
		game.StyleRedEyePterosaurBeak:          redEyePterosaurBeakStyle,
		game.StyleFinishFlag:                   finishFlagStyle,
		game.StyleFinishPole:                   finishPoleStyle,
		game.StyleTreeCanopy:                   treeCanopyStyle,
		game.StyleTreeHighlight:                treeHighlightStyle,
		game.StyleTreeTrunk:                    treeTrunkStyle,
		game.StyleFernLeaf:                     fernLeafStyle,
		game.StyleFernHighlight:                fernHighlightStyle,
		game.StyleFernStem:                     fernStemStyle,
		game.StyleGrazerBody:                   grazerBodyStyle,
		game.StyleGrazerHighlight:              grazerHighlightStyle,
		game.StyleGrazerShadow:                 grazerShadowStyle,
		game.StyleGrazerHorn:                   grazerHornStyle,
		game.StyleGrazerHead:                   grazerHeadStyle,
		game.StyleGrazerBeak:                   grazerBeakStyle,
		game.StyleGrazerLeg:                    grazerLegStyle,
	})

	stage := a.model.Stage()
	title := fmt.Sprintf("JURASSIC BALLOON — %d/%d %s", a.model.StageIndex()+1, game.StageCount(), strings.ToUpper(stage.Name))
	putCentered(a.screen, 0, width, title, titleStyle)
	putCentered(a.screen, height-1, width, status, infoStyle)
	if !a.model.Popped() && !a.model.Completed() && width >= 34 && height >= 12 {
		flight := fmt.Sprintf("distance %3.0f/%3.0f  •  heat %3.0f%%  •  v %+0.1f", a.model.Distance(), stage.GoalDistance, a.model.Heat()*100, a.model.VerticalVelocity())
		putCentered(a.screen, height-3, width, flight, flightStyle)
	}
	a.screen.Show()
}

func (a *App) drawStartScreen(width, height int) {
	titleStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 226, 140)).Bold(true)
	artStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(238, 202, 111)).Bold(true)
	textStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(238, 245, 242))
	blurbStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(205, 229, 232))
	lines := []struct {
		text  string
		style tcell.Style
	}{
		{"JURASSIC BALLOON", titleStyle},
		{"", textStyle},
		{"", artStyle},
		{"", artStyle},
		{"", artStyle},
		{"", artStyle},
		{"", artStyle},
		{"", artStyle},
		{"", textStyle},
		{"Stranded in prehistory, you stitched your parachute", blurbStyle},
		{"into a leaky balloon. The only way home is forward.", blurbStyle},
		{"", textStyle},
		{"ENTER / SPACE / W   start", textStyle},
		{fmt.Sprintf("1–%d   test stage", game.StageCount()), textStyle},
		{fmt.Sprintf("M   music: %s", enabledLabel(a.music)), textStyle},
		{fmt.Sprintf("S   sound effects: %s", enabledLabel(a.effects)), textStyle},
		{"Q / ESC   quit", textStyle},
	}
	startY := max((height-len(lines))/2, 0)
	for index, line := range lines {
		y := startY + index
		if y >= height {
			break
		}
		putCentered(a.screen, y, width, line.text, line.style.Background(skyColor(y, height)))
	}
	a.drawStartBalloon((width-9)/2, startY+2, height)
}

func (a *App) drawStartBalloon(startX, startY, screenHeight int) {
	rows := []string{
		"  ▗▄▄▄▖  ",
		" ▟█████▙ ",
		"▐██╳████▌",
		" ▜█████▛ ",
		"  ▝▄┬▄▘  ",
		"   ╰☺╯   ",
	}
	fill := color.NewRGBColor(145, 91, 55)
	highlight := color.NewRGBColor(181, 119, 67)
	shadow := color.NewRGBColor(101, 64, 55)
	outline := color.NewRGBColor(238, 202, 111)
	patch := color.NewRGBColor(190, 52, 55)
	for y, row := range rows {
		screenY := startY + y
		if screenY < 0 || screenY >= screenHeight {
			continue
		}
		background := skyColor(screenY, screenHeight)
		for x, r := range []rune(row) {
			if r == ' ' {
				continue
			}
			style := tcell.StyleDefault.Foreground(outline).Background(background).Bold(true)
			switch r {
			case '█':
				shade := fill
				if x <= 2 {
					shade = highlight
				} else if x >= 6 {
					shade = shadow
				}
				style = tcell.StyleDefault.Foreground(shade).Background(background).Bold(true)
			case '╳':
				style = tcell.StyleDefault.Foreground(patch).Background(fill).Bold(true)
			case '▙', '▌', '▛', '▘':
				style = tcell.StyleDefault.Foreground(shadow).Background(background).Bold(true)
			case '☺':
				style = tcell.StyleDefault.Foreground(tcell.ColorLightCyan).Background(background).Bold(true)
			}
			a.screen.SetContent(startX+x, screenY, r, nil, style)
		}
	}
}

func enabledLabel(enabled bool) string {
	if enabled {
		return "ON"
	}
	return "OFF"
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
					existing, existingStyle, _ := screen.Get(bounds.Min.X+x, bounds.Min.Y+y)
					background = underlyingColor(existing, existingStyle, background)
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
