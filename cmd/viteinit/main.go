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
	manifest := errorCheck()

	// this one is more of a programmer error than a user one the
	// manifest.errors should never be nil because I initialize my
	// slices not just declare them.
	if manifest.errors == nil {
		log.Fatal("manifest errors is nil")
	}

	if len(manifest.errors) > 0 {
		for _, err := range manifest.errors {
			if _, err = fmt.Fprintf(os.Stderr, "error: %v\n", err); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(1)
	}

	fmt.Println("Project name: ", manifest.projectName)

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
				filepath.Join(manifest.referenceAppPath, "src"),
				filepath.Join(manifest.referenceAppPath, ".eslintrc.cjs"),
				filepath.Join(manifest.referenceAppPath, ".eslintignore"),
				filepath.Join(manifest.referenceAppPath, ".prettierrc.cjs"),
				filepath.Join(manifest.referenceAppPath, ".prettierignore"),
				filepath.Join(manifest.referenceAppPath, "vite.config.ts"),
				filepath.Join(manifest.referenceAppPath, ".gitignore"),
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

type projectManifest struct {
	projectName      string
	referenceAppPath string
	errors           []error
}

func errorCheck() projectManifest {
	wg := sync.WaitGroup{}

	type errorCheckResult struct {
		projectName      string
		referenceAppPath string
		err              error
	}

	resultChan := make(chan errorCheckResult)
	errors := []error{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		currentDir, err := os.Getwd()
		if err != nil {
			resultChan <- errorCheckResult{
				err: fmt.Errorf("error: %v", err),
			}
			return
		}

		projectName := filepath.Base(currentDir)
		// vite prompts for name if project name is not lowercase
		// and I want to avoid dealing with prompts
		for _, char := range projectName {
			if unicode.IsUpper(char) {
				resultChan <- errorCheckResult{
					err: fmt.Errorf(
						"error: project name: '%s' cannot have upper case letters", projectName,
					),
				}
				return
			}
		}

		resultChan <- errorCheckResult{
			projectName: projectName,
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		referenceAppPath := os.Getenv(VITEINIT_REFERENCE_PATH)
		if referenceAppPath == "" {
			resultChan <- errorCheckResult{
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
				resultChan <- errorCheckResult{
					referenceAppPath: "",
					err: fmt.Errorf(
						"directory from env var '%s' does not exist", referenceAppPath,
					),
				}
				return
			}
			resultChan <- errorCheckResult{
				referenceAppPath: "",
				err:              fmt.Errorf("error: %v", err),
			}
			return
		}

		resultChan <- errorCheckResult{
			referenceAppPath: referenceAppPath,
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	manifest := projectManifest{}

outer:
	for {
		select {
		case r, ok := <-resultChan:
			switch {
			case !ok:
				break outer
			case r.err != nil:
				errors = append(errors, r.err)
			case r.referenceAppPath != "":
				manifest.referenceAppPath = r.referenceAppPath
			case r.projectName != "":
				manifest.projectName = r.projectName
			}
		}
	}
	manifest.errors = errors

	return manifest
}
