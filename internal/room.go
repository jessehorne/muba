package internal

type Room struct {
	Game        *Game
	ID          string
	Coords      Vector2
	Title       string
	Description string
	Users       []*User
	Minions     []*Minion
	Towers      []*Tower
}

func NewRoom(g *Game, id string, coords Vector2, title, desc string) *Room {
	return &Room{
		Game:        g,
		ID:          id,
		Coords:      coords,
		Title:       title,
		Description: desc,
		Users:       make([]*User, 0),
		Minions:     make([]*Minion, 0),
		Towers:      make([]*Tower, 0),
	}
}
