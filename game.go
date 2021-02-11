package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Building building
	Visible  bool
	Position position
}

type board struct {
	Cards      [][]card // [y][x]
	XMin, XMax int
	YMin, YMax int
}

func (b *board) cardBlocks(c *card) [2]*card {
	left := c.Position.X - 1
	right := c.Position.X + 1
	line := c.Position.Y - 1

	var blocked [2]*card

	if line >= b.YMin && line <= b.YMax {
		if left >= b.XMin && left <= b.XMax {
			if b.Cards[line][left].Building.ID != "" {
				blocked[0] = &b.Cards[line][left]
			}
		}
		if right >= b.XMin && right <= b.XMax {
			if b.Cards[line][right].Building.ID != "" {
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

	if line >= b.YMin && line <= b.YMax {
		if left >= b.XMin && left <= b.XMax {
			if b.Cards[line][left].Building.ID != "" {
				blocked = true
			}
		}
		if right >= b.XMin && right <= b.XMax {
			if b.Cards[line][right].Building.ID != "" {
				blocked = true
			}
		}
	}
	return blocked
}

type game struct {
	CurrentPlayer int
	Board         board
}

func loadGameContent() (importData, error) {
	data := importData{}
	filename := "data.json"
	jsonFile, err := os.Open(filename)
	defer jsonFile.Close()
	if err != nil {
		fmt.Printf("failed to open json file: %s, error: %v", filename, err)
		return data, err
	}

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return data, err
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return data, err
	}
	return data, nil
}

func deployBoard() {
	// boxContents, err := loadGameContent()
	// if err != nil {
	// 	fmt.Printf("failed load box content, error: %v", err)
	// 	return
	// }
}
