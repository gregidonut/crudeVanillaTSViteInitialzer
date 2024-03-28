package main

import (
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/runcommand"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}

	projectName := filepath.Base(currentDir)
	// vite prompts for name if project name is not lowercase
	// and I want to avoid dealing with prompts
	for _, char := range projectName {
		if unicode.IsUpper(char) {
			log.Fatalf("Error: project name: '%s' cannot have upper case letters\n", projectName)
		}
	}

	commands := []runcommand.Command{
		{
			Comment: "create vite app...",
			Cmd:     "npm",
			Args:    []string{"create", "vite@latest", ".", "--", "--template", "vanilla-ts"},
		},
		{
			Comment: "installing initial vite app...",
			Cmd:     "npm",
			Args: []string{
				"install",
				"eslint",
				"prettier",
				"@typescript-eslint/eslint-plugin",
				"@typescript-eslint/parser",
				"eslint-config-prettier",
				"eslint-plugin-import",
				"@types/node",
				"--save-dev",
			},
		},
		{
			Comment: "removing public, src, gitignore and root html...",
			Cmd:     "rm",
			Args:    []string{"-r", "public", "src", ".gitignore", "index.html"},
		},
	}
	for _, command := range commands {
		if err := command.RunCmd(); err != nil {
			log.Fatal(err)
		}
	}
}
