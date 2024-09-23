//go:build !windows

package env

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func loadVars() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell, "-l", "-c", "env")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		panic(errors.Join(fmt.Errorf("failed to get env variables from shell %s", shell), err))
	}

	env := string(out)
	env = env[:len(env)-1]

	for _, line := range strings.Split(env, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if err := os.Setenv(parts[0], parts[1]); err != nil {
			panic(errors.Join(fmt.Errorf("failed to set env variable %s", parts[0]), err))
		}
	}
}
