package runcommand

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Command struct {
	Comment string
	Cmd     string
	Args    []string
}

func (c *Command) RunCmd() error {
	fmt.Printf("\n\n============\n%s\n\n", c.Comment)
	cmd := exec.Command(c.Cmd, c.Args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("\ncmd.Run(): %s failed with %s\n", c.Cmd, err))
	}

	fmt.Println("\n\n...done")
	return nil
}
