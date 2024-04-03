package mode

import (
	"fmt"
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/runcommand"
	"github.com/gregidonut/crudeVanillaTSViteInitialzer/cmd/viteinit/utils"
	"log"
	"path/filepath"
	"sync"
)

// RunDefaultMode is the default logic that runs after
// initial crude error check if no mode is selected.
func RunDefaultMode(manifest utils.ProjectManifest) error {

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
				filepath.Join(manifest.ReferenceAppPath, "src"),
				filepath.Join(manifest.ReferenceAppPath, ".eslintrc.cjs"),
				filepath.Join(manifest.ReferenceAppPath, ".eslintignore"),
				filepath.Join(manifest.ReferenceAppPath, ".prettierrc.cjs"),
				filepath.Join(manifest.ReferenceAppPath, ".prettierignore"),
				filepath.Join(manifest.ReferenceAppPath, "vite.config.ts"),
				filepath.Join(manifest.ReferenceAppPath, ".gitignore"),
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
					fmt.Sprintf("s|<title>[^<]*<\\/title>|<title>%s<\\/title>|g", manifest.ProjectName),
					"src/index.html",
				},
			},
			{
				Comment: "replacing h1 tag inner text in index.html...",
				Cmd:     "sed",
				Args: []string{
					"-i",
					fmt.Sprintf("s|<h1>[^<]*<\\/h1>|<h1>%s<\\/h1>|g", manifest.ProjectName),
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
  I can mimic a regular static file server
- got to hello world equivalent index.html
- also added some styles with colors  that's sort of inspired by the
  dracula theme`,
			},
		},
	}

	for _, command := range gitCommands {
		if err := command.RunCmd(); err != nil {
			log.Fatal(err)
		}

	}

	return nil
}
