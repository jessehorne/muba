package internal

const (
	TeamOne = iota
	TeamTwo
)

type Team struct {
	Name  string
	Users []*User
}

func NewTeam(name string) *Team {
	return &Team{
		Name:  name,
		Users: make([]*User, 0),
	}
}
