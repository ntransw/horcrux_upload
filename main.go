package main

import (
	"log"

	"github.com/jesseduffield/horcrux/pkg/commands"
)

func main() {
	commands.Split("./test.txt", ".", 3, 3)

	paths, err := commands.GetHorcruxPathsInDir(".")
	if err != nil {
		log.Fatal(err)
	}
	commands.Bind(paths, "", true)
}
