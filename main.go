package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	data, err := loadGameContent()
	if err != nil {
		fmt.Println("Error in loading game content")
	} else {
		dataJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("%s\n", string(dataJSON))
	}
}
