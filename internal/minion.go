package internal

import (
	"fmt"
	"math/rand"
)

const (
	MinionTypeSmall int = iota
	MinionTypeBig
)

type MinionWave struct {
	Team          *Team
	Game          *Game
	Coords        Vector2
	Velocity      Vector2
	Type          int
	Lane          int
	Health        int
	CurrentHealth int
	Speed         int
	CoinGive      int
	Damage        int
	AttackSpeed   int

	walkTimer float64

	InCombat    bool
	combatTimer float64

	targetMinion *MinionWave
	targetUser   *User
	targetTower  *Tower

	lastHitter *User
}

func NewMinionWave(minionType int, lane int, t *Team, g *Game) *MinionWave {
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

	return &MinionWave{
		Team:          t,
		Game:          g,
		Type:          minionType,
		Lane:          lane,
		Coords:        g.Map.RoomIDToCoords(t.Base),
		Velocity:      vel,
		Health:        g.Map.Data.Minions.Health,
		CurrentHealth: g.Map.Data.Minions.Health,
		Speed:         g.Map.Data.Minions.Speed,
		CoinGive:      g.Map.Data.Minions.Gives.LastHitter.Coin,
		Damage:        damage,
		AttackSpeed:   atkSpeed,
	}
}

func (m *MinionWave) Update(dt float64) {
	// check if enemy minions are in the room, if so target
	if m.targetMinion == nil {
		var enemyMinions []*MinionWave
		for _, m2 := range m.Room().Minions {
			if m2.Team.ID != m.Team.ID {
				enemyMinions = append(enemyMinions, m2)
			}
		}
		if len(enemyMinions) > 0 {
			i := rand.Intn(len(enemyMinions))
			m.targetMinion = enemyMinions[i]
		}
	}

	// check if enemy players in room if not already attacking one, if so target at random
	if m.targetMinion == nil && m.targetUser == nil {
		enemyUsers := []*User{}
		for _, u := range m.Room().Users {
			if u.Team != m.Team {
				enemyUsers = append(enemyUsers, u)
			}
		}

		if len(enemyUsers) > 0 {
			i := rand.Intn(len(enemyUsers))
			m.targetUser = enemyUsers[i]
		}
	}

	// check if enemy towers in room if the user isn't attacking a user
	if m.targetMinion == nil && m.targetUser == nil && m.targetTower == nil {
		for _, t := range m.Room().Towers {
			if t.Team != m.Team {
				m.targetTower = t
				break
			}
		}
	}

	// determine if in combat
	if m.targetMinion != nil || m.targetUser != nil || m.targetTower != nil {
		m.InCombat = true
	} else {
		m.InCombat = false
	}

	// do damage
	if m.InCombat {
		m.combatTimer += dt

		if m.combatTimer >= float64(m.AttackSpeed) {
			m.combatTimer = 0

			// try to do damage
			if m.targetMinion != nil {
				m.AttackMinion(m.targetMinion)
			} else if m.targetUser != nil {
				m.AttackPlayer(m.targetUser)
			} else if m.targetTower != nil {
				m.AttackTower(m.targetTower)
			}
		}
	} else {
		m.walkTimer += dt
		if m.walkTimer >= float64(m.Speed) {
			m.walkTimer = 0
			m.MoveToNextRoom()
		}
	}
}

func (m *MinionWave) AttackMinion(m2 *MinionWave) {
	m2.CurrentHealth -= m.Damage
	sendToAll(m.Game.GetUsersInRoom(m.Coords), fmt.Sprintf("A %s has done %d damage to a %s.\n",
		m.GetColoredNameAndHealth(), m.Damage, m2.GetColoredNameAndHealth()))
	m2.lastHitter = nil // set to nil because a minion hit last, not a user
}

func (m *MinionWave) AttackPlayer(u *User) {
	u.CurrentHealth -= m.Damage
	c := colors[m.Team.ID].Sprint(m.Team.Name)
	sendToOne(u.conn, fmt.Sprintf("You've been struck by a %s minion wave for %d damage!\n", c, m.Damage))
	sendToAllExcept(u.conn, m.Game.GetUsersInRoom(m.Coords), fmt.Sprintf("%s has been struck by a %s minion wave for %d damage!\n",
		u.Username, c, m.Damage))
}

func (m *MinionWave) AttackTower(t *Tower) {
	t.CurrentHealth -= m.Damage
	c := colors[m.Team.ID].Sprint(m.Team.Name)
	sendToAll(m.Game.GetUsersInRoom(m.Coords), fmt.Sprintf("The tower was hit by a %s minion wave for %d damage!\n",
		c, m.Damage))
}

func (m *MinionWave) MoveToNextRoom() {
	// check if they can go to next room, if not, update velocity
	nextRoomCoords := m.Coords.Add(m.Velocity)
	oldRoomCoords := m.Coords
	nextRoom := m.Game.Map.CoordsToRoomID(nextRoomCoords)
	if nextRoom != "" {
		// remove from previous room list of minions and add to next
		mi := -1
		for i, mm := range m.Room().Minions {
			if m == mm {
				mi = i
				break
			}
		}
		if mi != -1 {
			r := m.Room()
			r.Minions = append(r.Minions[:mi], r.Minions[mi+1:]...)
		}
		m.Coords = nextRoomCoords
		r := m.Room()
		r.Minions = append(r.Minions, m)
		sendToAll(m.Game.GetUsersInRoom(m.Coords), fmt.Sprintf("A %s minion wave approaches.\n", colors[m.Team.ID].Sprint(m.Team.Name)))
		sendToAll(m.Game.GetUsersInRoom(oldRoomCoords), fmt.Sprintf("A %s minion wave exits.\n", colors[m.Team.ID].Sprint(m.Team.Name)))
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
					m.Velocity = NewVector2(0, -1)
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

func (m *MinionWave) Room() *Room {
	return m.Game.Map.CoordsToRoom(m.Coords)
}

func (m *MinionWave) FormattedHealth() string {
	c := colors["green"]
	perc := float64(m.CurrentHealth) / float64(m.Health)
	if perc < 0.2 {
		c = colors["red"]
	} else if perc < 0.7 {
		c = colors["yellow"]
	}
	msg := c.Sprint(m.CurrentHealth) + "/" + colors["green"].Sprint(m.Health)
	return msg
}

func (m *MinionWave) GetColoredNameAndHealth() string {
	c := colors[m.Team.ID].Sprint("minion") + " wave"
	return fmt.Sprintf("%s (%s)", c, m.FormattedHealth())
}
