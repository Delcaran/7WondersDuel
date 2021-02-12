package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type tokenChoice struct {
	Choose int
	Pick   int
}

type importData struct {
	Wonders []wonder
	Decks   []deck
	Tokens  []token
	Coins   resource
}

type player struct {
	Points        int
	Coins         resource
	BonusShields  resource // in addition of those from buildings
	MilitaryPower int
	Wonders       []wonder
	Raw           []building
	Manufactured  []building
	Military      []building
	Scientific    []building
	Civilian      []building
	Commercial    []building
	Guild         []building
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
		if left >= 0 && left <= b.XMax {
			if b.Cards[line][left].Building != nil {
				blocked = true
			}
		}
		if right >= 0 && right <= b.XMax {
			if b.Cards[line][right].Building != nil {
				blocked = true
			}
		}
	}
	return blocked
}
func (b *board) print() {
	for y := 0; y <= b.YMax; y++ {
		for x := 0; x < b.XMax; x++ {
			if x != 0 {
				fmt.Printf(" ")
			}
			if b.Cards[y][x].Building != nil {
				fmt.Printf("%s", b.Cards[y][x].Building.ID)
			} else {
				fmt.Printf("*")
			}
		}
		fmt.Printf("\n")
	}
}

type game struct {
	CurrentPlayer int
	Board         board
}

func loadGameContent() (importData, error) {
	data := importData{}
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
	return data, nil
}

func loadBoardLayout(age int, data *importData) board {
	//rand.Seed(time.Now())
	rand.Seed(13234)
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
	deckCards := len(ageDeck.Buildings)
	var usedCards []int

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
						newRand := rand.Intn(deckCards)
						for {
							randOk := true
							for _, n := range usedCards {
								if n == newRand {
									randOk = false
								}
							}
							if randOk {
								break
							} else {
								newRand = rand.Intn(deckCards)
							}
						}
						usedCards = append(usedCards, newRand)
						newLine[c].Building = &ageDeck.Buildings[newRand]
						newLine[c].Position.X = c
						newLine[c].Position.Y = line
						newLine[c].Visible = (text[c] == 'O') // uppercase letter o
						//fmt.Printf("%s in ( y : %d , x : %d)\n", newLine[c].Building.ID, newLine[c].Position.Y, newLine[c].Position.X)
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

func deployBoard() {
	boxContents, err := loadGameContent()
	if err != nil {
		log.Fatal(err)
	}
	board := loadBoardLayout(1, &boxContents)
	board.print()
}
