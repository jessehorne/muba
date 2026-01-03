package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"

	"github.com/google/uuid"
)

var (
	directions = []string{"north", "south", "east", "west", "n", "s", "e", "w", "ne", "nw", "se", "sw"}

	DirStringToVector = map[string]Vector2{
		"north": {0, -1}, "n": {0, -1},
		"south": {0, 1}, "s": {0, 1},
		"east": {1, 0}, "e": {1, 0},
		"west": {-1, 0}, "w": {-1, 0},
		"ne": {1, -1},
		"nw": {-1, -1},
		"se": {1, 1},
		"sw": {-1, 1},
	}

	DirVectorToString = map[Vector2]string{
		{0, -1}:  "north",
		{0, 1}:   "south",
		{1, 0}:   "east",
		{-1, 0}:  "west",
		{1, -1}:  "ne",
		{-1, -1}: "nw",
		{1, 1}:   "se",
		{-1, 1}:  "sw",
	}

	DirShortToLong = map[string]string{
		"north": "north",
		"south": "south",
		"east":  "east",
		"west":  "west",
		"ne":    "north east",
		"nw":    "north west",
		"se":    "south east",
		"sw":    "south west",
	}
)

type User struct {
	conn      net.Conn
	reader    *bufio.Reader
	ID        string
	Username  string
	ChampType int
	GameStats *GameStats
	Coords    Vector2
	Server    *Server
	Team      *Team
	Ready     bool

	Health        int
	CurrentHealth int
	HealPerSecond int
}

func NewUser(conn net.Conn, s *Server) *User {
	return &User{
		ID:     uuid.New().String(),
		conn:   conn,
		reader: bufio.NewReader(conn),
		Server: s,
	}
}

func (u *User) Handle() error {
	u.conn.Write([]byte(u.Server.welcomeMessage))

	u.getUsernameLoop()
	u.announceUserJoin()
	u.getTeamLoop()
	u.champSelectLoop()
	u.inputLoop()

	return nil
}

func (u *User) announceUserJoin() {
	// let everyone know you joined
	sendToAllExcept(u.conn, u.Server.Users, fmt.Sprintf("A user by the name of '%s' has joined your game.", u.Username))

	// show joining Player list of other players and their teams
	sendToOne(u.conn, getFormattedTeamLayout(u.Server.Users))
}

func (u *User) getUsernameLoop() {
	sendToOne(u.conn, "What should we call you?\n")
	for u.Username == "" {
		data := getInput(u.reader)
		if len(data) != 0 {
			u.Username = data
		}
	}
}

func (u *User) getTeamLoop() {
	teams := make([]string, 0)
	teamsFormatted := make([]string, 0)
	for teamID, t := range u.Server.Game.Teams {
		teams = append(teams, t.ID)
		teamsFormatted = append(teamsFormatted, fmt.Sprintf("%s (%d)", teamID, len(t.Users)))
	}
	team := ""
	sendToOne(u.conn, fmt.Sprintf("Pick a team: %s\n", strings.Join(teamsFormatted, ", ")))
	for team == "" {
		team = getInput(u.reader)
		if !slices.Contains(teams, team) {
			team = ""
			sendToOne(u.conn, fmt.Sprintf("There is no team named %s.\n", team))
			continue
		}

		tryTeam, ok := u.Server.Game.Teams[team]
		if !ok {
			team = ""
			sendToOne(u.conn, fmt.Sprintf("There is no team named %s.\n", team))
			continue
		}

		if len(tryTeam.Users) == tryTeam.Size {
			sendToOne(u.conn, "That team is full.\n")
			team = ""
			continue
		}

		u.Team = tryTeam
		u.Coords = u.Server.Game.Map.RoomIDToCoords(tryTeam.Base)
		tryTeam.Users = append(tryTeam.Users, u)
	}
	sendToAllExcept(u.conn, u.Server.Users, fmt.Sprintf("%s joined %s team!\n", u.Username, colors[u.Team.ID].Sprint(u.Team.Name)))
	sendToOne(u.conn, fmt.Sprintf("You joined %s team!\n", colors[u.Team.ID].Sprint(u.Team.Name)))
	log.Println("[INFO]", fmt.Sprintf("%s joined %s team", u.Username, team))
}

