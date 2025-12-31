package internal

import "net"

func sendToAllExcept(conn net.Conn, users map[string]*User, msg string) {
	for _, user := range users {
		if user.conn.RemoteAddr().String() != conn.RemoteAddr().String() {
			user.conn.Write([]byte(msg))
		}
	}
}

func sendToAll(conn net.Conn, users map[string]*User, msg string) {
	for _, user := range users {
		if user.conn.RemoteAddr().String() != conn.RemoteAddr().String() {
			user.conn.Write([]byte(msg))
		}
	}
}

func sendToOne(conn net.Conn, msg string) {
	conn.Write([]byte(msg))
}
