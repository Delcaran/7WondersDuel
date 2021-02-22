package gui

import (
	"7WondersDuel/game"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createBoardTable(game *game.Game) *tview.Table {
	boardTable := tview.NewTable().SetBorders(false).SetSelectable(true, true)
	minX := 10000
	maxY := 0
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := game.Board.Cards[r][c]
			if card.Building != nil {
				if r > maxY {
					maxY = r
				}
				if c < minX {
					minX = c
				}
				cardName := "?"
				color := tcell.ColorWhite
				if card.Visible {
					cardName = card.Building.Name
				}
				if game.Board.CardBlocked(&card) {
					color = tcell.ColorYellow
				}
				cell := tview.NewTableCell(cardName).
					SetTextColor(color).
					SetAlign(tview.AlignCenter)
				boardTable = boardTable.SetCell(r, c, cell)
			}
		}
	}
	boardTable.Select(maxY, minX) // first "buildable" card to the left

	return boardTable
}

func Gui(game *game.Game) *tview.Application {
	// now we should populate the game?
	game.CurrentAge = 1
	game.DeployBoard()

	title := fmt.Sprintf("7 Wonders Duel - Age %d", game.CurrentAge)

	boardTable := createBoardTable(game)
	infoFrame := tview.NewFrame(nil).AddText(title, true, tview.AlignCenter, tcell.ColorGreen)
	mainGrid := tview.NewGrid().SetBorders(true)
	mainGrid.AddItem(boardTable, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(infoFrame, 0, 1, 1, 1, 0, 0, false)
	app := tview.NewApplication().SetRoot(mainGrid, true)

	// GUI done, it's time to play
	return app
}
