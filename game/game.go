package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type tokenChoice struct {
	Choose int
	Pick   int
}

type gameContent struct {
	Wonders []wonder
	Decks   []deck
	Tokens  []token
	Coins   int
}

func (d *gameContent) prepareContent() {
	// Shuffle wonders
	rand.Shuffle(len(d.Wonders), func(i, j int) {
		d.Wonders[i], d.Wonders[j] = d.Wonders[j], d.Wonders[i]
	})

	// Shuffle tokens
	rand.Shuffle(len(d.Tokens), func(i, j int) {
		d.Tokens[i], d.Tokens[j] = d.Tokens[j], d.Tokens[i]
	})

	// Prepare buildings
	for _, deck := range d.Decks {
		deck.prepareBuildings()
	}
}

// Player informations
type Player struct {
	Name          string
	Points        int
	Coins         int
	BonusShields  int // in addition of those from buildings
	MilitaryPower int
	Wonders       []wonder
	Buildings     []*building
	Links         []string
	Tokens        []token
}

func calculateDynamicProduction(input []*Production, inputIndex int, output []Production, tmp Production) {
	if inputIndex < len(input) {
		in := input[inputIndex]
		if in.Wood > 0 {
			newTmp := tmp
			newTmp.Wood += in.Wood
			calculateDynamicProduction(input, inputIndex+1, output, newTmp)
		}
		if in.Clay > 0 {
			newTmp := tmp
			newTmp.Clay += in.Clay
			calculateDynamicProduction(input, inputIndex+1, output, newTmp)
		}
		if in.Stone > 0 {
			newTmp := tmp
			newTmp.Stone += in.Stone
			calculateDynamicProduction(input, inputIndex+1, output, newTmp)
		}
		if in.Papyrus > 0 {
			newTmp := tmp
			newTmp.Papyrus += in.Papyrus
			calculateDynamicProduction(input, inputIndex+1, output, newTmp)
		}
		if in.Glass > 0 {
			newTmp := tmp
			newTmp.Glass += in.Glass
			calculateDynamicProduction(input, inputIndex+1, output, newTmp)
		}
	} else {
		// fine del ramo
		output = append(output, tmp)
		return
	}
}

// AvailableResources calculates dynamic and static production
func (p *Player) AvailableResources() (Production, []Production) {
	var fixed Production
	var toBeAnalized []*Production
	genericBuildings := []genericBuilding{}
	for _, b := range p.Buildings {
		genericBuildings = append(genericBuildings, b)
	}
	for _, w := range p.Wonders {
		genericBuildings = append(genericBuildings, &w)
	}
	for _, g := range genericBuildings {
		in := g.getProduction()
		if in.Choice {
			toBeAnalized = append(toBeAnalized, in)
		} else {
			fixed.Wood += in.Wood
			fixed.Clay += in.Clay
			fixed.Stone += in.Stone
			fixed.Glass += in.Glass
			fixed.Papyrus += in.Papyrus
		}
	}
	// Create all combos of dynamic resources
	var dynamic []Production
	var tmp Production
	calculateDynamicProduction(toBeAnalized, 0, dynamic, tmp)
	return fixed, dynamic
}

func (p *Player) calculatePrices(opponent *Player) cost {
	var Prices cost
	opponentFixedProduction, _ := opponent.AvailableResources() // only raw and manufactured are fixed
	Prices.Wood = 2 + opponentFixedProduction.Wood
	Prices.Clay = 2 + opponentFixedProduction.Clay
	Prices.Stone = 2 + opponentFixedProduction.Stone
	Prices.Glass = 2 + opponentFixedProduction.Glass
	Prices.Papyrus = 2 + opponentFixedProduction.Papyrus
	for _, b := range p.Buildings {
		for _, t := range b.Trade {
			if t == "wood" {
				Prices.Wood = 1
			}
			if t == "clay" {
				Prices.Clay = 1
			}
			if t == "stone" {
				Prices.Stone = 1
			}
			if t == "glass" {
				Prices.Glass = 1
			}
			if t == "papyrus" {
				Prices.Papyrus = 1
			}
		}
	}
	return Prices
}

// CalculateSellIncome returns coins gained from selling a building
func (p *Player) CalculateSellIncome() int {
	var Coins int
	Coins = 2
	for _, b := range p.Buildings {
		if b.Type == "commercial" {
			Coins++
		}
	}
	return Coins
}

