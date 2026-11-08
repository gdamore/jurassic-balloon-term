// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

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
	if !a.started {
		if a.credits {
			a.drawCreditsScreen(width, height)
		} else {
			a.drawStartScreen(width, height)
		}
		a.screen.Show()
		return
	}
	balloonStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(238, 202, 111)).Bold(true)
	envelopeFill := color.NewRGBColor(145, 91, 55)
	balloonFillStyle := tcell.StyleDefault.Background(envelopeFill)
	balloonHighlightFillStyle := tcell.StyleDefault.Background(color.NewRGBColor(181, 119, 67))
	balloonShadowFillStyle := tcell.StyleDefault.Background(color.NewRGBColor(101, 64, 55))
	balloonShadeStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(174, 112, 59)).Bold(true)
	flameStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 142, 24)).Bold(true)
	flameHotStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 245, 128)).Bold(true)
	flameEmberStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(225, 62, 18)).Bold(true)
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
	popStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 116, 24)).Bold(true)
	popFlashStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 250, 205)).Bold(true)
	popHotStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 214, 42)).Bold(true)
	popEmberStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(218, 53, 24)).Bold(true)
	popDebrisStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(151, 85, 49)).Bold(true)
	popBasketStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(238, 202, 111)).Bold(true)
	farTerrainStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(43, 85, 96))
	farTerrainLavaStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(178, 55, 25)).Background(color.NewRGBColor(43, 85, 96)).Bold(true)
	farTerrainLavaGlowStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 118, 28)).Background(color.NewRGBColor(43, 85, 96)).Bold(true)
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

	drawSprites(a.screen, a.model.Sprites(), map[int]tcell.Style{
		game.StyleBalloon:                      balloonStyle,
		game.StyleBalloonFill:                  balloonFillStyle,
		game.StyleBalloonHighlightFill:         balloonHighlightFillStyle,
		game.StyleBalloonShadowFill:            balloonShadowFillStyle,
		game.StyleBalloonShade:                 balloonShadeStyle,
		game.StyleFlame:                        flameStyle,
		game.StyleFlameHot:                     flameHotStyle,
		game.StyleFlameEmber:                   flameEmberStyle,
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
		game.StylePopFlash:                     popFlashStyle,
		game.StylePopHot:                       popHotStyle,
		game.StylePopEmber:                     popEmberStyle,
		game.StylePopDebris:                    popDebrisStyle,
		game.StylePopBasket:                    popBasketStyle,
		game.StyleFarTerrain:                   farTerrainStyle,
		game.StyleFarTerrainLava:               farTerrainLavaStyle,
		game.StyleFarTerrainLavaGlow:           farTerrainLavaGlowStyle,
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
	if !a.model.Popped() && !a.model.Completed() {
		a.drawFlightGauges(width, height)
	}

	switch {
	case !a.model.Popped() && !a.model.Completed() && a.overlayAudio && a.overlayTitle != "" && time.Now().Before(a.overlayUntil):
		drawCenterOverlay(a.screen, a.overlayTitle, a.overlayDetail)
	case a.model.Paused:
		drawCenterOverlay(a.screen, "PAUSED", "P resumes")
	case a.model.PopAnimationComplete():
		drawCenterOverlay(a.screen, "POPPED!", "CLICK / SPACE / R / ENTER retries")
	case a.model.Completed() && a.model.StageIndex()+1 == game.StageCount():
		drawCenterOverlay(a.screen, "JOURNEY COMPLETE!", "CLICK / SPACE / R / ENTER restarts")
	case a.model.Grounded():
		drawCenterOverlay(a.screen, "READY TO LAUNCH", "SPACE / CLICK to fire the burner")
	case a.overlayTitle != "" && time.Now().Before(a.overlayUntil):
		drawCenterOverlay(a.screen, a.overlayTitle, a.overlayDetail)
	}
	a.screen.Show()
}

func (a *App) drawFlightGauges(width, height int) {
	if width < 20 || height < 16 {
		return
	}
	barHeight := 11
	if height < 30 {
		barHeight = 7
	}
	top := max(2, (height-barHeight)/2)
	stage := a.model.Stage()
	distance := 0.0
	if stage.GoalDistance > 0 {
		distance = a.model.Distance() / stage.GoalDistance
	}
	drawPositiveGauge(a.screen, 1, top, barHeight, 'D', distance, color.NewRGBColor(247, 189, 74))
	drawPositiveGauge(a.screen, 4, top, barHeight, 'H', a.model.Heat(), color.NewRGBColor(255, 142, 24))
	drawVelocityGauge(a.screen, 7, top, barHeight, a.model.VerticalVelocityRatio())
}

func drawPositiveGauge(screen tcell.Screen, x, top, height int, label rune, value float64, fill color.Color) {
	drawGaugeLabel(screen, x, top+height, label, fill)
	value = min(max(value, 0), 1)
	level := value * float64(height)
	full := int(level)
	partial := level - float64(full)
	for row := 0; row < height; row++ {
		y := top + height - 1 - row
		r := '┊'
		fg := color.NewRGBColor(88, 128, 136)
		if row < full {
			r, fg = '█', fill
		} else if row == full && partial > 0 {
			r, fg = fractionalBlock(partial), fill
		}
		drawHUDRune(screen, x, y, r, fg, r != '┊')
		drawHUDRune(screen, x+1, y, r, fg, r != '┊')
	}
}

func drawVelocityGauge(screen tcell.Screen, x, top, height int, value float64) {
	drawGaugeLabel(screen, x, top+height, 'V', color.NewRGBColor(182, 225, 231))
	value = min(max(value, -1), 1)
	middle := height / 2
	for row := 0; row < height; row++ {
		r := '┊'
		fg := color.NewRGBColor(88, 128, 136)
		if row == middle {
			r = '─'
		}
		drawHUDRune(screen, x, top+row, r, fg, false)
		drawHUDRune(screen, x+1, top+row, r, fg, false)
	}
	span := middle
	if value < 0 {
		filled := int(math.Ceil(-value * float64(span)))
		for row := 0; row < filled; row++ {
			y := top + middle - 1 - row
			drawHUDRune(screen, x, y, '█', color.NewRGBColor(102, 221, 235), true)
			drawHUDRune(screen, x+1, y, '█', color.NewRGBColor(102, 221, 235), true)
		}
	} else if value > 0 {
		filled := int(math.Ceil(value * float64(span)))
		for row := 0; row < filled; row++ {
			y := top + middle + 1 + row
			drawHUDRune(screen, x, y, '█', color.NewRGBColor(245, 174, 74), true)
			drawHUDRune(screen, x+1, y, '█', color.NewRGBColor(245, 174, 74), true)
		}
	}
}

func fractionalBlock(value float64) rune {
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	index := int(math.Ceil(min(max(value, 0), 1)*float64(len(blocks)))) - 1
	if index < 0 {
		return ' '
	}
	return blocks[index]
}

func drawGaugeLabel(screen tcell.Screen, x, y int, label rune, foreground color.Color) {
	drawHUDRune(screen, x, y, label, foreground, true)
}

func drawHUDRune(screen tcell.Screen, x, y int, r rune, foreground color.Color, bold bool) {
	_, height := screen.Size()
	existing, existingStyle, _ := screen.Get(x, y)
	background := underlyingColor(existing, existingStyle, skyColor(y, height))
	style := tcell.StyleDefault.Foreground(foreground).Background(background)
	if bold {
		style = style.Bold(true)
	}
	screen.SetContent(x, y, r, nil, style)
}

func drawCenterOverlay(screen tcell.Screen, title, detail string) {
	_, height := screen.Size()
	drawOverlayAt(screen, max((height-4)/2, 1), title, detail)
}

func drawOverlayAt(screen tcell.Screen, startY int, title, detail string) {
	width, height := screen.Size()
	boxWidth := max(len([]rune(title)), len([]rune(detail))) + 4
	boxWidth = min(boxWidth, width)
	boxHeight := 4
	startX := max((width-boxWidth)/2, 0)
	startY = min(max(startY, 0), max(height-boxHeight, 0))
	background := color.NewRGBColor(26, 45, 57)
	boxStyle := tcell.StyleDefault.Background(background)
	for y := startY; y < min(startY+boxHeight, height); y++ {
		screen.FillArea(startX, y, boxWidth, 1, ' ', boxStyle)
	}
	titleStyle := boxStyle.Foreground(color.NewRGBColor(255, 226, 140)).Bold(true)
	detailStyle := boxStyle.Foreground(color.NewRGBColor(220, 236, 236))
	putCenteredIn(screen, startX, startY+1, boxWidth, title, titleStyle)
	putCenteredIn(screen, startX, startY+2, boxWidth, detail, detailStyle)
}

func putCenteredIn(screen tcell.Screen, x, y, width int, text string, style tcell.Style) {
	textWidth := len([]rune(text))
	putString(screen, x+max((width-textWidth)/2, 0), y, text, style)
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
		{"SPACE / CLICK   start", textStyle},
		{fmt.Sprintf("1–%d   test stage", game.StageCount()), textStyle},
		{fmt.Sprintf("M   music: %s", enabledLabel(a.music)), textStyle},
		{fmt.Sprintf("S   sound: %s", enabledLabel(a.effects)), textStyle},
		{fmt.Sprintf("− / +   volume: %s", volumeLabel(a.volume)), textStyle},
		{"C   credits", textStyle},
		{"Q / ESC   quit", textStyle},
	}
	startY := max((height-len(lines))/2, 0)
	a.drawTitleScenery(width, height, startY, startY+len(lines))
	for index, line := range lines {
		y := startY + index
		if y >= height {
			break
		}
		putCentered(a.screen, y, width, line.text, line.style.Background(skyColor(y, height)))
	}
	a.drawStartBalloon((width-9)/2, startY+2, height)
}

