package main

type wonder struct {
	ID           string
	Name         string
	Production   production
	Construction construction
	Cost         cost
	TokenChoice  tokenChoice
}

type bonus struct {
	Best   []string
	Coin   resource
	Points int
}

type building struct {
	ID           string
	Name         string
	Type         string
	Cost         cost
	Production   production
	Costruction  construction
	Bonus        bonus
	CreationLink string
	CreatedLink  string
	Points       int
	Science      string
}

type deck struct {
	Age       int
	Buildings []building
}
