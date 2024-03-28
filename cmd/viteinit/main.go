package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const (
	NPX = "npx"
	RM  = "rm"
)

func main() {
	runCmd(
		"create vite app...",
		NPX,
		[]string{"create-vite@latest", ".", "--template", "vanilla-ts", "--force"},
	)
	runCmd(
		"removing public, src, gitignore and root html...",
		RM,
		[]string{"-r", "public", "src", ".gitignore", "index.html"},
	)
}

func runCmd(comment, command string, args []string) {
	fmt.Println(comment)
	cmd := exec.Command(command, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
