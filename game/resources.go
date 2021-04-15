package game

// Cost is the resources needed for build
type Cost struct {
	Coins   int
	Wood    int
	Clay    int
	Stone   int
	Glass   int
	Papyrus int
}

// Production is what a building provides
type Production struct {
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

// ToMap easy maps label and values
func (p *Production) ToMap() map[string]int {
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

// Construction is what you gain after built
type Construction struct {
	Points       int
	Turn         bool
	Coins        int
	CoinsRemoved int
	Discard      string
	Tokens       tokenChoice
	Shield       int
	Production   Production
	ForEach      forEach
}

type token struct {
	ID          string
	Name        string
	Description string
}
