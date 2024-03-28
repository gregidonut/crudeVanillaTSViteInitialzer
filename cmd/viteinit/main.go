package main

import (
	"fmt"
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/runcommand"
	"log"
	"os"
	"path/filepath"
	"sync"
	"unicode"
)

const VITEINIT_REFERENCE_PATH = "VITEINIT_REFERENCE_PATH"

func main() {
	fmt.Println("---------------------")
	fmt.Println()
	referenceAppPath, errors := errorCheck()
	if len(errors) > 0 || errors == nil {
		for _, err := range errors {
			if _, err = fmt.Fprintf(os.Stderr, "Error: %v\n", err); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(1)
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
		{
			Comment: "copying src, .eslintrc, .eslintignore, .prettierrc, .prettierignore, vite.config and .gitignore",
			Cmd:     "cp",
			Args: []string{
				"-r",
				filepath.Join(referenceAppPath, "src"),
				filepath.Join(referenceAppPath, ".eslintrc.cjs"),
				filepath.Join(referenceAppPath, ".eslintignore"),
				filepath.Join(referenceAppPath, ".prettierrc.cjs"),
				filepath.Join(referenceAppPath, ".prettierignore"),
				filepath.Join(referenceAppPath, "vite.config.ts"),
				filepath.Join(referenceAppPath, ".gitignore"),
				".",
			},
		},
	}
	for _, command := range commands {
		if err := command.RunCmd(); err != nil {
			log.Fatal(err)
		}
	}
}

func errorCheck() (string, []error) {
	wg := sync.WaitGroup{}

	type result struct {
		referenceAppPath string
		err              error
	}

	resultChan := make(chan result)
	errors := []error{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		currentDir, err := os.Getwd()
		if err != nil {
			resultChan <- result{
				err: fmt.Errorf("error: %v", err),
			}
			return
		}

		projectName := filepath.Base(currentDir)
		// vite prompts for name if project name is not lowercase
		// and I want to avoid dealing with prompts
		for _, char := range projectName {
			if unicode.IsUpper(char) {
				resultChan <- result{
					err: fmt.Errorf(
						"error: project name: '%s' cannot have upper case letters", projectName,
					),
				}
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		referenceAppPath := os.Getenv(VITEINIT_REFERENCE_PATH)
		if referenceAppPath == "" {
			resultChan <- result{
				referenceAppPath: "",
				err: fmt.Errorf(
					"could not get reference path at '%s' env var", VITEINIT_REFERENCE_PATH,
				),
			}
			return
		}
		_, err := os.Stat(referenceAppPath)
		if err != nil {
			if os.IsNotExist(err) {
				resultChan <- result{
					referenceAppPath: "",
					err: fmt.Errorf(
						"directory from env var '%s' does not exist", referenceAppPath,
					),
				}
				return
			}
			resultChan <- result{
				referenceAppPath: "",
				err:              fmt.Errorf("error: %v", err),
			}
			return
		}

		resultChan <- result{
			referenceAppPath: referenceAppPath,
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	referenceAppPath := ""

outer:
	for {
		select {
		case r, ok := <-resultChan:
			if !ok {
				break outer
			}
			if r.err != nil {
				errors = append(errors, r.err)
				continue
			}
			referenceAppPath = r.referenceAppPath
		}
	}

	return referenceAppPath, errors
}
