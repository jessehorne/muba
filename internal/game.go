package internal

type Game struct {
	Time    float32 // game runtime in seconds
	Started bool    // if the game has been started
	Teams   map[string]*Team
	Players map[string]*Player
}

func NewGame() *Game {
	return &Game{
		Teams: map[string]*Team{
			"red":  NewTeam("red"),
			"blue": NewTeam("blue"),
		},
		Players: make(map[string]*Player),
	}
}

func (g *Game) AddPlayer(p *Player) {
	_, ok := g.Players[p.ID]
	if !ok {
		g.Players[p.ID] = p
	}
}

func (g *Game) RemovePlayer(p *Player) {
	_, ok := g.Players[p.ID]
	if !ok {
		delete(g.Players, p.ID)
	}
}
