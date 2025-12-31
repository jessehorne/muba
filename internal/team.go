package internal

const (
	TeamOne = iota
	TeamTwo
)

type Team struct {
	Name    string
	Players []*Player
}

func NewTeam(name string) *Team {
	return &Team{
		Name:    name,
		Players: []*Player{},
	}
}
