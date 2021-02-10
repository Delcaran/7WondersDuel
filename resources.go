package main

import (
	"fmt"
)

type resource int
func (r resource) printValue(name string) {
    if(r > 0) {
        fmt.Printf("\t%v %s\n", r, name)
    }
}
func (r resource) hasValue() bool {
    return (r > 0)
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
    fmt.Println("Costs:")
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
	if (p.Coin.hasValue() || p.Wood.hasValue() || p.Clay.hasValue()|| p.Stone.hasValue()|| p.Glass.hasValue()|| p.Papyrus.hasValue() || p.Shield.hasValue()) {
		fmt.Printf("Produces")
		if(p.Choice) {
			fmt.Printf(" one of the following")
		}
		fmt.Println(" :")
		p.Coin.printValue("Coin")
		p.Wood.printValue("Wood")
		p.Clay.printValue("Clay")
		p.Stone.printValue("Stone")
		p.Glass.printValue("Glass")
		p.Papyrus.printValue("Papyrus")
		p.Shield.printValue("Shield")
	}
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
	fmt.Println("When built:")
	if(c.Points > 0) {
		fmt.Printf("\t%d Points\n", c.Points)
	}
	if(c.Turn) {
		fmt.Printf("\tExtra turn\n")
	}
	c.Coins.printValue("Coins")
	if(c.CoinsRemoved > 0) {
		fmt.Printf("\tRemoves %d coins from the opponent\n", c.CoinsRemoved)
	}
	if(len(c.Discard) > 0) {
		fmt.Printf("\tOpponent discard one building of type %s\n", c.Discard)
	}
	c.Tokens.printContent()
	c.Shield.printValue("Shield")
    c.Production.printContent()
}
