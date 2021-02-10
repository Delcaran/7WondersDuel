package main

import "fmt"

type wonder struct {
    Name string
    Production production
    Construction construction
    Cost cost
    TokenChoice tokenChoice
}
func (w wonder) print() {
    fmt.Printf("\"%s\"\n", w.Name)
    w.Cost.printContent()
    w.Construction.printContent()
}