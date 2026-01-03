package internal

import (
	"fmt"
	"log"
	"time"
)

const (
	TeamOne = iota
	TeamTwo
)

type Team struct {
	Game    *Game
	ID      string
	Name    string
	Base    string // base location
	Size    int
	Users   []*User
	Minions []*Minion
}

func NewTeam(team string, g *Game) *Team {
	var data MapDataTeam
	for _, t := range g.Map.Data.Teams {
		if t.ID == team {
			data = t
			break
		}
	}
	return &Team{
		Game:  g,
		ID:    data.ID,
		Name:  data.Name,
		Base:  data.Base,
		Size:  data.Size,
		Users: make([]*User, 0),
	}
}

func (t *Team) SpawnSmallMinionWave() {
	log.Println("[INFO] spawning minion wave")
	go func() {
		for i := 0; i < t.Game.Map.Data.Minions.SmallCount; i++ {
			for l := 0; l < 3; l++ {
				newMinion := NewMinion(MinionTypeSmall, l, t, t.Game)
				log.Println("[INFO] creating minion", t.Name, l, newMinion.Coords)
				t.Minions = append(t.Minions, newMinion)
			}
			c := colors[t.ID].Sprint(t.Name)
			sendToAll(t.Game.GetUsersInRoom(t.Game.Map.RoomIDToCoords(t.Base)), fmt.Sprintf("A %s minion wave has arrived.\n", c))
			time.Sleep(1 * time.Second)
		}
	}()
}

func (t *Team) Update(dt float64) {
	for _, m := range t.Minions {
		m.Update(dt)
	}
}
