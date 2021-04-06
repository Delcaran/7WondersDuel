package gui

import (
	"7WondersDuel/game"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getTypeColor(cardtype string) tcell.Color {
	color := tcell.ColorWhite
	switch cardtype {
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
	return color
}

// Print the board just for reference. All action will be made from right panel
func fillBoardTable(game *game.Game, boardTable *tview.Table) {
	boardTable = boardTable.Clear().SetBorders(false).SetSelectable(false, false)
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := game.Board.Cards[r][c]
			if card.Building != nil {
				cardName := "XXXXXXXXXX"
				color := tcell.ColorWhite
				if card.Visible {
					cardName = card.Building.Name
					color = getTypeColor(card.Building.Type)
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
		}
	}
}

func fillPlayerInfoArea(g *game.Game, player int, view *tview.Table, frame *tview.Frame) {
	var text string
	var textColor string
	flags := "-"
	if player == 0 {
		text = "YOU"
		textColor = "blue"
	} else {
		text = "OPPONENT"
		textColor = "red"
	}
	if player == g.CurrentPlayer {
		flags = "u"
	}
	fulltext := fmt.Sprintf("[%s::%s]%s[white]", textColor, flags, text)
	frame.Clear().AddText(fulltext, true, tview.AlignCenter, tcell.ColorWhite)
	view.SetBorders(false).SetSelectable(false, false).Clear()
	fixedRes, _ := g.Players[player].AvailableResources()
	cText := 0
	cVal := 1
	data := fixedRes.ToMap()
	labels := []string{"Coins", "Wood", "Stone", "Clay", "Glass", "Papyrus"}
	for riga, label := range labels {
		value := data[label]
		var color tcell.Color
		switch label {
		case "Coins":
			color = tcell.ColorYellow
			value = g.Players[player].Coins // non derivano da produzione, ma sono possedute
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

func fillActions(g *game.Game, view *tview.Table) {
	view.SetBorders(false).SetSelectable(false, false).Clear()
	row := 0
	for r := 0; r <= g.Board.YMax; r++ {
		for c := 0; c < g.Board.XMax; c++ {
			card := g.Board.Cards[r][c]
			if card.Visible && !g.Board.CardBlocked(&card) {
				view.SetCell(row, 0, tview.NewTableCell(card.Building.Name).SetTextColor(getTypeColor(card.Building.Type)))
				row++
			}
		}
	}
}

// Gui creates and returs main window ready to be displayed
func Gui(game *game.Game) *tview.Application {
	// create components & layout
	app := tview.NewApplication()
	youInfo := tview.NewTable()
	youInfoFrame := tview.NewFrame(youInfo)
	opponentInfo := tview.NewTable()
	opponentInfoFrame := tview.NewFrame(opponentInfo)
	mainLeftBottom := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainLeftBottom.AddItem(youInfoFrame, 0, 1, false)
	mainLeftBottom.AddItem(opponentInfoFrame, 0, 1, false)
	mainLeft := tview.NewFlex().SetDirection(tview.FlexRow)
	boardTable := tview.NewTable()
	mainLeft.AddItem(boardTable, 0, 1, false)
	mainLeft.AddItem(mainLeftBottom, 0, 1, false)
	actionsTable := tview.NewTable()
	mainRight := tview.NewFrame(actionsTable).AddText("Actions", true, tview.AlignCenter, tcell.ColorWhite)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainFlex.AddItem(mainLeft, 0, 1, false)
	mainFlex.AddItem(mainRight, 0, 1, true) // "Actions" has focus because all commands are here
	main := tview.NewFrame(mainFlex)

	refreshFunc := func() {
		if game.CurrentAge > 3 {
			app.Stop()
		} else {
			if game.CurrentPlayer > 0 {
				game.CurrentPlayer = 0
			} else {
				game.CurrentPlayer = 1
			}
			game.DeployBoard()
			title := fmt.Sprintf("7 Wonders Duel - Age %d", game.CurrentAge)
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
			main.Clear()
			main.AddText(title, true, tview.AlignCenter, titleColor)
			fillBoardTable(game, boardTable)
			fillPlayerInfoArea(game, 0, youInfo, youInfoFrame)
			fillPlayerInfoArea(game, 1, opponentInfo, opponentInfoFrame)
			fillActions(game, actionsTable)
		}
	}
	app.SetRoot(main, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// BEGIN DEBUG FUNCTIONALITY
		switch event.Rune() {
		case 'n':
			game.CurrentAge++
			refreshFunc()
		case 'q':
			app.Stop()
		case 'c':
			game.Players[0].Coins += 2
			refreshFunc()
		}
		// END DEBUG FUNCTIONALITY
		return event
	})

	// GUI done, it's time to play

	// initialization
	rand.Seed(time.Now().UnixNano())
	game.CurrentAge = 1
	game.CurrentPlayer = rand.Intn(2)
	refreshFunc()
	game.CurrentRound = 1

	return app
}
