package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Current Directory:", currentDir)
}
