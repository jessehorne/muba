package internal

const (
	ChampTypeFighter = iota
	ChampTypeWizard
	ChampTypeArcher
	ChampTypeNinja
)

var (
	ChampStringToType = map[string]int{
		"fighter": ChampTypeFighter,
		"wizard":  ChampTypeWizard,
		"archer":  ChampTypeArcher,
		"ninja":   ChampTypeNinja,
	}
)

var (
	GameStatsFighter = GameStats{
		Strength:    0,
		Defense:     0,
		MoveSpeed:   0,
		AttackSpeed: 0,
		Aura:        0,
		Health:      0,
		Gold:        0,
		Experience:  0,
		Level:       0,
	}

	GameStatsWizard = GameStats{
		Strength:    0,
		Defense:     0,
		MoveSpeed:   0,
		AttackSpeed: 0,
		Aura:        0,
		Health:      0,
		Gold:        0,
		Experience:  0,
		Level:       0,
	}

	GameStatsArcher = GameStats{
		Strength:    0,
		Defense:     0,
		MoveSpeed:   0,
		AttackSpeed: 0,
		Aura:        0,
		Health:      0,
		Gold:        0,
		Experience:  0,
		Level:       0,
	}

	GameStatsNinja = GameStats{
		Strength:    0,
		Defense:     0,
		MoveSpeed:   0,
		AttackSpeed: 0,
		Aura:        0,
		Health:      0,
		Gold:        0,
		Experience:  0,
		Level:       0,
	}
)

type GameStats struct {
	Strength    int // how hard you hit with melee/weapons
	Defense     int // how much you can take a punch
	MoveSpeed   int // how fast you move
	AttackSpeed int // how fast you attack
	Aura        int // how good you are at magick
	Health      int // your hit points
	Gold        int // gold aka currency
	Experience  int // total experience
	Level       int // real level
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
