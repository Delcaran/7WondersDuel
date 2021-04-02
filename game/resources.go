package game

type cost struct {
	Coins   int
	Wood    int
	Clay    int
	Stone   int
	Glass   int
	Papyrus int
}

type production struct {
	Coins   int
	Wood    int
	Clay    int
	Stone   int
	Glass   int
	Papyrus int
	Shield  int
	Choice  bool
}

const (
	coins = iota
	wood
	clay
	stone
	glass
	papyrus
	shield
)

func (p *production) ToMap() map[string]int {
	var m = map[string]int{
		"Coins":   p.Coins,
		"Wood":    p.Wood,
		"Clay":    p.Clay,
		"Stone":   p.Stone,
		"Glass":   p.Glass,
		"Papyrus": p.Papyrus,
		"Shield":  p.Shield,
	}
	return m
}

type forEach struct {
	Building string
	Coins    int
}

type construction struct {
	Points       int
	Turn         bool
	Coins        int
	CoinsRemoved int
	Discard      string
	Tokens       tokenChoice
	Shield       int
	Production   production
	ForEach      forEach
}

type token struct {
	ID          string
	Name        string
	Description string
}
