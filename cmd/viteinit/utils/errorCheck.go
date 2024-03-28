package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"unicode"
)

const VITEINIT_REFERENCE_PATH = "VITEINIT_REFERENCE_PATH"

type ProjectManifest struct {
	ProjectName      string
	ReferenceAppPath string
	Errors           []error
}

func ErrorCheck() ProjectManifest {
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

	manifest := ProjectManifest{}

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
				manifest.ReferenceAppPath = r.referenceAppPath
			case r.projectName != "":
				manifest.ProjectName = r.projectName
			}
		}
	}
	manifest.Errors = errors

	return manifest
}
