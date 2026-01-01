package internal

import (
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
)

var (
	currentColor = "white"

	colors = map[string]*color.Color{
		"white":  color.New(color.FgWhite),
		"red":    color.New(color.FgRed),
		"blue":   color.New(color.FgBlue),
		"green":  color.New(color.FgGreen),
		"yellow": color.New(color.FgYellow),
		"orange": color.New(color.Attribute(208)),
		"purple": color.New(color.FgMagenta),
		"grey":   color.New(color.FgHiBlack),
		"cyan":   color.New(color.FgCyan),
	}
)

func init() {
	color.NoColor = false
}

func sendToAllExcept(conn net.Conn, users map[string]*User, msg string) {
	for _, user := range users {
		if user.conn.RemoteAddr().String() != conn.RemoteAddr().String() {
			user.conn.Write([]byte(colors[currentColor].Sprint(msg)))
		}
	}
}

func sendToAll(users map[string]*User, msg string) {
	for _, user := range users {
		user.conn.Write([]byte(colors[currentColor].Sprint(msg)))
	}
}

func sendToOne(conn net.Conn, msg string) {
	conn.Write([]byte(colors[currentColor].Sprint(msg)))
}

func getFormattedTeamLayout(users map[string]*User) string {
	msg := "Players: "

	others := make([]string, 0)
	for _, us := range users {
		if us.Team == nil || us.Username == "" {
			continue
		}
		if us.Team.ID == "red" {
			others = append(others, fmt.Sprintf("%s (%s)", us.Username, colors["red"].Sprint(us.Team.Name)))
		}
	}

	for _, us := range users {
		if us.Team == nil || us.Username == "" {
			continue
		}
		if us.Team.ID == "blue" {
			others = append(others, fmt.Sprintf("%s (%s)", us.Username, colors["blue"].Sprint(us.Team.Name)))
		}
	}

	for _, us := range users {
		if us.Team != nil || us.Username == "" {
			continue
		}
		others = append(others, fmt.Sprintf("%s (no team)", us.Username))
	}

	if len(others) > 0 {
		msg = msg + strings.Join(others, ", ")
	}

	return msg + "\n"
}
