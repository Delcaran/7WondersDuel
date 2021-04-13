package gui

import (
	"7WondersDuel/game"
	"fmt"
	"strconv"

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
func fillBoard(game *game.Game, board *tview.Grid) {
	board.Clear().SetBorders(false).SetColumns(game.Board.XMax).SetRows(game.Board.YMax)
	for r := 0; r <= game.Board.YMax; r++ {
		for c := 0; c < game.Board.XMax; c++ {
			card := &game.Board.Cards[r][c]
			if !card.Picked && card.Building != nil {
				var color tcell.Color
				var cardName, extendedInfo, flag string
				extendedInfoView := tview.NewTextView().SetDynamicColors(true)
				extendedInfoView.SetTextAlign(tview.AlignLeft).SetTitleAlign(tview.AlignCenter).SetBorder(true)
				if card.Visible {
					color = getTypeColor(card)
					// default production for this kind of card
					switch card.Building.Type {
					case "raw":
						extendedInfo = getBuildingProduction(card)
					case "manufactured":
						extendedInfo = getBuildingProduction(card)
					case "commercial":
						extendedInfo = getBuildingTrade(card)
					case "military":
						extendedInfo = fmt.Sprintf("Shields: [red]%d[white]", card.Building.Production.Shield)
					case "guild":
						extendedInfo = "TODO"
					case "civilian":
						extendedInfo = getBuildingOTOG(card)
					case "scientific":
						extendedInfo = getBuildingScience(card)
					default:
						extendedInfo = ""
					}
					// extra output
					extendedInfo = fmt.Sprintf("%s\n%s", extendedInfo, getBuildingLinks(card))
					if card.Building.Type != "civilian" {
						extendedInfo = fmt.Sprintf("%s %s", extendedInfo, getBuildingOTOG(card))
					}
					// cost
					extendedInfo = fmt.Sprintf("%s\n%s", extendedInfo, getBuildingCost(card))
					// links
					if len(card.Building.Linked) > 0 {
						extendedInfo = fmt.Sprintf("%s\n%s", extendedInfo, card.Building.Linked)
					}
					if !game.Board.CardBlocked(card) {
						flag = "b"
						extendedInfoView.SetBorderAttributes(tcell.AttrBold)
					} else {
						flag = ""
						extendedInfoView.SetBorderAttributes(tcell.AttrDim)
					}
					cardName = fmt.Sprintf("[::%s]%s[-:-:-]", flag, card.Building.Name)
				} else {
					switch game.CurrentAge {
					case 1:
						color = tcell.ColorDarkGoldenrod
					case 2:
						color = tcell.ColorLightBlue
					case 3:
						color = tcell.ColorViolet
						if card.Building.Type == "guild" {
							color = tcell.ColorPurple
						}
					default:
						color = tcell.ColorWhite
					}
					extendedInfoView.SetBackgroundColor(color)
				}
				extendedInfoView.SetText(extendedInfo).SetTitle(cardName).SetTitleColor(color).SetBorderColor(color)
				board.AddItem(extendedInfoView, r, c, 1, 2, 0, 0, false)
			}
		}
	}
}

func fillPlayerInfoArea(g *game.Game, player int, view *tview.Table, frame *tview.Frame) {
	var textColor string
	var borderColor tcell.Color
	text := g.Players[player].Name
	flags := "-"
	if player == 0 {
		textColor = "blue"
		borderColor = tcell.ColorBlue
	} else {
		textColor = "red"
		borderColor = tcell.ColorRed
	}
	if player == g.CurrentPlayer {
		flags = "b"
		frame.SetBorderColor(borderColor).SetBorderAttributes(tcell.AttrBold | tcell.AttrReverse)
	} else {
		frame.SetBorderColor(borderColor).SetBorderAttributes(tcell.AttrDim)
	}
	fulltext := fmt.Sprintf("[%s::%s]%s[white]", textColor, flags, text)
	frame.Clear().SetTitle(fulltext)
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
	// TODO: points
	// TODO: military power
	// TODO: dynamic production
	// TODO: trade bonus
	// TODO: tokens
	// TODO: wonders
	// TODO: end-game bonus (da gilde ecc)
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

func appendBlock(text *string, block string) {
	if len(block) > 0 {
		if len(*text) > 0 {
			*text = fmt.Sprintf("%s ", *text)
		}
		*text = fmt.Sprintf("%s%s", *text, block)
	}
}

func getBuildingCost(card *game.Card) string {
	cost := card.Building.Cost
	var txtCost string
	totCost := 0
	totCost += appendValue(&txtCost, cost.Coins, "yellow")
	totCost += appendValue(&txtCost, cost.Wood, "brown")
	totCost += appendValue(&txtCost, cost.Stone, "grey")
	totCost += appendValue(&txtCost, cost.Clay, "orange")
	totCost += appendValue(&txtCost, cost.Glass, "lightblue")
	totCost += appendValue(&txtCost, cost.Papyrus, "goldenrod")
	if totCost > 0 {
		txtCost = fmt.Sprintf("[white]C:%s[white]", txtCost)
	}
	return txtCost
}

func getBuildingProduction(card *game.Card) string {
	prod := card.Building.Production
	var txtProd string
	totProd := 0
	totProd += appendValue(&txtProd, prod.Wood, "brown")
	totProd += appendValue(&txtProd, prod.Stone, "grey")
	totProd += appendValue(&txtProd, prod.Clay, "orange")
	totProd += appendValue(&txtProd, prod.Glass, "lightblue")
	totProd += appendValue(&txtProd, prod.Papyrus, "goldenrod")
	totProd += appendValue(&txtProd, prod.Shield, "red")
	// TODO: casi speciali
	if totProd > 0 {
		txtProd = fmt.Sprintf("[white]P:%s[white]", txtProd)
	}
	return txtProd
}

func getBuildingTrade(card *game.Card) string {
	txtTrade := ""
	if len(card.Building.Trade) > 0 {
		for _, t := range card.Building.Trade {
			if len(txtTrade) > 0 {
				txtTrade = fmt.Sprintf("%s,", txtTrade)
			}
			txtTrade = fmt.Sprintf("%s[goldenrod]%s[white]", txtTrade, t)
		}
		txtTrade = fmt.Sprintf("T:%s", txtTrade)
	}
	return txtTrade
}

func getBuildingScience(card *game.Card) string {
	txtScience := ""
	if len(card.Building.Science) > 0 {
		txtScience = fmt.Sprintf("S:[green]%s[white]", card.Building.Science)
	}
	return txtScience
}

func getBuildingOTOG(card *game.Card) string {
	var txtOTOG string
	otog := card.Building.Construction
	totOTOG := 0
	totOTOG += appendValue(&txtOTOG, otog.Coins, "yellow")
	totOTOG += appendValue(&txtOTOG, card.Building.Points, "blue")
	// TODO: casi speciali
	if totOTOG > 0 {
		txtOTOG = fmt.Sprintf("[white]G:%s[white]", txtOTOG)
	}
	return txtOTOG
}

func getBuildingLinks(card *game.Card) string {
	txtLinks := ""
	if len(card.Building.Links) > 0 {
		txtLinks = fmt.Sprintf("B:%s", card.Building.Links)
	}
	return txtLinks
}

func getBuildingSummary(card *game.Card) string {
	txtCost := getBuildingCost(card)
	if len(card.Building.Linked) > 0 {
		if len(txtCost) > 0 {
			txtCost = fmt.Sprintf("%s/", txtCost)
		}
		txtCost = fmt.Sprintf("%s%s", txtCost, card.Building.Linked)
	}

	txtProd := getBuildingProduction(card)
	txtOTOG := getBuildingOTOG(card)
	txtScience := getBuildingScience(card)
	txtTrade := getBuildingTrade(card)
	txtLinks := getBuildingLinks(card)

	text := txtCost
	appendBlock(&text, txtProd)
	appendBlock(&text, txtOTOG)
	appendBlock(&text, txtScience)
	appendBlock(&text, txtTrade)
	appendBlock(&text, txtLinks)

	return text
}

type componentsGUI struct {
	app                                          *tview.Application
	main, p1InfoFrame, p2InfoFrame, actionsFrame *tview.Frame
	p1Info, p2Info                               *tview.Table
	board                                        *tview.Grid
	activeCardsList, actionsList                 *tview.List
	mainFlex, bottomFlex, topFlex, actionsFlex   *tview.Flex
}

func fillActions(g *game.Game, gui *componentsGUI) {
	cardSelectors := "0123456789"
	cardRunes := []rune(cardSelectors)
	row := 0
	view := gui.activeCardsList
	gui.actionsFrame.SetTitle(fmt.Sprintf("%s, it's your turn", g.Players[g.CurrentPlayer].Name))
	view.Clear()
	for r := 0; r <= g.Board.YMax; r++ {
		for c := 0; c < g.Board.XMax; c++ {
			card := &g.Board.Cards[r][c]
			if !card.Picked && card.Visible && !g.Board.CardBlocked(card) {
				actions := gui.actionsList
				text := fmt.Sprintf("[%s]%s[white]", getTypeColorString(card), card.Building.Name)
				buyable, free, coins := evaluateBuildability(g, card)
				sellIncome := g.Players[g.CurrentPlayer].CalculateSellIncome()
				view.AddItem(text, getBuildingSummary(card), cardRunes[row], func() {
					gui.app.SetFocus(actions)
					actions.Clear()
					actions.AddItem("BACK", "Bo back to building selection", 'b', func() {
						actions.Clear()
						refresh(g, gui)
					})
					if buyable { // check if player can construct the card
						var subtext string
						if free || coins <= 0 {
							subtext = "[white]You can build for free"
						} else {
							subtext = fmt.Sprintf("[white]You can build spending [yellow]%d[white] coins", coins)
						}
						actions.AddItem(fmt.Sprintf("Construct %s", text), subtext, 'c', func() {
							g.Construct(card)
							actions.Clear()
							refresh(g, gui)
						})
					}
					actions.AddItem(fmt.Sprintf("Discard %s", text), fmt.Sprintf("[white]You will earn [yellow]%d[white] coins", sellIncome), 'd', func() {
						g.Discard(card)
						actions.Clear()
						refresh(g, gui)
					})
					actions.AddItem(fmt.Sprintf("Use %s to construct a Wonder", text), "", 'w', func() {
						g.ConstructWonder(card)
						actions.Clear()
						refresh(g, gui)
					})
				})
				row++
			}
		}
	}
}

func drawMain(game *game.Game, gui *componentsGUI) {
	gui.app.SetRoot(gui.main, true)
	gui.app.SetFocus(gui.activeCardsList)

	gui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// BEGIN DEBUG FUNCTIONALITY
		switch event.Rune() {
		case 'n':
			game.CurrentRound = 0
			game.CurrentAge++
			refresh(game, gui)
		case 'q':
			gui.app.Stop()
		}
		// END DEBUG FUNCTIONALITY
		return event
	})

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
	gui.main.Clear()
	gui.main.AddText(title, true, tview.AlignCenter, titleColor)
	fillBoard(game, gui.board)
	// TODO: stabilire se *io* sono sempre il player 0 o meno...
	fillPlayerInfoArea(game, 0, gui.p1Info, gui.p1InfoFrame)
	fillPlayerInfoArea(game, 1, gui.p2Info, gui.p2InfoFrame)
	fillActions(game, gui)
}

