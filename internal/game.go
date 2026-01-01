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
	Time         float32 // game runtime in seconds
	Started      bool    // if the game has been started
	Teams        map[string]*Team
	Server       *Server
	Map          *Map
	Running      bool
	dt           float64
	ReadyToClose chan struct{}

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

func (g *Game) StartGameLoop() {
	oldTime := time.Now()
	for g.Running {
		newTime := time.Now()
		g.dt = newTime.Sub(oldTime).Seconds()

		// todo
	}
	g.ReadyToClose <- struct{}{}
}

func (g *Game) StartCountdownStage() {
	sendToAll(g.Server.Users, "All players are ready. Let the countdown begin.")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	secondsCounter := 10

	for range ticker.C {
		addon := "..."
		if secondsCounter == 1 {
			addon = "! GO!!!"
		}
		sendToAll(g.Server.Users, fmt.Sprintf("%d%s", secondsCounter, addon))

		if secondsCounter == 1 {
			g.StartRunningStage()
			return
		}
		secondsCounter--
	}
}

func (g *Game) StartRunningStage() {
	g.Stage = GameStageRunning
	sendToAll(g.Server.Users, "The game has begun!")
}

func (g *Game) CheckIfReadyToStart() {
	allReady := true
	for _, t := range g.Teams {
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
