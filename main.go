package main

import "7WondersDuel/game"
import "7WondersDuel/gui"

func main() {
	match game.Game
	match.Age = 1
	app := gui.Gui(&match)
	app.Run()
}