func refresh(game *game.Game, gui *componentsGUI) {
	if game.CurrentAge > 3 {
		gui.app.Stop()
		// display the winner
		fmt.Printf("\n\n%s WINS!\n\n", game.Players[game.CurrentPlayer].Name)
	} else {
		if game.CurrentAge == 0 && (game.Players[0].Name == "" || game.Players[1].Name == "") {
			game.Players[0].Name = "Leonida"
			game.Players[1].Name = "Serse"
			form := tview.NewForm().
				AddInputField("Player 1", game.Players[0].Name, 20, nil, func(text string) {
					game.Players[0].Name = text
				}).
				AddInputField("Player 2", game.Players[1].Name, 20, nil, func(text string) {
					game.Players[1].Name = text
				}).
				AddButton("Start", func() {
					refresh(game, gui)
				}).
				AddButton("Quit", func() {
					gui.app.Stop()
				})
			form.SetBorder(true).SetTitle("Enter player's names").SetTitleAlign(tview.AlignCenter)
			gui.app.SetRoot(form, true).SetFocus(form)
		} else {
			if game.CurrentRound == 0 {
				// choose who begins this age
				var txt string
				if game.CurrentAge == 0 {
					txt = "Who will begin the game?"
				} else {
					txt = fmt.Sprintf("Who will begin Age %d?", game.CurrentAge)
				}
				modal := tview.NewModal().SetText(txt).AddButtons([]string{game.Players[0].Name, game.Players[1].Name}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					game.CurrentPlayer = buttonIndex
					drawMain(game, gui)
				})
				gui.app.SetRoot(modal, true).SetFocus(modal)
			}
		}
	}
}

