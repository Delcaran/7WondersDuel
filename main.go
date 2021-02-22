package main

import "7WondersDuel/game"
import "7WondersDuel/gui"

func main() {
	var match game.Game
	app := gui.Gui(&match)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
