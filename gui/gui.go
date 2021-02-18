package gui

import (
	"7WondersDuel/game"
	"fmt"
	"github.com/rivo/tview"
)

func gui(game *game) Application {
	box := tview.NewBox().SetBorder(true).SetTitle("7 Wonders Duel")
	app := tview.NewApplication().SetRoot(box, true)
	return app
}
