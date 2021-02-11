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

type game struct {
	CurrentPlayer int
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
