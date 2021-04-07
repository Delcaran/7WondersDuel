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

func getTypeColor(card *game.Card) tcell.Color {
	color := tcell.ColorWhite
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
	return color
}

func getTypeColorString(card *game.Card) string {
	color := "white"
	switch card.Building.Type {
	case "raw":
		color = "brown"
	case "manufactured":
		color = "darkgrey"
	case "commercial":
		color = "goldenrod"
	case "military":
		color = "red"
	case "guild":
		color = "purple"
	case "civilian":
		color = "blue"
	case "scientific":
		color = "green"
	default:
		color = "white"
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
					color = getTypeColor(&card)
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

func evaluateBuildability(g *game.Game, card *game.Card) (bool, bool, int) {
	otherPlayer := 0
	if g.CurrentPlayer == otherPlayer {
		otherPlayer = 1
	}

	buyable, free, coins := g.Players[g.CurrentPlayer].CalculateBuyingCost(card.Building, &g.Players[otherPlayer])
	return buyable, free, coins
}

func appendValue(text *string, value int, color string) int {
	if value > 0 {
		if len(*text) > 0 {
			*text = fmt.Sprintf("%s,", *text)
		}
		*text = fmt.Sprintf("%s[%s]%d[white]", *text, color, value)
	}
	return value
}

func getBuildingSummary(card *game.Card) string {
	var txtCost, txtProd, txtOTOG string
	cost := card.Building.Cost
	prod := card.Building.Production
	otog := card.Building.Construction

	totCost := 0
	totCost += appendValue(&txtCost, cost.Coins, "yellow")
	totCost += appendValue(&txtCost, cost.Wood, "brown")
	totCost += appendValue(&txtCost, cost.Stone, "grey")
	totCost += appendValue(&txtCost, cost.Clay, "orange")
	totCost += appendValue(&txtCost, cost.Glass, "lightblue")
	totCost += appendValue(&txtCost, cost.Papyrus, "goldenrod")
	// TODO: links
	if totCost > 0 {
		txtCost = fmt.Sprintf("[white]C: %s[white]", txtCost)
	}

	totProd := 0
	totProd += appendValue(&txtProd, prod.Wood, "brown")
	totProd += appendValue(&txtProd, prod.Stone, "grey")
	totProd += appendValue(&txtProd, prod.Clay, "orange")
	totProd += appendValue(&txtProd, prod.Glass, "lightblue")
	totProd += appendValue(&txtProd, prod.Papyrus, "goldenrod")
	totProd += appendValue(&txtProd, prod.Shield, "red")
	// TODO: casi speciali
	if totProd > 0 {
		txtProd = fmt.Sprintf("[white]P: %s[white]", txtProd)
	}

	totOTOG := 0
	totOTOG += appendValue(&txtOTOG, otog.Coins, "yellow")
	totOTOG += appendValue(&txtOTOG, card.Building.Points, "blue")
	// TODO: casi speciali
	if totOTOG > 0 {
		txtOTOG = fmt.Sprintf("[white]G: %s[white]", txtOTOG)
	}

	text := txtCost
	if len(txtProd) > 0 {
		if len(text) > 0 {
			text = fmt.Sprintf("%s - ", text)
		}
		text = fmt.Sprintf("%s%s", text, txtProd)
	}
	if len(txtOTOG) > 0 {
		if len(text) > 0 {
			text = fmt.Sprintf("%s - ", text)
		}
		text = fmt.Sprintf("%s%s", text, txtOTOG)
	}
	if len(card.Building.Science) > 0 {
		if len(text) > 0 {
			text = fmt.Sprintf("%s - ", text)
		}
		text = fmt.Sprintf("%sS: [green]%s[white]", text, card.Building.Science)
	}
	if len(card.Building.Trade) > 0 {
		if len(text) > 0 {
			text = fmt.Sprintf("%s - ", text)
		}
		txtTrade := ""
		for _, t := range card.Building.Trade {
			if len(txtTrade) > 0 {
				txtTrade = fmt.Sprintf("%s,", txtTrade)
			}
			txtTrade = fmt.Sprintf("%s[goldenrod]%s[white]", txtTrade, t)
		}
		text = fmt.Sprintf("%s[white]T: %s", text, txtTrade)
	}

	return text
}

func fillActions(g *game.Game, view *tview.List, actions *tview.List, actionsFrame *tview.Frame) {
	cardSelectors := "0123456789"
	cardRunes := []rune(cardSelectors)
	view.Clear()
	row := 0
	for r := 0; r <= g.Board.YMax; r++ {
		for c := 0; c < g.Board.XMax; c++ {
			card := g.Board.Cards[r][c]
			if card.Visible && !g.Board.CardBlocked(&card) {
				text := fmt.Sprintf("[%s]%s[white]", getTypeColorString(&card), card.Building.Name)
				buyable, free, coins := evaluateBuildability(g, &card)
				sellIncome := g.Players[g.CurrentPlayer].CalculateSellIncome()
				view.AddItem(text, getBuildingSummary(&card), cardRunes[row], func() {
					actions.Clear()
					if buyable { // check if player can construct the card
						var subtext string
						if free || coins <= 0 {
							subtext = "[white]You can build for free"
						} else {
							subtext = fmt.Sprintf("[white]You have to spend [yellow]%d[white] extra coins", coins)
						}
						actions.AddItem(fmt.Sprintf("Construct %s", text), subtext, 'c', func() {})
					}
					actions.AddItem(fmt.Sprintf("Discard %s", text), fmt.Sprintf("[white]You will earn [yellow]%d[white] coins", sellIncome), 'd', func() {})
					actions.AddItem(fmt.Sprintf("Use %s to construct a Wonder", text), "", 'w', func() {})
				})
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
	activeCardsList := tview.NewList()
	actionsList := tview.NewList()
	mainRightTop := tview.NewFrame(activeCardsList).AddText("USABLE CARDS", true, tview.AlignCenter, tcell.ColorWhite)
	mainRightBottom := tview.NewFrame(actionsList)
	mainRight := tview.NewFlex().SetDirection(tview.FlexRow)
	mainRight.AddItem(mainRightTop, 0, 1, false)
	mainRight.AddItem(mainRightBottom, 0, 1, false)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainFlex.AddItem(mainLeft, 0, 1, false)
	mainFlex.AddItem(mainRight, 0, 1, true) // "USABLE CARDS" has focus because all commands start from here
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
			fillActions(game, activeCardsList, actionsList, mainRightBottom)
		}
	}
	app.SetRoot(main, true)
	app.SetFocus(activeCardsList)

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
