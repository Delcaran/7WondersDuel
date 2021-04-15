package game

import (
	"math/rand"
	"time"
)

// Wonder is a special type of building
type Wonder struct {
	ID           string
	Name         string
	Production   Production   // perpetual gains
	Construction Construction // one-time-only gains
	Cost         Cost
	TokenChoice  tokenChoice
	Built        bool
}

func (b *Wonder) getProduction() *Production {
	return &b.Production
}
func (b *Wonder) getConstruction() *Construction {
	return &b.Construction
}
func (b *Wonder) getCost() *Cost {
	return &b.Cost
}

type bonus struct {
	Best   []string
	Coin   int
	Points int
}

type building struct {
	ID           string
	Name         string
	Type         string
	Cost         Cost
	Production   Production   // perpetual gains
	Construction Construction // one-time-only gains
	Bonus        bonus
	Trade        []string
	Linked       string // sfrutta questa catena
	Links        string // crea questa catena
	Points       int
	Science      string
}

func (b *building) getProduction() *Production {
	return &b.Production
}
func (b *building) getConstruction() *Construction {
	return &b.Construction
}
func (b *building) getCost() *Cost {
	return &b.Cost
}

// GenericBuilding is a common interface for wonders and normal buildings
type GenericBuilding interface {
	getProduction() *Production
	getConstruction() *Construction
	getCost() *Cost
}

// GetProduction return production of building or wonder
func GetProduction(g GenericBuilding) *Production {
	return g.getProduction()
}

// GetConstruction return construction of building or wonder
func GetConstruction(g GenericBuilding) *Construction {
	return g.getConstruction()
}

// GetCost return cost of building or wonder
func GetCost(g GenericBuilding) *Cost {
	return g.getCost()
}

type deck struct {
	Age       int
	Buildings []building
}

func (d *deck) removeBuilding(i int) {
	if i < len(d.Buildings) {
		//fmt.Printf("Remove %s %s\n", d.Buildings[i].Type, d.Buildings[i].Name)
		d.Buildings[i] = d.Buildings[len(d.Buildings)-1]
		d.Buildings = d.Buildings[:len(d.Buildings)-1]
	}
}

func (d *deck) prepareBuildings() {
	const replaceCards = 3
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(d.Buildings), func(i, j int) {
		d.Buildings[i], d.Buildings[j] = d.Buildings[j], d.Buildings[i]
	})
	// delete extra cards
	if d.Age == 3 {
		guilds := 0
		standard := 0
		for i, b := range d.Buildings {
			if guilds >= replaceCards && standard >= replaceCards {
				break
			}
			if b.Type == "guild" {
				if guilds < replaceCards {
					d.removeBuilding(i)
					guilds++
				}
			} else {
				if standard < replaceCards {
					d.removeBuilding(i)
					standard++
				}
			}
		}
	} else {
		d.Buildings = d.Buildings[replaceCards:]
	}
}
