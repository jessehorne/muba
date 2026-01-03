package internal

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
