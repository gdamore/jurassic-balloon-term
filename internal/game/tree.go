// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import (
	"fmt"
	"time"

	"github.com/gdamore/jurassic-balloon-term/internal/sprite"
)

const (
	treeLayer  = 9
	treeHeight = 10
)

type treeTrack struct {
	sprite   *sprite.Sprite
	distance float64
	kind     treeKind
}

type treeKind uint8

const (
	treeBroadleaf treeKind = iota
	treePalm
	treeCycad
)

func newTree(index int, distance float64) treeTrack {
	kind := treeKindForIndex(index)
	name, trunkStart, leftRows, rightRows := treeArt(kind)
	tree := sprite.New(fmt.Sprintf("tree-%s-%d", name, index+1))
	left := treeFrame(kind, trunkStart, leftRows...)
	left.Duration, left.Next = 650*time.Millisecond, "right"
	right := treeFrame(kind, trunkStart, rightRows...)
	right.Duration, right.Next = 650*time.Millisecond, "left"
	for name, frame := range map[string]sprite.Frame{"left": left, "right": right} {
		if err := tree.AddFrame(name, frame); err != nil {
			panic(err)
		}
	}
	if err := tree.SetFrame("left"); err != nil {
		panic(err)
	}
	tree.SetLayer(treeLayer)
	tree.SetVelocity(-worldScrollSpeed, 0)
	return treeTrack{sprite: tree, distance: distance, kind: kind}
}

func treeKindForIndex(index int) treeKind {
	pattern := [...]treeKind{treeBroadleaf, treeBroadleaf, treePalm, treeBroadleaf, treeCycad}
	return pattern[index%len(pattern)]
}

func treeArt(kind treeKind) (string, int, []string, []string) {
	switch kind {
	case treePalm:
		return "palm", 2, []string{
				"▗▄▙▄█▄▟▄▖",
				"▝▀▜██▛▀▘ ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"         ",
			}, []string{
				"▗▄▙▄█▄▟▄▖",
				" ▝▀▜██▛▀▘",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"    █    ",
				"         ",
			}
	case treeCycad:
		return "cycad", 5, []string{
				"   ▄   ",
				"▗▄▟█▙▄▖",
				" ▀███▀ ",
				"▗▄▟█▙▄▖",
				" ▝███▘ ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"       ",
			}, []string{
				"  ▄    ",
				"▗▄▟█▙▄▖",
				" ▀███▀ ",
				"▗▄▟█▙▄▖",
				" ▝███▘ ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"  ▐█▌  ",
				"       ",
			}
	default:
		return "broadleaf", 5, []string{
				"  ▄▄   ",
				" ▟██▙  ",
				"▟████▙ ",
				"██████▌",
				"▀████▀ ",
				"  ▐▌   ",
				"  ▐▌   ",
				"  ▐▌   ",
				"  ▐▌   ",
				"       ",
			}, []string{
				"   ▄▄  ",
				"  ▟██▙ ",
				" ▟████▙",
				"▐██████",
				" ▀████▀",
				"  ▐▌   ",
				"  ▐▌   ",
				"  ▐▌   ",
				"  ▐▌   ",
				"       ",
			}
	}
}

func treeFrame(kind treeKind, trunkStart int, rows ...string) sprite.Frame {
	frame := styledFrame(StyleTreeCanopy, rows...)
	for y := trunkStart; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			index := y*frame.Width + x
			if frame.Cells[index].Rune != 0 {
				frame.Cells[index].Style = StyleTreeTrunk
			}
		}
	}
	highlights := map[treeKind][][2]int{
		treeBroadleaf: {{2, 1}, {3, 1}, {1, 2}, {2, 2}, {1, 3}},
		treePalm:      {{2, 0}, {4, 0}, {1, 1}, {2, 1}},
		treeCycad:     {{3, 0}, {1, 1}, {2, 1}, {2, 3}},
	}[kind]
	for _, point := range highlights {
		x, y := point[0], point[1]
		index := y*frame.Width + x
		if frame.Cells[index].Rune != 0 {
			frame.Cells[index].Style = StyleTreeHighlight
		}
	}
	return frame
}
