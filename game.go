package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type tokenChoice struct {
    Choose int
    Pick int
}
func (t tokenChoice) printContent() {
	if(t.Choose > 0 && t.Pick > 0 && t.Pick < t.Choose) {
		fmt.Printf("\tExtract %d tokens and pick %d\n", t.Choose, t.Pick)
	}
}

type importData struct {
    Wonders []wonder
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