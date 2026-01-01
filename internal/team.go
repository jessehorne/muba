package internal

const (
	TeamOne = iota
	TeamTwo
)

type Team struct {
	ID    string
	Name  string
	Base  string // base location
	Size  int
	Users []*User
}

func NewTeam(data MapDataTeam) *Team {
	return &Team{
		ID:    data.ID,
		Name:  data.Name,
		Base:  data.Base,
		Size:  data.Size,
		Users: make([]*User, 0),
	}
}
