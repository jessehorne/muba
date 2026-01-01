package internal

import (
	"bufio"
	"strings"
)

func getInput(reader *bufio.Reader) string {
	msg, _ := reader.ReadString('\n')
	if len(msg) > 0 {
		msg = strings.TrimSpace(msg)
	}
	return msg
}
