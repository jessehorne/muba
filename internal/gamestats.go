package internal

const (
	ChampTypeFighter = iota
	ChampTypeWizard
	ChampTypeArcher
	ChampTypeNinja
)

var (
	ChampStringToType map[string]int
	ChampTypeToString map[int]string
	ChampTypeColored  map[string]string

	GameStatsFighter = GameStats{
		Strength:      10,
		Defense:       10,
		MoveSpeed:     7,
		AttackSpeed:   4,
		Aura:          1,
		Health:        700,
		HealPerSecond: 1,
		Gold:          100,
		Experience:    0,
		Level:         1,
	}

	GameStatsWizard = GameStats{
		Strength:      1,
		Defense:       3,
		MoveSpeed:     5,
		AttackSpeed:   4,
		Aura:          10,
		Health:        500,
		HealPerSecond: 1,
		Gold:          100,
		Experience:    0,
		Level:         1,
	}

	GameStatsArcher = GameStats{
		Strength:      3,
		Defense:       3,
		MoveSpeed:     4,
		AttackSpeed:   2,
		Aura:          0,
		Health:        400,
		HealPerSecond: 1,
		Gold:          100,
		Experience:    0,
		Level:         1,
	}

	GameStatsNinja = GameStats{
		Strength:      5,
		Defense:       1,
		MoveSpeed:     3,
		AttackSpeed:   3,
		Aura:          0,
		Health:        300,
		HealPerSecond: 1,
		Gold:          100,
		Experience:    0,
		Level:         1,
	}
)

func initGamestats() {
	ChampStringToType = map[string]int{
		"fighter": ChampTypeFighter,
		"wizard":  ChampTypeWizard,
		"archer":  ChampTypeArcher,
		"ninja":   ChampTypeNinja,
	}

	ChampTypeToString = map[int]string{
		ChampTypeFighter: "fighter",
		ChampTypeWizard:  "wizard",
		ChampTypeArcher:  "archer",
		ChampTypeNinja:   "ninja",
	}

	ChampTypeColored = map[string]string{
		"fighter": colors["purple"].Sprint("fighter"),
		"wizard":  colors["cyan"].Sprint("wizard"),
		"archer":  colors["green"].Sprint("archer"),
		"ninja":   colors["yellow"].Sprint("ninja"),
	}
}

type GameStats struct {
	Strength      int // how hard you hit with melee/weapons
	Defense       int // how much you can take a punch
	MoveSpeed     int // how fast you move
	AttackSpeed   int // how fast you attack
	Aura          int // how good you are at magick
	Health        int // your hit points
	HealPerSecond int // how much you heal per second
	Gold          int // gold aka currency
	Experience    int // total experience
	Level         int // real level
}

func NewGameStats(champType int) *GameStats {
	switch champType {
	case ChampTypeFighter:
		return &GameStatsFighter
	case ChampTypeWizard:
		return &GameStatsWizard
	case ChampTypeArcher:
		return &GameStatsArcher
	case ChampTypeNinja:
		return &GameStatsNinja
	default:
		return nil
	}
}
