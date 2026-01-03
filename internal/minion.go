package internal

import (
	"fmt"
)

const (
	MinionTypeSmall int = iota
	MinionTypeBig
)

type Minion struct {
	Team        *Team
	Game        *Game
	Coords      Vector2
	Velocity    Vector2
	Type        int
	Lane        int
	Health      int
	Speed       int
	CoinGive    int
	Damage      int
	AttackSpeed int

	walkTimer float64

	InCombat bool
}

func NewMinion(minionType int, lane int, t *Team, g *Game) *Minion {
	damage := g.Map.Data.Minions.SmallDamage
	if minionType == MinionTypeBig {
		damage = g.Map.Data.Minions.BigDamage
	}

	atkSpeed := g.Map.Data.Minions.SmallAttackSpeed
	if minionType == MinionTypeBig {
		atkSpeed = g.Map.Data.Minions.BigAttackSpeed
	}

	var vel Vector2
	if t.ID == "red" {
		if lane == LaneTop {
			vel = NewVector2(0, -1)
		} else if lane == LaneMid {
			vel = NewVector2(1, -1)
		} else if lane == LaneBottom {
			vel = NewVector2(1, 0)
		}
	} else if t.ID == "blue" {
		if lane == LaneTop {
			vel = NewVector2(-1, 0)
		} else if lane == LaneMid {
			vel = NewVector2(-1, 1)
		} else if lane == LaneBottom {
			vel = NewVector2(0, 1)
		}
	}

	return &Minion{
		Team:        t,
		Game:        g,
		Type:        minionType,
		Lane:        lane,
		Coords:      g.Map.RoomIDToCoords(t.Base),
		Velocity:    vel,
		Health:      g.Map.Data.Minions.Health,
		Speed:       g.Map.Data.Minions.Speed,
		CoinGive:    g.Map.Data.Minions.Gives.LastHitter.Coin,
		Damage:      damage,
		AttackSpeed: atkSpeed,
	}
}

func (m *Minion) Update(dt float64) {
	m.walkTimer += dt
	if m.walkTimer >= float64(m.Speed) {
		m.walkTimer = 0

		if !m.InCombat {
			m.MoveToNextRoom()
		}
	}
}

func (m *Minion) MoveToNextRoom() {
	// check if they can go to next room, if not, update velocity
	nextRoomCoords := m.Coords.Add(m.Velocity)
	oldRoomCoords := m.Coords
	nextRoom := m.Game.Map.CoordsToRoomID(nextRoomCoords)
	if nextRoom != "" {
		m.Coords = nextRoomCoords
		sendToAll(m.Game.GetUsersInRoom(m.Coords), fmt.Sprintf("A %s minion approaches.\n", colors[m.Team.ID].Sprint(m.Team.Name)))
		sendToAll(m.Game.GetUsersInRoom(oldRoomCoords), fmt.Sprintf("A %s minion exits.\n", colors[m.Team.ID].Sprint(m.Team.Name)))
	} else {
		// determine if we need to update velocity because minion hit a wall
		if m.Lane == LaneTop {
			if m.Team.ID == "red" {
				if m.Coords == NewVector2(0, 0) {
					m.Velocity = NewVector2(1, 0)
					m.MoveToNextRoom()
				}
			} else if m.Team.ID == "blue" {
				if m.Coords == NewVector2(0, 0) {
					m.Velocity = NewVector2(0, 1)
					m.MoveToNextRoom()
				}
			}
		} else if m.Lane == LaneBottom {
			if m.Team.ID == "red" {
				if m.Coords == NewVector2(6, 6) {
					m.Velocity = NewVector2(0, 1)
					m.MoveToNextRoom()
				}
			} else if m.Team.ID == "blue" {
				if m.Coords == NewVector2(6, 6) {
					m.Velocity = NewVector2(-1, 0)
					m.MoveToNextRoom()
				}
			}
		}
	}
}
