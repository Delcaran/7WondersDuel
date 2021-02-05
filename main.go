package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type resource int

func (r resource) printValue(name string) {
    if(r > 0) {
        fmt.Printf("%v %s\n", r, name)
    }
}

type printContent interface {
    printContent()
}

type cost struct {
    Coin resource
    Wood resource
    Clay resource
    Stone resource
    Glass resource
    Papyrus resource   
}
func (c cost) printContent() {
    fmt.Printf("Costs: \n")
    c.Coin.printValue("Coin")
    c.Wood.printValue("Wood")
    c.Clay.printValue("Clay")
    c.Stone.printValue("Stone")
    c.Glass.printValue("Glass")
    c.Papyrus.printValue("Papyrus")
}

type production struct {
    Coin resource
    Wood resource
    Clay resource
    Stone resource
    Glass resource
    Papyrus resource
    Shield resource   
    Choice bool
}
func (p production) printContent() {
    fmt.Printf("Produces ")
    if(p.Choice) {
        fmt.Printf("one of the following")
    }
    fmt.Printf(" :\n")
    p.Coin.printValue("Coin")
    p.Wood.printValue("Wood")
    p.Clay.printValue("Clay")
    p.Stone.printValue("Stone")
    p.Glass.printValue("Glass")
    p.Papyrus.printValue("Papyrus")
    p.Shield.printValue("Shield")
}

type tokenChoice struct {
    Choose int
    Pick int
}

type construction struct {
    Points int
    Turn bool
    Coins resource
    CoinsRemoved resource   
    Discard string    
    Tokens tokenChoice
    Shield resource
    Production production
}
func (c construction) printContent() {
    fmt.Printf("Gives you:\n")
    /*
    if(c.Production != nil) {
        fmt.Printf("And ")
        c.Production.printContent()
    }
    */
}

type wonder struct {
    Name string
    Production production
    Construction construction
    Cost cost
    TokenChoice tokenChoice
}
func (w wonder) print() {
    fmt.Printf("\"%s\" \n", w.Name)
    w.Cost.printContent()
    w.Construction.printContent()
}

type importData struct {
    Wonders []wonder
}

func main() {
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

    data := importData{}
    if err := json.Unmarshal(jsonData, &data); err != nil {
        fmt.Printf("failed to unmarshal json file, error: %v", err)
        return
    }

    // Print
    for _, wonder := range data.Wonders {
        wonder.print()
        fmt.Printf("\n")
    }
}