func (u *User) champSelectLoop() {
	sendToOne(u.conn, fmt.Sprintf("Select a champ (%s, %s, %s, or %s): ",
		ChampTypeColored["fighter"], ChampTypeColored["wizard"],
		ChampTypeColored["archer"], ChampTypeColored["ninja"]))
	var champSelected bool
	var champType string
	for !champSelected {
		data := getInput(u.reader)
		if len(data) == 0 {
			continue
		}

		if slices.Contains([]string{"fighter", "wizard", "archer", "ninja"}, data) {
			champType = data
			u.ChampType = ChampStringToType[data]
			u.GameStats = NewGameStats(u.ChampType)
			u.Health = u.GameStats.Health
			u.CurrentHealth = u.GameStats.Health
			u.HealPerSecond = u.GameStats.HealPerSecond
			champSelected = true
		}
	}
	msg := fmt.Sprintf("Player '%s' has chosen to be a %s", u.GetUsernameWithTeam(), ChampTypeColored[champType])
	sendToAllExcept(u.conn, u.Server.Users, msg)
	sendToOne(u.conn, fmt.Sprintf("You have chosen to be a %s", ChampTypeColored[champType]))
}

func (u *User) inputLoop() {
	for {
		data := getInput(u.reader)
		if len(data) == 0 {
			continue
		}
		if data == "team" {
			sendToOne(u.conn, "Team: ")
			sendToOne(u.conn, u.Team.Name)
		} else if data == "ready" {
			u.SetReady()
		} else if data == "look" || data == "l" {
			u.Look()
		} else if data == "health" || data == "h" {
			sendToOne(u.conn, u.FormattedHealth()+"\n")
		} else if data == "me" {
			sendToOne(u.conn, u.FormattedQuickStat())
		} else if slices.Contains(directions, data) {
			u.Move(data)
		}
	}
}

func (u *User) GetUsernameWithTeam() string {
	if u.Team == nil {
		return u.Username + " (no team)"
	}
	return u.Username + " (" + colors[u.Team.ID].Sprint(u.Team.Name) + ")"
}

func (u *User) SetReady() {
	u.Ready = true
	sendToAll(u.Server.Users, fmt.Sprintf("Player '%s' is ready!\n", u.GetUsernameWithTeam()))
	u.Server.Game.CheckIfReadyToStart()
	log.Printf("[INFO] player %s is ready and their coords are %v\n", u.Username, u.Coords)
}

func (u *User) IsInBase() bool {
	return u.Server.Game.Map.CoordsToRoomID(u.Coords) == u.Team.ID
}

// Move attempts to move the user in a specific direction.
func (u *User) Move(dir string) {
	vel, ok := DirStringToVector[dir]
	if !ok {
		sendToOne(u.conn, "Invalid directional command. Please contact an admin.")
		return
	}

	nextRoomCoords := u.Coords.Add(vel)
	nextRoom := u.Server.Game.Map.CoordsToRoomID(nextRoomCoords)
	if nextRoom == "" {
		sendToOne(u.conn, "You can't go that way.")
		return
	}

	u.Coords = nextRoomCoords
	sendToAllExcept(u.conn, u.Server.Game.GetUsersInRoom(u.Coords), fmt.Sprintf("Player %s went %s.", u.GetUsernameWithTeam(), DirShortToLong[dir]))
	u.Look()
}

// Look tells the player what they see in the current room.
func (u *User) Look() {
	msg := fmt.Sprintf("Coordinates: %d,%d\n", u.Coords.X, u.Coords.Y)

	r := u.Server.Game.Map.CoordsToRoom(u.Coords)
	if r != nil {
		if len(r.Towers) > 0 {
			towerAddon := colors[r.Towers[0].Team.ID].Sprint(r.Towers[0].Team.Name)
			msg = msg + fmt.Sprintf("There sits a large %s tower.\n", towerAddon)
		}
	}

	sendToOne(u.conn, msg)
}

func (u *User) FormattedQuickStat() string {
	msg := "%s\t[%s]\t%s\t(%d,%d)\n"
	return fmt.Sprintf(msg,
		u.GetUsernameWithTeam(), u.FormattedHealth(),
		ChampTypeColored[ChampTypeToString[u.ChampType]], u.Coords.X, u.Coords.Y)
}

func (u *User) EnemyTeam() *Team {
	for _, t := range u.Server.Game.Teams {
		if t != u.Team {
			return t
		}
	}
	return nil
}

func (u *User) FormattedHealth() string {
	c := colors["green"]
	perc := float64(u.CurrentHealth / u.Health)
	if perc < 0.2 {
		c = colors["red"]
	} else if perc < 0.7 {
		c = colors["yellow"]
	}
	msg := c.Sprint(u.CurrentHealth) + "/" + colors["green"].Sprint(u.Health)
	return msg
}