// CalculateBuyingCost returns how many coins are needed to construct a building
func (p *Player) CalculateBuyingCost(b *building, opponent *Player) (bool, bool, int) {
	var MissingResources cost
	// check links for free building
	for _, l := range p.Links {
		if l == b.Linked {
			return true, true, 0 // Buyable?, Free?, Coins
		}
	}
	// check impossible to build building (not enough coins)
	if b.Cost.Coins > p.Coins {
		return false, false, b.Cost.Coins - p.Coins // Buyable?, Free?, Coins
	}
	// check available resources and prices for missing ones
	Buyable := true
	FixedProduction, DynamicProduction := p.AvailableResources()
	MissingResources = b.Cost
	MissingResources.Wood -= FixedProduction.Wood
	MissingResources.Clay -= FixedProduction.Clay
	MissingResources.Stone -= FixedProduction.Stone
	MissingResources.Papyrus -= FixedProduction.Papyrus
	MissingResources.Glass -= FixedProduction.Glass

	// Is Fixed Production enough?
	Buyable = Buyable && MissingResources.Wood <= 0
	Buyable = Buyable && MissingResources.Clay <= 0
	Buyable = Buyable && MissingResources.Stone <= 0
	Buyable = Buyable && MissingResources.Papyrus <= 0
	Buyable = Buyable && MissingResources.Glass <= 0
	if Buyable {
		return Buyable, false, b.Cost.Coins // Buyable?, Free?, Coins
	}
	// Not enough... calculate with Dynamic Production
	for _, o := range DynamicProduction {
		Buyable := true
		Buyable = Buyable && (MissingResources.Wood-o.Wood) <= 0
		Buyable = Buyable && (MissingResources.Clay-o.Clay) <= 0
		Buyable = Buyable && (MissingResources.Stone-o.Stone) <= 0
		Buyable = Buyable && (MissingResources.Papyrus-o.Papyrus) <= 0
		Buyable = Buyable && (MissingResources.Glass-o.Glass) <= 0
		if Buyable {
			return Buyable, false, b.Cost.Coins // Buyable?, Free?, Coins
		}
	}
	// Still not enough: look into trading ...
	Prices := p.calculatePrices(opponent)
	MissingResources.Coins += MissingResources.Wood * Prices.Wood
	MissingResources.Coins += MissingResources.Clay * Prices.Clay
	MissingResources.Coins += MissingResources.Stone * Prices.Stone
	MissingResources.Coins += MissingResources.Papyrus * Prices.Papyrus
	MissingResources.Coins += MissingResources.Glass * Prices.Glass
	// Ora MissingResources.Coins contiene quante monete mi servono per comprare tutto
	// Detraggo da questo costo le risorse dinamiche
	minCoins := MissingResources.Coins
	for _, o := range DynamicProduction {
		tmp := MissingResources.Coins
		tmp -= o.Wood * Prices.Wood
		tmp -= o.Clay * Prices.Clay
		tmp -= o.Stone * Prices.Stone
		tmp -= o.Papyrus * Prices.Papyrus
		tmp -= o.Glass * Prices.Glass
		if tmp < minCoins {
			minCoins = tmp
		}
	}
	// ora minCoins contiene il minimo che devo pagare per poter costruire
	if minCoins <= p.Coins {
		return true, false, minCoins // Buyable?, Free?, Coins
	}

	return false, false, MissingResources.Coins // Buyable?, Free?, Coins
}

type position struct {
	X int // 0 is left
	Y int // 0 is top
}

// Card represents a physical card on the board
type Card struct {
	Building *building
	Visible  bool
	Picked   bool
	Position position
}

type board struct {
	Cards          [][]Card // [y][x]
	AvailableCards int
	XMax           int
	YMax           int
}

func (b *board) cardBlocks(c *Card) []*Card {
	left := c.Position.X - 1
	right := c.Position.X + 1
	line := c.Position.Y - 1

	var blocked []*Card

	if line >= 0 && line <= b.YMax {
		if left >= 0 && left < b.XMax {
			if b.Cards[line][left].Building != nil {
				blocked = append(blocked, &b.Cards[line][left])
			}
		}
		if right >= 0 && right < b.XMax {
			if b.Cards[line][right].Building != nil {
				blocked = append(blocked, &b.Cards[line][right])
			}
		}
	}
	return blocked
}
func (b *board) CardBlocked(c *Card) bool {
	left := c.Position.X - 1
	right := c.Position.X + 1
	line := c.Position.Y + 1

	blocked := false

	if line >= 0 && line <= b.YMax {
		if left >= 0 && left < b.XMax {
			if b.Cards[line][left].Building != nil && !b.Cards[line][left].Picked {
				blocked = true
			}
		}
		if right >= 0 && right < b.XMax {
			if b.Cards[line][right].Building != nil && !b.Cards[line][right].Picked {
				blocked = true
			}
		}
	}
	return blocked
}

