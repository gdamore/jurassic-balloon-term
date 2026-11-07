// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package main

import (
	"fmt"
	"os"

	"github.com/gdamore/jurassic-balloon-term/internal/ui"
	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create terminal screen: %v\n", err)
		os.Exit(1)
	}

	if err := ui.New(screen).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run Jurassic Balloon Term: %v\n", err)
		os.Exit(1)
	}
}