// Gui creates and returs main window ready to be displayed
func Gui(game *game.Game) *tview.Application {
	// create components & layout
	var myGUI componentsGUI

	myGUI.topFlex = tview.NewFlex().SetDirection(tview.FlexColumn) // parte superiore: plancia e comandi
	myGUI.board = tview.NewGrid()
	myGUI.topFlex.AddItem(myGUI.board, 0, 2, false)
	myGUI.actionsFlex = tview.NewFlex().SetDirection(tview.FlexColumn) // parte superiore destra: carte disponibili e azioni
	myGUI.actionsFrame = tview.NewFrame(myGUI.actionsFlex)
	myGUI.topFlex.AddItem(myGUI.actionsFrame, 0, 1, false)
	myGUI.activeCardsList = tview.NewList()
	myGUI.actionsList = tview.NewList()
	myGUI.actionsFlex.AddItem(myGUI.activeCardsList, 0, 1, true)
	myGUI.actionsFlex.AddItem(myGUI.actionsList, 0, 1, false)

	myGUI.bottomFlex = tview.NewFlex().SetDirection(tview.FlexColumn) // parte inferiore: info dei due giocatori
	myGUI.p1Info = tview.NewTable()
	myGUI.p1InfoFrame = tview.NewFrame(myGUI.p1Info)
	myGUI.bottomFlex.AddItem(myGUI.p1InfoFrame, 0, 1, false)
	myGUI.p2Info = tview.NewTable()
	myGUI.p2InfoFrame = tview.NewFrame(myGUI.p2Info)
	myGUI.bottomFlex.AddItem(myGUI.p2InfoFrame, 0, 1, false)

	myGUI.mainFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	myGUI.mainFlex.AddItem(myGUI.topFlex, 0, 1, false)
	myGUI.mainFlex.AddItem(myGUI.bottomFlex, 0, 1, false)

	myGUI.main = tview.NewFrame(myGUI.mainFlex)

	myGUI.actionsFrame.SetBorder(true).SetTitleAlign(tview.AlignCenter).SetTitleColor(tcell.ColorWhite)
	myGUI.p1InfoFrame.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	myGUI.p2InfoFrame.SetBorder(true).SetTitleAlign(tview.AlignCenter)

	myGUI.app = tview.NewApplication()
	// GUI done, it's time to play

	// initialization
	game.CurrentPlayer = -1
	refresh(game, &myGUI)

	return myGUI.app
}
