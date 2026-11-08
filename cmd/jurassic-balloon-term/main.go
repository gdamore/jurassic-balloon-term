// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package main

import (
	"fmt"
	"os"

	"github.com/gdamore/jurassic-balloon-term/internal/ui"
	"github.com/gdamore/tcell/v3"
)

func main() {
	sound, err := newGameSound()
	if err != nil {
		fmt.Fprintf(os.Stderr, "audio disabled: %v\n", err)
	} else {
		defer sound.Close()
	}

	screen, err := tcell.NewScreen(tcell.OptAdvancedKeys(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "create terminal screen: %v\n", err)
		os.Exit(1)
	}

	app := ui.New(screen)
	if sound != nil {
		app = ui.New(screen, sound)
	}
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run Jurassic Balloon Term: %v\n", err)
		os.Exit(1)
	}
}
