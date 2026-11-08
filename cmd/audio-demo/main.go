// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

// audio-demo auditions the audio engine without running the game.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/jurassic-balloon-term/internal/audio"
	"github.com/gdamore/jurassic-balloon-term/internal/sound"
)

func main() {
	engine, err := audio.New(audio.NewOtoBackend(), sound.Catalog(), audio.Options{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "start audio demo: %v\n", err)
		os.Exit(1)
	}
	defer engine.Close()

	if err := engine.Resume(); err != nil {
		fmt.Fprintf(os.Stderr, "resume audio: %v\n", err)
		os.Exit(1)
	}
	_ = engine.SetEnvironment(audio.EnvironmentPrehistoricDay)
	_ = engine.SetMusic(audio.MusicClearSkies)

	printHelp()
	scanner := bufio.NewScanner(os.Stdin)
	burner := false
	musicGain := float32(0.35)
	for {
		fmt.Print("audio> ")
		if !scanner.Scan() {
			break
		}
		command := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch command {
		case "c", "clear":
			_ = engine.SetMusic(audio.MusicClearSkies)
		case "x", "threat":
			_ = engine.SetMusic(audio.MusicThreat)
		case "f", "failure":
			_ = engine.SetMusic(audio.MusicFailure)
		case "v", "victory":
			_ = engine.SetMusic(audio.MusicVictory)
		case "a", "day":
			_ = engine.SetEnvironment(audio.EnvironmentPrehistoricDay)
		case "s", "storm":
			_ = engine.SetEnvironment(audio.EnvironmentStorm)
		case "n", "silence":
			_ = engine.SetMusic(audio.MusicNone)
			_ = engine.SetEnvironment(audio.EnvironmentNone)
		case "b", "burner":
			burner = !burner
			_ = engine.SetContinuous(audio.ContinuousBurner, burner)
			fmt.Printf("burner %s\n", onOff(burner))
		case "k", "basket", "creak":
			_ = engine.Play(audio.CueBasketCreak)
		case "p", "pop":
			_ = engine.Play(audio.CueBalloonPop)
		case "y", "pterosaur", "pterodactyl":
			_ = engine.Play(audio.CuePterosaurCall)
		case "d", "dinosaur":
			_ = engine.Play(audio.CueDistantDinosaur)
		case "r", "carnivore":
			_ = engine.Play(audio.CueDistantCarnivore)
		case "t", "thunder":
			_ = engine.Play(audio.CueThunderNear)
		case "l", "lightning", "crack":
			_ = engine.Play(audio.CueLightningCrack)
		case "+":
			musicGain = min(musicGain+0.05, 1)
			_ = engine.SetBusGain(audio.BusMusic, musicGain)
			fmt.Printf("music gain %.0f%%\n", musicGain*100)
		case "-":
			musicGain = max(musicGain-0.05, 0)
			_ = engine.SetBusGain(audio.BusMusic, musicGain)
			fmt.Printf("music gain %.0f%%\n", musicGain*100)
		case "h", "help", "?":
			printHelp()
		case "q", "quit", "exit":
			return
		case "":
			continue
		default:
			fmt.Printf("unknown command %q; enter h for help\n", command)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read command: %v\n", err)
	}
}

func printHelp() {
	fmt.Print("Jurassic Balloon audio demo (embedded CC0 music and effects)\n\n" +
		"Enter a command followed by Return:\n" +
		"  c  clear-skies music       x  approaching-threat music\n" +
		"  f  failure sting           v  victory sting\n" +
		"  a  prehistoric ambience    s  storm ambience\n" +
		"  n  music and ambience off  b  toggle burner loop\n" +
		"  k  basket creak            p  balloon pop\n" +
		"  y  pterosaur cry\n" +
		"  d  distant dinosaur        r  foreground carnivore\n" +
		"  t  nearby thunder           l  foreground lightning crack\n" +
		"  +  raise music volume      -  lower music volume\n" +
		"  h  help                    q  quit\n\n")
}

func onOff(value bool) string {
	if value {
		return "on"
	}
	return "off"
}
