package main

import (
	"os"

	"7WondersDuel/game"
	"7WondersDuel/gui"
	"7WondersDuel/tui"
)

func main() {
	var match game.Game

	useTui := len(os.Args[1:]) >= 1 && os.Args[1] == "tui"

	if useTui {
		app := tui.Tui(&match)
		if err := app.Run(); err != nil {
			panic(err)
		}
	} else {
		app := gui.Gui()
		(*app).Run()
	}
}