func (a *App) drawCreditsScreen(width, height int) {
	titleStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(255, 226, 140)).Bold(true)
	headingStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(182, 225, 231)).Bold(true)
	textStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(220, 236, 236))
	dimStyle := tcell.StyleDefault.Foreground(color.NewRGBColor(167, 202, 207))
	lines := []struct {
		text  string
		style tcell.Style
	}{
		{"CREDITS", titleStyle},
		{"", textStyle},
		{"JURASSIC BALLOON", headingStyle},
		{"Created, designed, and programmed by", textStyle},
		{"Garrett D'Amore", titleStyle},
		{"", textStyle},
		{"REQUIRED AUDIO ATTRIBUTION", headingStyle},
		{"“Soliloquy” — Matthew Pablo · CC BY 3.0 · edited", textStyle},
		{"“R4_00491 Jurassic World Alien Planet”", textStyle},
		{"kevp888 / Kevin Luce · CC BY 4.0 · excerpted and edited", dimStyle},
		{"“Lightning and Thunder Clap” — akelley6", textStyle},
		{"CC BY 4.0 · excerpted and edited", dimStyle},
		{"", textStyle},
		{"ADDITIONAL AUDIO THANKS", headingStyle},
		{"Sorth · Kistol · congusbongus · kindland · AntumDeluge", textStyle},
		{"rubberduck · Alessandrobraz · SketchMan3", textStyle},
		{"", textStyle},
		{"BUILT WITH", headingStyle},
		{"Tcell · Oto · oggvorbis", textStyle},
		{"Inspired by Escape from Proxima 5", dimStyle},
		{"", textStyle},
		{"Full licenses and source links: THIRD_PARTY_NOTICES.md", dimStyle},
		{"", textStyle},
		{"C / Q / ESC   return", titleStyle},
	}
	startY := max((height-len(lines))/2, 0)
	a.drawTitleScenery(width, height, startY, startY+len(lines))
	for index, line := range lines {
		y := startY + index
		if y >= height {
			break
		}
		putCentered(a.screen, y, width, line.text, line.style.Background(skyColor(y, height)))
	}
}

