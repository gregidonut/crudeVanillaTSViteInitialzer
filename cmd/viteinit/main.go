package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	projectName := filepath.Base(currentDir)
	fmt.Println("project name:", projectName)
}
