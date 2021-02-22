package gui

import (
	"7WondersDuel/game"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Print the board just for reference. All action will be made from right panel
func createBoardTable(game *game.Game) *tview.Table {
	boardTable := tview.NewTable().SetBorders(false).SetSelectable(false, false)
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := game.Board.Cards[r][c]
			if card.Building != nil {
				cardName := "?"
				color := tcell.ColorWhite
				if card.Visible {
					cardName = card.Building.Name
					switch card.Building.Type {
					case "raw":
						color = tcell.ColorBrown
					case "manufactured":
						color = tcell.ColorGrey
					case "commercial":
						color = tcell.ColorYellow
					case "military":
						color = tcell.ColorRed
					case "guild":
						color = tcell.ColorPurple
					case "civilian":
						color = tcell.ColorBlue
					case "scientific":
						color = tcell.ColorGreen
					default:
						color = tcell.ColorWhite
					}
					if game.Board.CardBlocked(&card) {
						// TODO colore invertito
					}
				}
				cell := tview.NewTableCell(cardName).
					SetTextColor(color).
					SetAlign(tview.AlignCenter)
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
