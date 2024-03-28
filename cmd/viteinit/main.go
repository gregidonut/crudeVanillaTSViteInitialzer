package main

import (
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/runcommand"
	"log"
)

const (
	NPX = "npx"
	RM  = "rm"
)

func main() {
	commands := []runcommand.Command{
		{
			Comment: "create vite app...",
			Cmd:     NPX,
			Args:    []string{"create-vite@latest", ".", "--template", "vanilla-ts", "--force"},
		},
		{
			Comment: "removing public, src, gitignore and root html...",
			Cmd:     RM,
			Args:    []string{"-r", "public", "src", ".gitignore", "index.html"},
		},
	}

	for _, command := range commands {
		if err := command.RunCmd(); err != nil {
			log.Fatal(err)
		}
	}
}
