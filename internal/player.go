package internal

type Player struct {
	Coords    Vector2
	GameStats *GameStats
	Team      *Team
}

func NewPlayer(champType int, t *Team) *Player {
	return &Player{
		Coords:    NewVector2(0, 0),
		GameStats: NewGameStats(champType),
		Team:      t,
	}
}
