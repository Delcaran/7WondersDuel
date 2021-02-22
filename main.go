package main

import "7WondersDuel/game"
import "7WondersDuel/gui"

func main() {
	var match game.Game
	match.CurrentAge = 1
	app := gui.Gui(&match)
	app.Run()
}
