package internal

import "time"

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
		Attacks: map[string]Attack{
			"1": {
				Name:        "Swing",
				Description: "A long and slow swing of the sword.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level*u.GameStats.Strength*10 + 50
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 5,
			},
			"2": {
				Name:        "Slash",
				Description: "A short but strong slash with the sword.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level*u.GameStats.Strength*10 + 10
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 3,
			},
			"3": {
				Name:        "Stab",
				Description: "A quick lunge with the sword.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level * u.GameStats.Strength * 10
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 3,
			},
			"4": {
				Name:        "Triple Swing Beatdown",
				Description: "Three fast swinging attacks that do a lot of damage.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level * u.GameStats.Strength * 20
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 3,
			},
		},
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
		Attacks: map[string]Attack{
			"1": {
				Name:        "Fireball",
				Description: "A large blast of fire for high damage.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level*u.GameStats.Aura*10 + 50
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 5,
			},
			"2": {
				Name:        "Sludge",
				Description: "A blast of sludge that slows an enemy player for 3 seconds.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						hs[0].Slow(u.GameStats.Level)
					}
				},
				Cooldown: 3,
			},
			"3": {
				Name:        "Heal",
				Description: "Heal yourself and friendly players.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						hs[0].DoDamage(-u.GameStats.Level * u.GameStats.Aura * 10)
					}
				},
				Cooldown: 3,
			},
			"4": {
				Name:        "Call to Xan'Roth",
				Description: "Call upon Xan'Roth to cast down lightning from the ethereal realm on an opponent.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level * u.GameStats.Aura * 20
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 3,
			},
		},
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
		Level:         1, Attacks: map[string]Attack{
			"1": {
				Name:        "Three Arrow Blast",
				Description: "Three quick arrows to an opponent.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level*u.GameStats.Strength*10 + 60
						hs[0].DoDamage(dmg)
					}
				},
				Cooldown: 2,
			},
			"2": {
				Name:        "Heavy Arrow",
				Description: "A slow but heavy hitting arrow that knocks an opponent out temporarily.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						hs[0].KnockOut(u.GameStats.Level)
					}
				},
				Cooldown: 5,
			},
			"3": {
				Name:        "Fire Arrow",
				Description: "Does damage to an opponent over time with fire!",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						go func() {
							dmg := u.GameStats.Level * u.GameStats.Strength
							for i := 0; i < u.GameStats.Level; i++ {
								hs[0].DoDamage(dmg)
								time.Sleep(1 * time.Second)
							}
						}()
					}
				},
				Cooldown: 3,
			},
			"4": {
				Name:        "Explosive Arrow",
				Description: "Launch an explosive arrow at an enemy doing damage to every opponent within range.",
				Do: func(u *User, hs []Hurtable) {
					if hs != nil {
						dmg := u.GameStats.Level * u.GameStats.Strength * 20
						for _, h := range hs {
							h.DoDamage(dmg)
						}
					}
				},
				Cooldown: 5,
			},
		},
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

type Attack struct {
	Name        string
	Description string
	Do          func(u *User, hs []Hurtable)
	Cooldown    int
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
	Attacks       map[string]Attack
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
