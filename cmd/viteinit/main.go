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
		log.Fatal("manifest.errors is nil")
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

	initialCommands := []runcommand.Command{
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
	for _, command := range initialCommands {
		if err := command.RunCmd(); err != nil {
			log.Fatal(err)
		}
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		sedCommandsForIndexHtml := []runcommand.Command{
			{
				Comment: "replacing title tag inner text in index.html...",
				Cmd:     "sed",
				Args: []string{
					"-i",
					fmt.Sprintf("s|<title>[^<]*<\\/title>|<title>%s<\\/title>|g", manifest.projectName),
					"src/index.html",
				},
			},
			{
				Comment: "replacing h1 tag inner text in index.html...",
				Cmd:     "sed",
				Args: []string{
					"-i",
					fmt.Sprintf("s|<h1>[^<]*<\\/h1>|<h1>%s<\\/h1>|g", manifest.projectName),
					"src/index.html",
				},
			},
		}

		for _, command := range sedCommandsForIndexHtml {
			if err := command.RunCmd(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		sedCommandsForPackageJson := runcommand.Command{
			Comment: "replacing build script, and appending lint script in package.json...",
			Cmd:     "sed",
			Args: []string{
				"-i",
				`s|"build": "tsc \&\& vite build"|"build": "eslint \&\& tsc \&\& vite build",\n    "lint": "eslint . --ext .ts"|`,
				"package.json",
			},
		}
		if err := sedCommandsForPackageJson.RunCmd(); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	gitCommands := []runcommand.Command{
		{
			Comment: "git init",
			Cmd:     "git",
			Args: []string{
				"init",
			},
		},
		{
			Comment: "git add -A",
			Cmd:     "git",
			Args: []string{
				"add",
				"-A",
			},
		},
		{
			Comment: "git commit",
			Cmd:     "git",
			Args: []string{
				"commit",
				"-m",
				"Initialize project",
				"-m",
				`- vanilla-ts vite, with prettier and eslint
- added node.gitignore github template to vite generated gitignore
- implemented a crude directory walk in the vite.config.ts file so that
  I can mimic regular static file server
- got to hello world equivalent index.html`,
			},
		},
	}

	for _, command := range gitCommands {
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
						"directory '%s' from env var does not exist", referenceAppPath,
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