func (b *board) removeCard(c *Card) {
	c.Picked = true
	c.Visible = false
	blocked := b.cardBlocks(c)
	for _, u := range blocked {
		if !b.CardBlocked(u) {
			u.Visible = true
		}
	}
	b.AvailableCards--
}

// Game contains all the information required to play
type Game struct {
	CurrentPhase    int
	CurrentRound    int
	CurrentAge      int
	CurrentPlayer   int
	Board           board
	Tokens          []token
	DiscardedTokens []token
	BoxContent      gameContent
	Players         [2]Player
}

func loadGameContent() (gameContent, error) {
	data := gameContent{}
	filename := "conf/data.json"
	jsonFile, err := os.Open(filename)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Fatal(err)
	}

	data.prepareContent()

	return data, nil
}

func loadBoardLayout(age int, data *gameContent) board {
	datAges, err := os.Open("conf/ages.dat")
	defer datAges.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(datAges)
	blankLines := 0
	line := -1 // will be incremented before processing each line
	var ageDeck *deck
	for _, fulldeck := range data.Decks {
		if fulldeck.Age == age {
			ageDeck = &fulldeck
			break
		}
	}
	lastCard := 0

	var newBoard board
	newBoard.AvailableCards = 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		lenght := len(text)
		if lenght == 0 {
			blankLines++
			line = -1
		} else {
			if (age - 1) == blankLines {
				// this is the right layout for the requested age
				newBoard.YMax = line
				newBoard.XMax = lenght // every line in layout MUST have the same number of chars
				newLine := make([]Card, newBoard.XMax)
				for c := 0; c < lenght; c++ {
					newLine[c].Building = nil
					if text[c] != ' ' {
						newLine[c].Building = &ageDeck.Buildings[lastCard]
						newLine[c].Picked = false
						newBoard.AvailableCards++
						newLine[c].Position.X = c
						newLine[c].Position.Y = line
						newLine[c].Visible = (text[c] == 'O') // uppercase letter o
						lastCard++
					}
				}
				newBoard.Cards = append(newBoard.Cards, newLine)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return newBoard
}

// DeployBoard generates the layout for playing
func (g *Game) DeployBoard() {
	if g.CurrentRound == 0 { // at beginning of each age
		if g.CurrentPhase == ReadyToPlay && g.CurrentAge == 0 { // at beginning of game
			g.CurrentAge = 1
			g.CurrentRound++
		}
		g.Board = loadBoardLayout(g.CurrentAge, &g.BoxContent)
	}
}

func (g *Game) calculatePoints() [2]int {
	var points [2]int
	// TODO implementare
	points[0] = rand.Intn(100)
	points[1] = rand.Intn(100)
	return points
}

func (g *Game) endRound() {
	if g.Board.AvailableCards == 0 {
		if g.CurrentAge == 3 { // end game, g.CurrentPlayer will cointain the winner. -1 for draw
			points := g.calculatePoints()
			if points[0] == points[1] {
				g.CurrentPlayer = -1
			} else {
				if points[0] > points[1] {
					g.CurrentPlayer = 0
				} else {
					g.CurrentPlayer = 1
				}
			}
		} else {
			g.CurrentAge++
			g.CurrentRound = 0

			p1idx := 0
			p2idx := 1
			p1 := &g.Players[p1idx]
			p2 := &g.Players[p2idx]
			var chooser int
			if p1.MilitaryPower == p2.MilitaryPower {
				// player who played last card chooses who begins next age
				chooser = g.CurrentPlayer
			} else {
				// player with weaker military power chooses who begins next age
				if p1.MilitaryPower < p2.MilitaryPower {
					chooser = 0
				} else {
					chooser = 1
				}
			}
			g.CurrentPlayer = chooser // TODO: per ora chi deve scegliere "sceglie" sempre se stesso, poi vedremo come fare
		}
	} else {
		g.CurrentRound++
		g.switchPlayer()
	}
}

// Construct acts when the current player wants to build a card
func (g *Game) Construct(card *Card) {
	player := &g.Players[g.CurrentPlayer]
	player.Buildings = append(player.Buildings, card.Building)
	g.Board.removeCard(card)
	g.endRound()
}

// ConstructWonder acts when the current player wants to build a wonder with the selected card
func (g *Game) ConstructWonder(card *Card) {
	//player := &g.Players[g.CurrentPlayer]
	g.Board.removeCard(card)
	g.endRound()
}

// Discard acts when the current player discards a card from the board
func (g *Game) Discard(card *Card) {
	player := &g.Players[g.CurrentPlayer]
	player.Coins += player.CalculateSellIncome()
	g.Board.removeCard(card)
	g.endRound()
}

// phases of play
const (
	PlayerNamesPhase = iota
	FirstPlayerSelectionPhase
	Player1Wonder1Phase
	Player2Wonder2Phase
	Player2Wonder1Phase
	Player1Wonder2Phase
	ReadyToPlay
	EndGame
)

// Initialize match
func (g *Game) Initialize() {
	rand.Seed(time.Now().UTC().UnixNano())
	g.Player1().Name = "Leonida"
	g.Player2().Name = "Serse"
	g.CurrentPhase = PlayerNamesPhase
	g.CurrentAge = 0
	g.CurrentPlayer = -1
	g.CurrentRound = 0
	var err error
	g.BoxContent, err = loadGameContent()
	if err != nil {
		log.Fatal(err)
	}
	// Setup initial data
	g.BoxContent.Coins = 14*1 + 10*3 + 7*6
	for p := range g.Players {
		g.BoxContent.Coins -= 7
		g.Players[p].Coins = 7
	}
}

// NextPhase advances to next stage before real play
func (g *Game) NextPhase() {
	switch g.CurrentPhase {
	case Player1Wonder1Phase:
		g.switchPlayer() // from Player 1 to player 2
	case Player2Wonder1Phase:
		g.switchPlayer() // from Player 2 to player 1
	}
	g.CurrentPhase++
}

// SetPlayer1Turn sets turn to player 1
func (g *Game) SetPlayer1Turn() {
	g.CurrentPlayer = 0
}

// SetPlayer2Turn sets turn to player 1
func (g *Game) SetPlayer2Turn() {
	g.CurrentPlayer = 1
}

// SetRandomPlayerTurn sets turn to random player
func (g *Game) SetRandomPlayerTurn() {
	g.CurrentPlayer = rand.Intn(2)
}

// Player1 pointer
func (g *Game) Player1() *Player {
	return &g.Players[0]
}

// Player2 pointer
func (g *Game) Player2() *Player {
	return &g.Players[1]
}

// IsFirst returns true if player is the first of the whole match
func (g *Game) IsFirst(p *Player) bool {
	return &g.Players[0] == p
}

// GetCurrentPlayer returns pointer to current player
func (g *Game) GetCurrentPlayer() *Player {
	return &g.Players[g.CurrentPlayer]
}

// GetOtherPlayer returns pointer to other player
func (g *Game) GetOtherPlayer() *Player {
	if g.CurrentPlayer == 0 {
		return &g.Players[1]
	}
	return &g.Players[0]
}

func (g *Game) switchPlayer() {
	if g.CurrentPlayer == 0 {
		g.CurrentPlayer = 1
	} else {
		g.CurrentPlayer = 0
	}
}

func test(d, p1, p2 int) {
	local := fmt.Sprintf("%d %d %d", d, p1, p2)
	fmt.Println(local)
}

// AddWonders adds some wonders to players
func (g *Game) AddWonders(selected []int, available []int) {
	addAvailable := g.CurrentPhase == Player1Wonder2Phase || g.CurrentPhase == Player2Wonder2Phase

	test(len(g.BoxContent.Wonders), len(g.Player1().Wonders), len(g.Player2().Wonders))

	// aggiunta
	for _, idx := range selected {
		g.GetCurrentPlayer().Wonders = append(g.GetCurrentPlayer().Wonders, g.BoxContent.Wonders[idx])
	}
	if addAvailable {
		for _, idx := range available { // qui dovrebbe essercene solo una
			g.GetOtherPlayer().Wonders = append(g.GetOtherPlayer().Wonders, g.BoxContent.Wonders[idx])
		}
	}

	test(len(g.BoxContent.Wonders), len(g.Player1().Wonders), len(g.Player2().Wonders))

	// rimozione
	for _, idx := range selected {
		g.BoxContent.Wonders = append(g.BoxContent.Wonders[:idx], g.BoxContent.Wonders[idx+1:]...)
	}
	if addAvailable {
		for _, idx := range available { // qui dovrebbe essercene solo una
			g.BoxContent.Wonders = append(g.BoxContent.Wonders[:idx], g.BoxContent.Wonders[idx+1:]...)
		}
	}

	test(len(g.BoxContent.Wonders), len(g.Player1().Wonders), len(g.Player2().Wonders))
}
