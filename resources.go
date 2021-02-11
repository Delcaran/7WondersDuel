package main

type resource int

type cost struct {
	Coins   resource
	Wood    resource
	Clay    resource
	Stone   resource
	Glass   resource
	Papyrus resource
}

type production struct {
	Coins   resource
	Wood    resource
	Clay    resource
	Stone   resource
	Glass   resource
	Papyrus resource
	Shield  resource
	Choice  bool
}

type forEach struct {
	Building string
	Coins    resource
}

type construction struct {
	Points       int
	Turn         bool
	Coins        resource
	CoinsRemoved resource
	Discard      string
	Tokens       tokenChoice
	Shield       resource
	Production   production
	ForEach      forEach
}

type token struct {
	ID          string
	Name        string
	Description string
}
