package main

import (
	"flag"
	"fmt"
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/mode"
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/utils"
	"log"
	"os"
)

func main() {
	manifest := utils.ErrorCheck()

	// this one is more of a programmer error than a user one the
	// manifest.errors should never be nil because I initialize my
	// slices not just declare them.
	if manifest.Errors == nil {
		log.Fatal("manifest.errors is nil")
	}

	if len(manifest.Errors) > 0 {
		for _, err := range manifest.Errors {
			if _, err = fmt.Fprintf(os.Stderr, "error: %v\n", err); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(1)
	}

	fmt.Println("Project pref: ", manifest.ProjectName)

	pref := ""
	flag.StringVar(&pref, "pref", "", "prefix the project pref string for the html's title and h1")
	flag.Parse()
	if pref != "" {
		manifest.ProjectName = fmt.Sprintf("%s %s", pref, manifest.ProjectName)
	}

	if err := mode.RunDefaultMode(manifest); err != nil {
		log.Fatal(err)
	}
}
