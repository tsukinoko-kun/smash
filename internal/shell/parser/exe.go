package parser

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"smash/internal/env"
	"smash/internal/system"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tsukinoko-kun/calc/calc"
)

func (e *exe) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if e.Name == "" {
		return nil
	}
	if alias, ok := env.Alias[e.Name]; ok {
		switch len(alias) {
		case 0:
			panic(fmt.Sprintf("invalid alias %s = %v", e.Name, alias))
		case 1:
			e.Name = alias[0]
		default:
			e.Name = alias[0]
			e.Args = append(alias[1:], e.Args...)
		}
	}
	if ok, err := e.internal(stdin, stdout, stderr); ok {
		if err != nil {
			os.Setenv("?", "1")
			os.Setenv("status", "1")
			_, _ = fmt.Fprintln(stderr, err.Error())
		} else {
			os.Setenv("?", "0")
			os.Setenv("status", "0")
		}
		return err
	} else {
		cmd := exec.Command(e.Name, e.Args...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		var err error
		if e.InternalParallel {
			err = cmd.Run()
		} else {
			i := trackRunning(cmd)
			err = cmd.Run()
			untrack(i)
		}
		if err != nil {
			os.Setenv("?", "1")
			os.Setenv("status", "1")
			_, _ = fmt.Fprintln(stderr, err.Error())
		} else {
			os.Setenv("?", "0")
			os.Setenv("status", "0")
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
	if len(e.Args) == 0 {
		os.Exit(0)
	} else {
		exitCode, err := strconv.Atoi(e.Args[0])
		if err != nil {
			return fmt.Errorf("exit: Illegal number: %s", e.Args[0])
		}
		os.Exit(exitCode)
	}
	return nil
}

func (e *exe) echo(stdout io.Writer) error {
	if len(e.Args) == 0 {
		_, _ = fmt.Fprintln(stdout)
	} else {
		_, _ = fmt.Fprintln(stdout, strings.Join(e.Args, " "))
	}
	return nil
}

func (e *exe) printf(stdout io.Writer) error {
	switch len(e.Args) {
	case 0:
		return errors.New("printf: usage: printf format [arg ...]")
	case 1:
		_, _ = fmt.Fprintln(stdout, e.Args[0])
	default:
		arg := make([]any, len(e.Args)-1)
		for i, a := range e.Args[1:] {
			arg[i] = a
		}
		_, _ = fmt.Fprintf(stdout, e.Args[0], arg...)
	}
	return nil
}

func (e *exe) time(stdout io.Writer, stderr io.Writer) error {
	if len(e.Args) == 0 {
		return nil
	}

	var cmd *exec.Cmd
	if len(e.Args) == 1 {
		cmd = exec.Command(e.Args[0])
	} else {
		cmd = exec.Command(e.Args[0], e.Args[1:]...)
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
	if len(e.Args) == 0 {
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
		root, err := calc.Ast(calc.Tokenize(strings.Join(e.Args, "")))
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
	switch len(e.Args) {
	case 0:
		return pushDir(env.GetUser().HomeDir)
	case 1:
		if e.Args[0] == "-" {
			return popDir()
		} else {
			return pushDir(e.Args[0])
		}
	default:
		return pushDir(strings.Join(e.Args, " "))
	}
}

func (e *exe) smashfetch(stdout io.Writer) error {
	_, _ = fmt.Fprintln(stdout, system.SmashFetch())
	return nil
}

func (s *exe) sleep() error {
	if len(s.Args) == 0 {
		return errors.New("sleep: missing operand")
	}
	var dur time.Duration
	for _, arg := range s.Args {
		d, err := time.ParseDuration(arg)
		if err != nil {
			return err
		}
		dur += d
	}
	time.Sleep(dur)
	return nil
}

const loremipsum = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.`

func (e *exe) loremipsum(stdout io.Writer) error {
	_, _ = fmt.Fprintln(stdout, loremipsum)
	return nil
}
