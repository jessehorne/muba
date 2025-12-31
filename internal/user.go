package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	conn     net.Conn
	reader   *bufio.Reader
	id       string
	username string
	player   *Player
	server   *Server
	team     *Team
}

func NewUser(conn net.Conn, s *Server) *User {
	return &User{
		id:     uuid.New().String(),
		conn:   conn,
		reader: bufio.NewReader(conn),
		server: s,
	}
}

func (u *User) Handle() error {
	u.conn.Write([]byte(u.server.welcomeMessage))

	u.getUsernameLoop()
	u.announceUserJoin()
	u.getTeamLoop()
	u.inputLoop()

	return nil
}

func (u *User) announceUserJoin() {
	for _, user := range u.server.Users {
		if user.conn.RemoteAddr().String() != u.conn.RemoteAddr().String() {
			user.conn.Write([]byte(fmt.Sprintf("A user by the name of '%s' has joined your game.", u.username)))
		} else {
			var users []string
			for _, uu := range u.server.Users {
				if user.conn.RemoteAddr().String() != u.conn.RemoteAddr().String() {
					addon := "(no team)"
					if uu.team != nil {
						addon = fmt.Sprintf("(%s)", uu.team.Name)
					}
					users = append(users, fmt.Sprintf("%s %s", uu.username, addon))
				}
			}
			msg := "You have joined the game.\n"
			if len(users) > 0 {
				msg = msg + fmt.Sprintf("Your fellow players are...\n%s\n", strings.Join(users, "\n"))
			}
			u.conn.Write([]byte(msg))
		}
	}
}

func (u *User) getUsernameLoop() {
	u.conn.Write([]byte("Please type a username: "))
	for u.username == "" {
		data := getInput(u.reader)
		if len(data) != 0 {
			u.username = data
		}
	}
}

func (u *User) getTeamLoop() {
	team := ""
	u.conn.Write([]byte("Pick a team (red or blue): "))
	for team == "" {
		team = getInput(u.reader)
		if team == "red" || team == "blue" {
			u.team = u.server.Teams[team]
		}
	}
	sendToAllExcept(u.conn, u.server.Users, fmt.Sprintf("%s joined %s team!\n", u.username, team))
	sendToOne(u.conn, fmt.Sprintf("You joined %s team!\n", team))
	log.Println("[INFO]", fmt.Sprintf("%s joined %s team", u.username, team))
}

func (u *User) inputLoop() {
	for {
		getInput(u.reader)
	}
}
