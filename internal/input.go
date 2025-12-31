package internal

import (
	"bufio"
)

func getInput(reader *bufio.Reader) string {
	msg, _ := reader.ReadString('\n')
	if len(msg) > 0 {
		msg = msg[:len(msg)-1]
	}
	return msg
}
