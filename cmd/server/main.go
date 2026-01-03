package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/jessehorne/muba/internal"
)

func init() {
	color.NoColor = false
}

func main() {
	serv, err := internal.NewServer("", "8080")
	if err != nil {
		log.Fatal(err)
	}
	err = serv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
