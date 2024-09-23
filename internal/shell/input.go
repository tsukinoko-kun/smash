package shell

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type (
	command interface {
		Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error
		SetName(name string)
		SetArgs(args []string)
		And() command
		Or() command
		Pipe() command
		Parallel()
		IsParallel() bool
	}

	exe struct {
		name     string
		args     []string
		parallel bool
	}

	and struct {
		commands []command
		current  command
		parallel bool
	}

	or struct {
		commands []command
		current  command
		parallel bool
	}

	pipe struct {
		from     command
		to       command
		parallel bool
	}
)

func expand(s string) string {
	if len(s) == 0 {
		return s
	}
	s = os.ExpandEnv(s)
	if s[0] == '~' {
		ss := u.HomeDir + s[1:]
		if _, err := os.Stat(ss); err == nil {
			s = filepath.Clean(ss)
		}
	} else if _, err := os.Stat(s); err == nil {
		s = filepath.Clean(s)
	}
	return s
}

func expandArray(arr []string) []string {
	if arr == nil {
		return nil
	}
	for i, s := range arr {
		arr[i] = expand(s)
	}
	return arr
}

func Parse(s string) ([]command, error) {
	commands := make([]command, 0)
	var currentCommand command = &exe{}
	commands = append(commands, currentCommand)

	parser := shellwords.NewParser()
	for {
		args, err := parser.Parse(s)
		if err != nil {
			return nil, err
		}

		switch len(args) {
		case 0:
			break
		case 1:
			currentCommand.SetName(expand(args[0]))
		default:
			currentCommand.SetName(expand(args[0]))
			currentCommand.SetArgs(expandArray(args[1:]))
		}

		if parser.Position < 0 {
			break
		}

		i := parser.Position
		for ; i < len(s); i++ {
			if isSpace(s[i]) {
				break
			}
		}
		if next, err := shift(&commands, currentCommand, string([]rune(s)[parser.Position:i])); err != nil {
			return nil, err
		} else {
			currentCommand = next
		}
		if i >= len(s) {
			break
		} else {
			s = string([]rune(s)[i+1:])
		}
	}

	return commands, nil
}

func shift(commands *[]command, current command, token string) (command, error) {
	switch token {
	case "&&":
		i := len(*commands) - 1
		next := current.And()
		(*commands)[i] = next
		return next, nil
	case "||":
		i := len(*commands) - 1
		next := current.Or()
		(*commands)[i] = next
		return next, nil
	case "|":
		i := len(*commands) - 1
		next := current.Pipe()
		(*commands)[i] = next
		return next, nil
	case "&":
		current.Parallel()
		next := &exe{}
		*commands = append(*commands, next)
		return next, nil
	case ";":
		next := &exe{}
		*commands = append(*commands, next)
		return next, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", token)
	}
}

func isSpace(r byte) bool {
	switch r {
	case ' ', '\t', '\r', '\n', '\v':
		return true
	}
	return false
}

func (e *exe) SetName(name string) {
	e.name = name
}

func (a *and) SetName(name string) {
	a.current.SetName(name)
}

func (o *or) SetName(name string) {
	o.current.SetName(name)
}

func (p *pipe) SetName(name string) {
	p.to.SetName(name)
}

func (e *exe) SetArgs(args []string) {
	e.args = args
}

func (a *and) SetArgs(args []string) {
	a.current.SetArgs(args)
}

func (o *or) SetArgs(args []string) {
	o.current.SetArgs(args)
}

func (p *pipe) SetArgs(args []string) {
	p.to.SetArgs(args)
}

func (e *exe) And() command {
	next := &exe{}
	a := &and{current: next, commands: []command{e, next}}
	return a
}

func (a *and) And() command {
	a.current = &exe{}
	a.commands = append(a.commands, a.current)
	return a
}

func (o *or) And() command {
	i := len(o.commands) - 1
	next := &exe{}
	a := &and{current: next, commands: []command{o, o.commands[i], next}}
	o.current = nil
	o.commands = o.commands[:i]
	return a
}

func (p *pipe) And() command {
	next := &exe{}
	a := &and{current: next, commands: []command{p, next}}
	return a
}

func (e *exe) Or() command {
	next := &exe{}
	o := &or{current: next, commands: []command{e, next}}
	return o
}

func (a *and) Or() command {
	next := &exe{}
	o := &or{current: next, commands: []command{a, next}}
	a.current = nil
	return o
}

func (o *or) Or() command {
	o.current = &exe{}
	o.commands = append(o.commands, o.current)
	return o
}

func (p *pipe) Or() command {
	next := &exe{}
	return &or{current: next, commands: []command{p, next}}
}

func (e *exe) Pipe() command {
	return &pipe{from: e, to: &exe{}}
}

func (a *and) Pipe() command {
	a.current = &pipe{from: a.current, to: &exe{}}
	a.commands[len(a.commands)-1] = a.current
	return a
}

func (o *or) Pipe() command {
	o.current = &pipe{from: o.current, to: &exe{}}
	o.commands[len(o.commands)-1] = o.current
	return o
}

func (p *pipe) Pipe() command {
	return &pipe{from: p, to: &exe{}}
}

func (a *and) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	for _, cmd := range a.commands {
		if err := cmd.Run(stdin, stdout, stderr); err != nil {
			return err
		}
	}
	return nil
}

func (o *or) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	for _, cmd := range o.commands {
		if err := cmd.Run(stdin, stdout, stderr); err == nil {
			return nil
		}
	}
	return fmt.Errorf("smash: all commands in or failed")
}

func (p *pipe) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	stdoutReader, stdoutWriter := io.Pipe()

	var fromErr error
	var toErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer stdoutWriter.Close()
		fromErr = p.from.Run(stdin, stdoutWriter, stderr)
	}()

	go func() {
		defer wg.Done()
		toErr = p.to.Run(stdoutReader, stdout, stderr)
	}()

	wg.Wait()

	if fromErr != nil {
		return fromErr
	}
	return toErr
}

func (e *exe) Parallel() {
	e.parallel = true
}

func (a *and) Parallel() {
	a.parallel = true
}

func (o *or) Parallel() {
	o.parallel = true
}

func (p *pipe) Parallel() {
	p.parallel = true
}

func (e *exe) IsParallel() bool {
	return e.parallel
}

func (a *and) IsParallel() bool {
	return a.parallel
}

func (o *or) IsParallel() bool {
	return o.parallel
}

func (p *pipe) IsParallel() bool {
	return p.parallel
}

func RunCommand(c command, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if c.IsParallel() {
		go c.Run(stdin, stdout, stderr)
		return nil
	} else {
		return c.Run(stdin, stdout, stderr)
	}
}
