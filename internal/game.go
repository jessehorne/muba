package internal

type Game struct {
	Time    float32 // game runtime in seconds
	Started bool    // if the game has been started
	Teams   map[string]*Team
	Players map[string]*Player
	Map     *Map
}

func NewGame(mapPath string) (*Game, error) {
	mapData, err := NewMap(mapPath)
	if err != nil {
		return nil, err
	}

	return &Game{
		Teams: map[string]*Team{
			"red":  NewTeam("red"),
			"blue": NewTeam("blue"),
		},
		Players: make(map[string]*Player),
		Map:     mapData,
	}, nil
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
