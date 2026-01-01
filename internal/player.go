package internal

import "github.com/google/uuid"

type Player struct {
	ID        string
	Coords    Vector2
	GameStats *GameStats
	Team      *Team
}

func NewPlayer(champType int, t *Team) *Player {
	return &Player{
		ID:        uuid.New().String(),
		Coords:    NewVector2(0, 0),
		GameStats: NewGameStats(champType),
		Team:      t,
	}
}
