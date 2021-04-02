package gui

import (
	"7WondersDuel/game"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Print the board just for reference. All action will be made from right panel
func createBoardTable(game *game.Game) *tview.Table {
	boardTable := tview.NewTable().SetBorders(false).SetSelectable(false, false)
	boardTable = boardTable.Clear()
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

func fillPlayerInfoArea(g *game.Game, player int, view *tview.Table) {
	view.SetBorders(false).SetSelectable(false, false)
	view.Clear()
	fixedRes, _ := g.Players[player].AvailableResources()
	cText := 0
	cVal := 1
	data := fixedRes.ToMap()
	labels := []string{"Coins", "Wood", "Stone", "Clay", "Glass", "Papyrus"}
	for riga, label := range labels {
		value := data[label]
		color := tcell.ColorWhite
		switch label {
		case "Coins":
			color = tcell.ColorYellow
		case "Wood":
			color = tcell.ColorBrown
		case "Stone":
			color = tcell.ColorGrey
		case "Clay":
			color = tcell.ColorOrange
		case "Glass":
			color = tcell.ColorLightBlue
		case "Papyrus":
			color = tcell.ColorGoldenrod
		}
		view = view.SetCell(riga, cText, tview.NewTableCell(label).SetTextColor(color).SetAlign(tview.AlignRight))
		view = view.SetCell(riga, cVal, tview.NewTableCell(strconv.Itoa(value)).SetTextColor(color).SetAlign(tview.AlignCenter))
	}
}

// Gui creates and returs main window ready to be displayed
func Gui(game *game.Game) *tview.Application {
	// create components & layout
	app := tview.NewApplication()
	youInfo := tview.NewTable()
	fillPlayerInfoArea(game, 0, youInfo)
	youInfoFrame := tview.NewFrame(youInfo).AddText("YOU", true, tview.AlignCenter, tcell.ColorBlue)
	opponentInfoFrame := tview.NewFrame(nil).AddText("OPPONENT", true, tview.AlignCenter, tcell.ColorRed)
	mainLeftBottom := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainLeftBottom.AddItem(youInfoFrame, 0, 1, false)
	mainLeftBottom.AddItem(opponentInfoFrame, 0, 1, false)
	mainLeft := tview.NewFlex().SetDirection(tview.FlexRow)
	boardTable := tview.NewTable()
	mainLeft.AddItem(boardTable, 0, 1, false)
	mainLeft.AddItem(mainLeftBottom, 0, 1, false)
	mainRight := tview.NewFrame(nil).AddText("Actions", true, tview.AlignCenter, tcell.ColorWhite)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainFlex.AddItem(mainLeft, 0, 1, false)
	mainFlex.AddItem(mainRight, 0, 1, true) // "Actions" has focus because all commands are here
	main := tview.NewFrame(mainFlex)

	refreshFunc := func() {
		game.CurrentAge++
		if game.CurrentAge > 3 {
			app.Stop()
		} else {
			game.DeployBoard()
			title := fmt.Sprintf("7 Wonders Duel - Age %d", game.CurrentAge)
			main.Clear()
			titleColor := tcell.ColorWhite
			switch game.CurrentAge {
			case 1:
				titleColor = tcell.ColorDarkGoldenrod
			case 2:
				titleColor = tcell.ColorLightBlue
			case 3:
				titleColor = tcell.ColorViolet
			default:
				titleColor = tcell.ColorRed // means "oh shit"
			}
			main.AddText(title, true, tview.AlignCenter, titleColor)
			mainLeft.RemoveItem(boardTable)
			mainLeft.RemoveItem(mainLeftBottom)
			boardTable = createBoardTable(game)
			mainLeft.AddItem(boardTable, 0, 1, false)
			mainLeft.AddItem(mainLeftBottom, 0, 1, false)
		}
	}
	app.SetRoot(main, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// BEGIN DEBUG FUNCTIONALITY
		switch event.Rune() {
		case 'n':
			refreshFunc()
		case 'q':
			app.Stop()
		}
		// END DEBUG FUNCTIONALITY
		return event
	})

	// GUI done, it's time to play

	// initialization
	game.CurrentAge = 0
	refreshFunc()

	return app
}
