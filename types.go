type resource int
const (
	wood resource = iota
	clay
	stone
	glass
	papyrus
	shield
)

type wonder struct {
	name String
	cost map[resource]int
	effects map[effect]int
	points int
	shields int
	built bool
	coins_gained coin
	coins_removed coin
	new_turn bool
}

type scientific_symbol int
const (
	astro scientific_symbol = iota
	wheel
	sundial
	mortar
	pendulum
	ink
)

type building_type int
const (
	raw building_type = iota
	manufactured
	civilian
	scientific
	commercial
	military
)

type guild_type int
const (
	builders guild_type = iota
	moneylenders
	scientists
	shipowners
	traders
	magistrates
	tacticians
)

type linking_symbol int
const (
	stable condition_type = iota
	archery
	castle
	helm
	vase
	barrel
	mask
	theater
	sun
	drop
	column
	moon
	sword
	music
	gear
	book
	lamp
)

type building struct {
	age int
	name String
	guild guild_type
	cost map[material]int
	required_link linking_symbol
	production map[material]
	created_link linking_symbol
	color building_type
	points int
	visible bool
	blocked_by []*building
}

type coin int

type progress_token int
const (
	agriculture progress_token = iota
	architecture
	economy
	law
	masonry
	mathematics
	philosophy
	strategy
	theology
	urbanism
)

type player struct {
	name String
	active bool
	coins coin
	wonders []wonder
	buildings []building
	progress_tokens []progress_token
	production map[resource]int
	price map[resource]int
}

type game struct {
	progress_tokens []progress_token
	buildable_buildings []*building
	discarded_buildings []*building
}
