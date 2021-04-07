package game

import (
	"math/rand"
	"time"
)

type wonder struct {
	ID           string
	Name         string
	Production   production   // perpetual gains
	Construction construction // one-time-only gains
	Cost         cost
	TokenChoice  tokenChoice
}

func (b *wonder) getProduction() *production {
	return &b.Production
}
func (b *wonder) getConstruction() *construction {
	return &b.Construction
}
func (b *wonder) getCost() *cost {
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
	Cost         cost
	Production   production   // perpetual gains
	Construction construction // one-time-only gains
	Bonus        bonus
	Trade        []string
	CreationLink string
	CreatedLink  string
	Points       int
	Science      string
}

func (b *building) getProduction() *production {
	return &b.Production
}
func (b *building) getConstruction() *construction {
	return &b.Construction
}
func (b *building) getCost() *cost {
	return &b.Cost
}

type genericBuilding interface {
	getProduction() *production
	getConstruction() *construction
	getCost() *cost
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
