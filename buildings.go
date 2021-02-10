package main

type wonder struct {
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
	Name         string
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
	Age          int
	Raw          []building
	Manufactured []building
	Military     []building
	Scientific   []building
	Civilian     []building
	Commercial   []building
	Guild        []building
}
