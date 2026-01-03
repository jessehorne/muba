package internal

import (
	"fmt"
	"log"
	"time"
)

const (
	GameStageWaitingForPlayers = iota
	GameStageCountdown
	GameStageRunning
	GameStageEnded
)

type Game struct {
	Time    float32 // game runtime in seconds
	Started bool    // if the game has been started
	Teams   map[string]*Team
	Server  *Server
	Map     *Map
	Running bool
	dt      float64

	StageTimerCounter float64
	Stage             int
}

func NewGame(s *Server, mapPath string) (*Game, error) {
	mapData, err := NewMap(mapPath)
	if err != nil {
		return nil, err
	}

	g := &Game{
		Server:  s,
		Map:     mapData,
		Running: true,
		Stage:   GameStageWaitingForPlayers,
	}

	g.LoadMap()

	return g, nil
}

func (g *Game) LoadMap() {
	// init things
	initColors()
	initGamestats()

	// init teams
	g.Teams = make(map[string]*Team)
	for _, t := range g.Map.Data.Teams {
		g.Teams[t.ID] = NewTeam(t.ID, g)
	}

	// init rooms
	roomCounter := 0
	for y, row := range g.Map.Data.Map {
		for x, roomID := range row {
			newRoom := NewRoom(
				g, roomID, NewVector2(x, y), "Tower "+roomID, "A tower that attacks the enemy.")
			g.Map.Rooms[y][x] = newRoom
			roomCounter++
		}
	}
	log.Printf("[INFO] created %d rooms\n", roomCounter)

	// set up bases
	// todo

	// load towers
	towersCounter := 0
	for _, t := range g.Map.Data.Towers.Red {
		coords := g.Map.RoomIDToCoords(t)
		room := g.Map.CoordsToRoom(coords)
		if room != nil {
			room.Towers = append(room.Towers, NewTower(g, room, coords, g.Teams["red"]))
			towersCounter++
		}
	}
	for _, t := range g.Map.Data.Towers.Blue {
		coords := g.Map.RoomIDToCoords(t)
		room := g.Map.CoordsToRoom(coords)
		if room != nil {
			room.Towers = append(room.Towers, NewTower(g, room, coords, g.Teams["blue"]))
			towersCounter++
		}
	}
	log.Printf("[INFO] created %d towers\n", towersCounter)

	// set up minion spawning
	// todo

	// set up camps
	// todo

	// set up bosses
	// todo

	// set up vision camps
	// todo
}

func (g *Game) SetPlayersStartingPositions() {
	log.Println("[INFO] setting player starting positions")
	for _, t := range g.Teams {
		for _, u := range t.Users {
			u.Coords = g.Map.RoomIDToCoords(t.ID)
		}
	}
}

func (g *Game) StartGame() {
	g.StartMinionSpawnerHandler()

	oldTime := time.Now()
	counter := 0.0
	g.Running = true
	for g.Running {
		newTime := time.Now()
		g.dt = newTime.Sub(oldTime).Seconds()
		counter += g.dt
		if counter >= 1/g.Map.Data.TickRate {
			g.UpdateRooms(counter)
			counter = 0
		}

		oldTime = newTime
	}
}

func (g *Game) UpdateRooms(dt float64) {
	for y := 0; y < len(g.Map.Data.Map); y++ {
		for x := 0; x < len(g.Map.Data.Map[y]); x++ {
			r := g.Map.Rooms[y][x]
			r.Update(dt)
		}
	}
}

func (g *Game) StartMinionSpawnerHandler() {
	log.Println("[INFO] starting minion spawner handler")
	go func() {
		//ticker := time.NewTicker(time.Duration(g.Map.Data.Minions.RespawnTime) * time.Second)
		//for _ = range ticker.C {
		//	for _, t := range g.Teams {
		//		t.SpawnSmallMinionWave()
		//	}
		//}
		for _, t := range g.Teams {
			t.SpawnSmallMinionWave()
		}
	}()
}

func (g *Game) StartCountdownStage() {
	sendToAll(g.Server.Users, "All players are ready. Let the countdown begin.")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	secondsCounter := 3

	for range ticker.C {
		addon := "..."
		if secondsCounter == 0 {
			addon = "! GO!!!"
		}

		if secondsCounter == 0 {
			g.StartRunningStage()
			return
		} else {
			sendToAll(g.Server.Users, fmt.Sprintf("%d%s\n", secondsCounter, addon))
		}
		secondsCounter--
	}
}

func (g *Game) StartRunningStage() {
	g.Stage = GameStageRunning
	g.Running = true
	go g.StartGame()
	sendToAll(g.Server.Users, "The game has begun!\n")
}

func (g *Game) CheckIfReadyToStart() {
	allReady := true
	for _, t := range g.Teams {
		if len(t.Users) != t.Size {
			allReady = false
		}
		for _, user := range t.Users {
			if !user.Ready {
				allReady = false
			}
		}
	}
	if allReady {
		log.Println("[INFO] All players ready...starting countdown!")
		g.StartCountdownStage()
	}
}

func (g *Game) GetUsersInRoom(coords Vector2) map[string]*User {
	usersInRoom := make(map[string]*User)
	for _, u := range g.Server.Users {
		if u.Coords == coords {
			usersInRoom[u.ID] = u
		}
	}
	return usersInRoom
}
