package shell

import (
	"errors"
	"fmt"
	"github.com/tsukinoko-kun/calc/calc"
	"io"
	"net/url"
	"os"
	"os/exec"
	"smash/internal/env"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Run(s string) error {
	defer func() {
		if r := recover(); r != nil {
			issueUrl, _ := url.Parse("https://github.com/tsukinoko-kun/smash/issues/new")
			q := issueUrl.Query()
			q.Add("title", "panic on command execution")
			q.Add("body", fmt.Sprintf("command:\n```shell\n%s\n```\n\nerror:\n```\n%v\n```\n", s, r))
			issueUrl.RawQuery = q.Encode()
			_, _ = fmt.Fprintf(os.Stderr, "smash: %v\nplease report this: %s\n", r, issueUrl.String())
		}
	}()
	Cancel()
	commands, err := Parse(s)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to parse: %s", s), err)
	}

	for _, c := range commands {
		_ = RunCommand(c, os.Stdin, os.Stdout, os.Stderr)
	}

	return nil
}

func (e *exe) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if e.name == "" {
		return nil
	}
	if alias, ok := env.Alias[e.name]; ok {
		switch len(alias) {
		case 0:
			panic(fmt.Sprintf("invalid alias %s = %v", e.name, alias))
		case 1:
			e.name = alias[0]
		default:
			e.name = alias[0]
			e.args = append(alias[1:], e.args...)
		}
	}
	if ok, err := e.internal(stdin, stdout, stderr); ok {
		if err != nil {
			_, _ = fmt.Fprintln(stderr, err.Error())
		}
		return err
	} else {
		cmd := exec.Command(e.name, e.args...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		var err error
		if e.parallel {
			err = cmd.Run()
		} else {
			i := trackRunning(cmd)
			err = cmd.Run()
			untrack(i)
		}
		if err != nil {
			_, _ = fmt.Fprintln(stderr, err.Error())
		}
		return err
	}
}

var (
	running    []*exec.Cmd
	runningMut sync.Mutex
)

func trackRunning(cmd *exec.Cmd) int {
	runningMut.Lock()
	defer runningMut.Unlock()
	running = append(running, cmd)
	return len(running) - 1
}

func untrack(i int) {
	runningMut.Lock()
	defer runningMut.Unlock()
	running[i] = nil
}

func Cancel() {
	runningMut.Lock()
	defer runningMut.Unlock()
	for _, cmd := range running {
		if cmd == nil || cmd.Process == nil {
			continue
		}

		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			continue
		}

		err := cmd.Process.Signal(os.Interrupt)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to send interrupt signal: %v\n", err)
		}
	}
	running = nil
}

func (e *exe) exit() error {
	if len(e.args) == 0 {
		os.Exit(0)
	} else {
		exitCode, err := strconv.Atoi(e.args[0])
		if err != nil {
			return fmt.Errorf("exit: Illegal number: %s", e.args[0])
		}
		os.Exit(exitCode)
	}
	return nil
}

func (e *exe) echo(stdout io.Writer) error {
	if len(e.args) == 0 {
		_, _ = fmt.Fprintln(stdout)
	} else {
		_, _ = fmt.Fprintln(stdout, strings.Join(e.args, " "))
	}
	return nil
}

func (e *exe) printf(stdout io.Writer) error {
	switch len(e.args) {
	case 0:
		return errors.New("printf: usage: printf format [arg ...]")
	case 1:
		_, _ = fmt.Fprintln(stdout, e.args[0])
	default:
		arg := make([]any, len(e.args)-1)
		for i, a := range e.args[1:] {
			arg[i] = a
		}
		_, _ = fmt.Fprintf(stdout, e.args[0], arg...)
	}
	return nil
}

func (e *exe) time(stdout io.Writer, stderr io.Writer) error {
	if len(e.args) == 0 {
		return nil
	}

	var cmd *exec.Cmd
	if len(e.args) == 1 {
		cmd = exec.Command(e.args[0])
	} else {
		cmd = exec.Command(e.args[0], e.args[1:]...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	start := time.Now()
	err := cmd.Run()
	end := time.Now()
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err.Error())
	}
	_, _ = fmt.Fprintf(stdout, "\treal %v\n", end.Sub(start))
	return err
}

func (e *exe) calc(stdin io.Reader, stdout io.Writer) error {
	if len(e.args) == 0 {
		// read string from stdin
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, stdin)
		root, err := calc.Ast(calc.Tokenize(strings.TrimSpace(sb.String())))
		if err != nil {
			return err
		}
		if v, err := root.Eval(); err != nil {
			return err
		} else {
			_, _ = fmt.Fprintf(stdout, "%v\n", v)
			return nil
		}
	} else {
		// read args
		root, err := calc.Ast(calc.Tokenize(strings.Join(e.args, "")))
		if err != nil {
			return err
		}
		if v, err := root.Eval(); err != nil {
			return err
		} else {
			_, _ = fmt.Fprintf(stdout, "%v\n", v)
			return nil
		}
	}
}

var (
	cdHistory    []string
	cdHistoryMut sync.Mutex
)

func popDir() error {
	cdHistoryMut.Lock()
	defer cdHistoryMut.Unlock()
	if len(cdHistory) == 0 {
		return nil
	}
	lastCd := cdHistory[len(cdHistory)-1]
	cdHistory = cdHistory[:len(cdHistory)-1]
	if err := os.Chdir(lastCd); err != nil {
		return err
	} else {
		_ = os.Setenv("PWD", lastCd)
		return nil
	}
}

func pushDir(p string) error {
	cdHistoryMut.Lock()
	defer cdHistoryMut.Unlock()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cdHistory = append(cdHistory, wd)

	if err := os.Chdir(p); err != nil {
		return err
	} else {
		_ = os.Setenv("PWD", p)
		return nil
	}
}

func (e *exe) cd() error {
	switch len(e.args) {
	case 0:
		return pushDir(u.HomeDir)
	case 1:
		if e.args[0] == "-" {
			return popDir()
		} else {
			return pushDir(e.args[0])
		}
	default:
		return pushDir(strings.Join(e.args, " "))
	}
}
