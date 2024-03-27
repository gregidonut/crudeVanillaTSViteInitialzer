package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please define project name as first argument")
		return
	}

	fmt.Println("First argument:", os.Args[1])
}
