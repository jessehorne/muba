package internal

import "fmt"

type Tower struct {
	Game          *Game
	Room          *Room
	Coords        Vector2
	Team          *Team
	Health        int
	CurrentHealth int
	AttackSpeed   int
	Damage        int
	Gives         MapDataGives
}

func NewTower(g *Game, r *Room, coords Vector2, t *Team) *Tower {
	return &Tower{
		Game:          g,
		Room:          r,
		Coords:        coords,
		Team:          t,
		Health:        g.Map.Data.Towers.Health,
		CurrentHealth: g.Map.Data.Towers.Health,
		AttackSpeed:   g.Map.Data.Towers.AttackSpeed,
		Damage:        g.Map.Data.Towers.Damage,
		Gives:         g.Map.Data.Towers.Gives,
	}
}

func (t *Tower) FormattedHealth() string {
	c := colors["green"]
	perc := float64(t.CurrentHealth) / float64(t.Health)
	if perc < 0.2 {
		c = colors["red"]
	} else if perc < 0.7 {
		c = colors["yellow"]
	}
	msg := c.Sprint(t.CurrentHealth) + "/" + colors["green"].Sprint(t.Health)
	return msg
}

func (t *Tower) GetColoredNameAndHealth() string {
	c := colors[t.Team.ID].Sprint("tower")
	return fmt.Sprintf("%s (%s)", c, t.FormattedHealth())
}
