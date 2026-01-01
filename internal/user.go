package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"slices"

	"github.com/google/uuid"
)

type User struct {
	conn     net.Conn
	reader   *bufio.Reader
	ID       string
	Username string
	Player   *Player
	Server   *Server
	Team     *Team
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
	u.conn.Write([]byte("Please type a Username: "))
	for u.Username == "" {
		data := getInput(u.reader)
		if len(data) != 0 {
			u.Username = data
		}
	}
}

func (u *User) getTeamLoop() {
	team := ""
	u.conn.Write([]byte("Pick a Team (red or blue): "))
	for team == "" {
		team = getInput(u.reader)
		if team != "red" && team != "blue" {
			continue
		}

		tryTeam, ok := u.Server.Game.Teams[team]
		if !ok {
			continue
		}

		if len(tryTeam.Users) == 5 {
			sendToOne(u.conn, "That Team is full.")
			continue
		}

		u.Team = tryTeam
		tryTeam.Users = append(tryTeam.Users, u)
	}
	sendToAllExcept(u.conn, u.Server.Users, fmt.Sprintf("%s joined %s Team!\n", u.Username, team))
	sendToOne(u.conn, fmt.Sprintf("You joined %s Team!\n", team))
	log.Println("[INFO]", fmt.Sprintf("%s joined %s Team", u.Username, team))
}

func (u *User) champSelectLoop() {
	sendToOne(u.conn, "Select a champ (fighter, wizard, archer, or ninja): ")
	var champSelected bool
	var champType string
	for !champSelected {
		data := getInput(u.reader)
		if len(data) == 0 {
			continue
		}

		if slices.Contains([]string{"fighter", "wizard", "archer", "ninja"}, data) {
			champType = data
			u.Player = NewPlayer(ChampStringToType[data], u.Team)
			champSelected = true
		}
	}
	c := "white"
	switch champType {
	case "fighter":
		c = "purple"
	case "wizard":
		c = "blue"
	case "archer":
		c = "green"
	case "ninja":
		c = "yellow"
	}
	msg := fmt.Sprintf("Player '%s' has chosen to be a %s", u.GetUsernameWithTeam(), colors[c].Sprint(champType))
	sendToAllExcept(u.conn, u.Server.Users, msg)
	sendToOne(u.conn, fmt.Sprintf("You have chosen to be a %s", colors[c].Sprint(champType)))
}

func (u *User) inputLoop() {
	for {
		data := getInput(u.reader)
		if len(data) == 0 {
			continue
		}
		if data == "Team" {
			sendToOne(u.conn, "Team: ")
			sendToOne(u.conn, u.Team.Name)
		}
	}
}

func (u *User) GetUsernameWithTeam() string {
	if u.Team == nil {
		return u.Username + " (no team)"
	}
	return u.Username + " (" + colors[u.Team.Name].Sprint(u.Team.Name) + ")"
}