type titleCloud struct {
	rows      []string
	xFraction float64
	y         int
	speed     float64
	highlight color.Color
	body      color.Color
	shadow    color.Color
}

func (a *App) drawTitleScenery(width, height, contentTop, contentBottom int) {
	if width < 30 || height < 28 {
		return
	}
	farHighlight := color.NewRGBColor(177, 204, 211)
	farBody := color.NewRGBColor(154, 186, 199)
	farShadow := color.NewRGBColor(111, 155, 174)
	nearHighlight := color.NewRGBColor(244, 249, 246)
	nearBody := color.NewRGBColor(218, 235, 235)
	nearShadow := color.NewRGBColor(169, 202, 210)
	clouds := []titleCloud{
		{
			rows: []string{
				"   ╭──╮   ",
				"╭──╯  ╰──╮",
				"╰────────╯",
			},
			xFraction: 0.12, y: 2, speed: 0.8,
			highlight: farHighlight, body: farBody, shadow: farShadow,
		},
		{
			rows: []string{
				"    ▄▄▄     ",
				" ▄▄█████▄▄  ",
				"▀██████████▀",
			},
			xFraction: 0.70, y: 5, speed: 1.25,
			highlight: nearHighlight, body: nearBody, shadow: nearShadow,
		},
		{
			rows: []string{
				"  ╭──╮    ╭╮  ",
				"╭─╯  ╰────╯╰─╮",
				"╰────────────╯",
			},
			xFraction: 1.10, y: 2, speed: 0.65,
			highlight: farHighlight, body: farBody, shadow: farShadow,
		},
	}
	for _, cloud := range clouds {
		if cloud.y+len(cloud.rows) >= contentTop {
			continue
		}
		cloudWidth := len([]rune(cloud.rows[len(cloud.rows)-1]))
		x := titleCloudX(width, cloudWidth, cloud.xFraction, cloud.speed, a.titleTime)
		drawTitleCloud(a.screen, x, cloud.y, cloud, height)
	}

	terrainRoom := height - contentBottom
	farMax := min(9, terrainRoom-1)
	if farMax >= 2 {
		drawTitleTerrain(a.screen, height, a.titleTime, 0.45, min(5, farMax), farMax,
			color.NewRGBColor(51, 92, 99), color.NewRGBColor(43, 78, 87), 0.08)
	}
	groundMax := min(5, terrainRoom-1)
	if groundMax >= 2 {
		drawTitleTerrain(a.screen, height, a.titleTime, 1.05, 2, groundMax,
			color.NewRGBColor(91, 133, 72), color.NewRGBColor(82, 61, 43), 1.6)
	}
}

