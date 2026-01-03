package internal

import (
	"fmt"

	"github.com/fatih/color"
)

type Room struct {
	Game        *Game
	ID          string
	Coords      Vector2
	Title       string
	Description string
	Users       []*User
	Minions     []*MinionWave
	Towers      []*Tower

	minionsFightingTimer float64
}

func NewRoom(g *Game, id string, coords Vector2, title, desc string) *Room {
	return &Room{
		Game:        g,
		ID:          id,
		Coords:      coords,
		Title:       title,
		Description: desc,
		Users:       make([]*User, 0),
		Minions:     make([]*MinionWave, 0),
		Towers:      make([]*Tower, 0),
	}
}

func (r *Room) Update(dt float64) {
	// update minions
	for _, m := range r.Minions {
		m.Update(dt)

		// check if any need to die
		if m.CurrentHealth <= 0 {
			c := colors[m.Team.ID].Sprint(m.Team.Name)
			sendToAll(r.Game.GetUsersInRoom(r.Coords), fmt.Sprintf("%s has come for a %s minion wave!\n", c,
				color.New(color.BgWhite).Add(color.FgBlack).Sprint("Death")))
			r.KillMinion(m)
			for _, m2 := range r.Minions {
				if m2.targetMinion == m {
					m2.targetMinion = nil
				}
			}
		}
	}

	// tell us if minions are fighting
	hasRed := false
	hasBlue := false
	for _, m := range r.Minions {
		if m.Team.ID == "red" {
			hasRed = true
		} else if m.Team.ID == "blue" {
			hasBlue = true
		}
	}
	if hasRed && hasBlue {
		r.minionsFightingTimer += dt

		if r.minionsFightingTimer > 3 {
			r.minionsFightingTimer = 0
			sendToAll(r.Game.GetUsersInRoom(r.Coords), "The minions from both sides are fighting each other...")
		}
	}
}

func (r *Room) KillMinion(m *MinionWave) {
	minionIndex := -1
	for i, mm := range r.Minions {
		if m == mm {
			minionIndex = i
			break
		}
	}
	if minionIndex != -1 {
		r.Minions = append(r.Minions[:minionIndex], r.Minions[minionIndex+1:]...)
	}
	m.Team.KillMinion(m)
}
