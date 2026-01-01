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

	g := &Game{
		Players: make(map[string]*Player),
		Map:     mapData,
	}

	g.LoadMap()

	return g, nil
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

func (g *Game) LoadMap() {
	// init teams
	g.Teams = make(map[string]*Team)
	for _, t := range g.Map.Data.Teams {
		g.Teams[t.ID] = NewTeam(t)
	}

	// set up base
	// todo

	// load towers
	// todo

	// set up minion spawning
	// todo

	// set up camps
	// todo

	// set up bosses
	// todo

	// set up vision camps
	// todo
}