func titleCloudX(width, cloudWidth int, initialFraction, speed float64, elapsed time.Duration) int {
	gap := 18.0
	period := float64(width+cloudWidth) + gap
	x := initialFraction*float64(width) - speed*elapsed.Seconds()
	for x < -float64(cloudWidth) {
		x += period
	}
	for x >= float64(width)+gap {
		x -= period
	}
	return int(math.Floor(x))
}

func drawTitleCloud(screen tcell.Screen, startX, startY int, cloud titleCloud, screenHeight int) {
	for y, row := range cloud.rows {
		foreground := cloud.body
		if y == 0 {
			foreground = cloud.highlight
		} else if y == len(cloud.rows)-1 {
			foreground = cloud.shadow
		}
		for x, r := range []rune(row) {
			screenX, screenY := startX+x, startY+y
			if r == ' ' || screenX < 0 {
				continue
			}
			width, _ := screen.Size()
			if screenX >= width {
				continue
			}
			style := tcell.StyleDefault.Foreground(foreground).Background(skyColor(screenY, screenHeight))
			screen.SetContent(screenX, screenY, r, nil, style)
		}
	}
}

func drawTitleTerrain(screen tcell.Screen, screenHeight int, elapsed time.Duration, speed float64, minHeight, maxHeight int, surface, fill color.Color, phase float64) {
	width, _ := screen.Size()
	offset := speed * elapsed.Seconds()
	previous := titleTerrainHeight(-1, offset, minHeight, maxHeight, phase)
	for x := 0; x < width; x++ {
		columnHeight := titleTerrainHeight(x, offset, minHeight, maxHeight, phase)
		top := screenHeight - columnHeight
		surfaceRune := '▄'
		switch {
		case columnHeight > previous:
			surfaceRune = '▟'
		case columnHeight < previous:
			surfaceRune = '▙'
		}
		existing, existingStyle, _ := screen.Get(x, top)
		background := underlyingColor(existing, existingStyle, skyColor(top, screenHeight))
		screen.SetContent(x, top, surfaceRune, nil, tcell.StyleDefault.Foreground(surface).Background(background))
		for y := top + 1; y < screenHeight; y++ {
			screen.SetContent(x, y, '█', nil, tcell.StyleDefault.Foreground(fill).Background(fill))
		}
		previous = columnHeight
	}
}

func titleTerrainHeight(x int, offset float64, minHeight, maxHeight int, phase float64) int {
	worldX := float64(x) + offset
	wave := math.Sin(worldX*0.075+phase) + 0.45*math.Sin(worldX*0.19+phase*1.7)
	ratio := min(max((wave+1.45)/2.9, 0), 1)
	return minHeight + int(math.Round(ratio*float64(maxHeight-minHeight)))
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
