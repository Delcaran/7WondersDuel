package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Wonder struct {
	name string
}

type Data struct {
	wonders []Wonder
}

func main() {
	fmt.Printf("ciao\n")
    filename := "data.json"
    jsonFile, err := os.Open(filename)
    defer jsonFile.Close()
    if err != nil {
        fmt.Printf("failed to open json file: %s, error: %v", filename, err)
        return
    }

    jsonData, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        fmt.Printf("failed to read json file, error: %v", err)
        return
    }

    data := Data{}
    if err := json.Unmarshal(jsonData, &data); err != nil {
        fmt.Printf("failed to unmarshal json file, error: %v", err)
        return
    }

    // Print
    for _, wonder := range data.wonders {
        fmt.Printf("Name: %s \n", wonder.name)
    }
}
