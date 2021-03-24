package gui

import (
	"7WondersDuel/game"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Print the board just for reference. All action will be made from right panel
func createBoardTable(game *game.Game) *tview.Table {
	boardTable := tview.NewTable().SetBorders(false).SetSelectable(true, true)
	boardTable.SetSelectable(false, false)
	//maxY := 0
	//minX := 10000
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := game.Board.Cards[r][c]
			//selectable := false
			if card.Building != nil {
				cardName := "XXXXXXXXXX"
				color := tcell.ColorWhite
				if card.Visible {
					//if r > maxY {
					//	maxY = r
					//}
					//if c < minX {
					//	minX = c
					//}
					cardName = card.Building.Name
					switch card.Building.Type {
					case "raw":
						color = tcell.ColorBrown
					case "manufactured":
						color = tcell.ColorDarkGrey
					case "commercial":
						color = tcell.ColorGoldenrod
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
				}
				cell := tview.NewTableCell(cardName).
					SetTextColor(color).
					SetAlign(tview.AlignCenter)
				bgColor := tcell.ColorWhite
				if !game.Board.CardBlocked(&card) {
					cell.SetBackgroundColor(bgColor)
				} else {
					if !card.Visible {
						switch game.CurrentAge {
						case 1:
							bgColor = tcell.ColorDarkGoldenrod
						case 2:
							bgColor = tcell.ColorLightBlue
						case 3:
							bgColor = tcell.ColorViolet
							if card.Building.Type == "guild" {
								bgColor = tcell.ColorPurple
							}
						default:
							bgColor = tcell.ColorWhite
						}
						cell.SetTextColor(bgColor).SetBackgroundColor(bgColor)
					}
				}
				boardTable = boardTable.SetCell(r, c, cell)
			}
			//boardTable.GetCell(r, c).SetSelectable(selectable)
		}
	}
	//boardTable.Select(maxY, minX)
	return boardTable
}

// Gui creates and returs main window ready to be displayed
func Gui(game *game.Game) *tview.Application {
	// now we should populate the game?
	game.CurrentAge = 1
	game.DeployBoard()

	title := fmt.Sprintf("7 Wonders Duel - Age %d", game.CurrentAge)

	boardTable := createBoardTable(game)

	youInfo := tview.NewFrame(nil).AddText("YOU", true, tview.AlignCenter, tcell.ColorBlue)
	opponentInfo := tview.NewFrame(nil).AddText("OPPONENT", true, tview.AlignCenter, tcell.ColorRed)
	mainLeftBottom := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainLeftBottom.AddItem(youInfo, 0, 1, false)
	mainLeftBottom.AddItem(opponentInfo, 0, 1, false)
	mainLeft := tview.NewFlex().SetDirection(tview.FlexRow)
	mainLeft.AddItem(boardTable, 0, 1, false)
	mainLeft.AddItem(mainLeftBottom, 0, 1, false)
	mainRight := tview.NewFrame(nil).AddText("Actions", true, tview.AlignCenter, tcell.ColorWhite)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainFlex.AddItem(mainLeft, 0, 1, false)
	mainFlex.AddItem(mainRight, 0, 1, true)
	main := tview.NewFrame(mainFlex).AddText(title, true, tview.AlignCenter, tcell.ColorGreen)
	app := tview.NewApplication().SetRoot(main, true)

	// GUI done, it's time to play
	return app
}
