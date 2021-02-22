package gui

import (
	"7WondersDuel/game"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createBoardTable(game *game.Game) *tview.Table {
	var title string
	fmt.Sprintf(title, "Age %d", game.CurrentAge)
	boardTable := tview.NewTable().SetBorders(false)
	boardTable.Box.SetTitle(title)
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := game.Board.Cards[r][c]
			if card.Building != nil {
				cardName := "?"
				color := tcell.ColorWhite
				if card.Visible {
					cardName = card.Building.Name
				}
				if game.Board.CardBlocked(&card) {
					color = tcell.ColorYellow
				}
				cell := tview.NewTableCell(cardName).SetTextColor(color).SetAlign(tview.AlignCenter)
				boardTable = boardTable.SetCell(r, c, cell)
			}
		}
	}
	return boardTable
}

func Gui(game *game.Game) *tview.Application {
	// now we should populate the game?
	game.CurrentAge = 1
	game.DeployBoard()

	// GUI done, it's time to play
	boardTable := createBoardTable(game)
	app := tview.NewApplication().SetRoot(boardTable, true)

	return app
}
