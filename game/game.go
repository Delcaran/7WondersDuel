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
	Coins   resource
}

func (d *gameContent) prepareContent() {
	rand.Seed(time.Now().UTC().UnixNano())

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

type player struct {
	Points        int
	Coins         resource
	BonusShields  resource // in addition of those from buildings
	MilitaryPower int
	Wonders       []wonder
	Buildings     []building
	Links         []string
}

func calculateDynamicProduction(input []*production, inputIndex int, output []production, tmp production) {
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

func (p *player) availableResources() (production, []production) {
	var fixed production
	var toBeAnalized []*production
	genericBuildings := []genericBuilding{}
	for _, b := range p.Buildings {
		genericBuildings = append(genericBuildings, &b)
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
	var dynamic []production
	var tmp production
	calculateDynamicProduction(toBeAnalized, 0, dynamic, tmp)
	return fixed, dynamic
}

func (p *player) calculatePrices(opponent *player) cost {
	var Prices cost
	opponentFixedProduction, _ := opponent.availableResources() // only raw and manufactured are fixed
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

func (p *player) calculateSellIncome() resource {
	var Coins resource
	Coins = 2
	for _, b := range p.Buildings {
		if b.Type == "commercial" {
			Coins++
		}
	}
	return Coins
}

func (p *player) calculateBuyingCost(b *building, opponent *player) (bool, bool, resource) {
	var MissingResources cost
	// check links for free building
	for _, l := range p.Links {
		if l == b.CreationLink {
			return true, true, 0 // Buyable?, Free?, Coins
		}
	}
	// check impossible to build building (not enough coins)
	if b.Cost.Coins > p.Coins {
		return false, false, b.Cost.Coins - p.Coins // Buyable?, Free?, Coins
	}
	// check available resources and prices for missing ones
	Buyable := true
	FixedProduction, DynamicProduction := p.availableResources()
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

type card struct {
	Building *building
	Visible  bool
	Position position
}

type board struct {
	Cards [][]card // [y][x]
	XMax  int
	YMax  int
}

func (b *board) cardBlocks(c *card) [2]*card {
	left := c.Position.X - 1
	right := c.Position.X + 1
	line := c.Position.Y - 1

	var blocked [2]*card

	if line >= 0 && line <= b.YMax {
		if left >= 0 && left <= b.XMax {
			if b.Cards[line][left].Building != nil {
				blocked[0] = &b.Cards[line][left]
			}
		}
		if right >= 0 && right <= b.XMax {
			if b.Cards[line][right].Building != nil {
				blocked[1] = &b.Cards[line][right]
			}
		}
	}
	return blocked
}
func (b *board) cardBlocked(c *card) bool {
	left := c.Position.X - 1
	right := c.Position.X + 1
	line := c.Position.Y + 1

	blocked := false

	if line >= 0 && line <= b.YMax {
		if left >= 0 && left < b.XMax {
			if b.Cards[line][left].Building != nil {
				blocked = true
			}
		}
		if right >= 0 && right < b.XMax {
			if b.Cards[line][right].Building != nil {
				blocked = true
			}
		}
	}
	return blocked
}
func (b *board) debugPrint() {
	for y := 0; y <= b.YMax; y++ {
		for x := 0; x < b.XMax; x++ {
			if x != 0 {
				fmt.Printf(" ")
			}
			card := b.Cards[y][x]
			if card.Building != nil {
				format := "%s"
				if !card.Visible {
					format = "?%s?"
				}
				if b.cardBlocked(&card) {
					format = fmt.Sprintf("#%s#", format)
				}
				fmt.Printf(format, card.Building.ID)
			} else {
				fmt.Printf("*")
			}
		}
		fmt.Printf("\n")
	}
}

type Game struct {
	CurrentAge      int
	CurrentPlayer   int
	Board           board
	Tokens          []token
	DiscardedTokens []token
	BoxContent      gameContent
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
				newLine := make([]card, newBoard.XMax)
				for c := 0; c < lenght; c++ {
					newLine[c].Building = nil
					if text[c] != ' ' {
						newLine[c].Building = &ageDeck.Buildings[lastCard]
						newLine[c].Position.X = c
						newLine[c].Position.Y = line
						newLine[c].Visible = (text[c] == 'O') // uppercase letter o
						//fmt.Printf("%s in ( y : %d , x : %d)\n", newLine[c].Building.ID, newLine[c].Position.Y, newLine[c].Position.X)
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

func (g *Game) deployBoard() {
	if g.CurrentAge == 1 {
		var err error
		g.BoxContent, err = loadGameContent()
		if err != nil {
			log.Fatal(err)
		}
	}
	g.Board = loadBoardLayout(g.CurrentAge, &g.BoxContent)
}
