package internal

import (
	"fmt"
	"log"
)

const (
	TeamOne = iota
	TeamTwo
)

type Team struct {
	Game              *Game
	ID                string
	Name              string
	Base              string // base location
	Size              int
	Users             []*User
	Minions           []*MinionWave
	CurrentMinionWave int
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
	t.CurrentMinionWave++
	log.Printf("[INFO] spawning minion wave #%d\n", t.CurrentMinionWave)
	go func() {
		for l := 0; l < 3; l++ {
			newMinion := NewMinionWave(MinionTypeSmall, l, t, t.Game)
			t.Minions = append(t.Minions, newMinion)
			r := t.Game.Map.CoordsToRoom(t.Game.Map.RoomIDToCoords(t.Base))
			r.Minions = append(r.Minions, newMinion)
			c := colors[t.ID].Sprint(t.Name)
			addon := "heading top"
			if l == LaneMid {
				addon = "heading mid"
			} else if l == LaneBottom {
				addon = "heading bot"
			}
			sendToAll(t.Game.GetUsersInRoom(t.Game.Map.RoomIDToCoords(t.Base)), fmt.Sprintf("A %s minion wave has arrived %s.\n",
				c, addon))
		}
	}()
}

func (t *Team) Update(dt float64) {
	for _, m := range t.Minions {
		m.Update(dt)
	}
}

func (t *Team) GetEnemyTeam() *Team {
	for _, tt := range t.Game.Teams {
		if tt != t {
			return tt
		}
	}
	return nil
}

func (t *Team) KillMinion(m *MinionWave) {
	minionIndex := -1
	for i, m2 := range t.Minions {
		if m == m2 {
			minionIndex = i
			break
		}
	}
	if minionIndex == -1 {
		t.Minions = append(t.Minions[:minionIndex], t.Minions[minionIndex+1:]...)
	}
}
