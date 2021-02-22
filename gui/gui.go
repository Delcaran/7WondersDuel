package gui

import (
	"7WondersDuel/game"
	"github.com/rivo/tview"
)

func Gui(game *game.Game) *tview.Application {
	game.DeployBoard()
	box := tview.NewBox().SetBorder(true).SetTitle("7 Wonders Duel")
	app := tview.NewApplication().SetRoot(box, true)
	return app
}
