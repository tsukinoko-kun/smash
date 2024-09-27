package parser

import (
	"fmt"
	"io"
	"sync"
)

type (
	Command interface {
		Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error
		PushName(name string)
		And() Command
		Or() Command
		Pipe() Command
		Parallel()
		IsParallel() bool
		GetLast() *exe
	}

	exe struct {
		Name             string
		Args             []string
		InternalParallel bool
	}

	and struct {
		Commands         []Command
		Current          Command
		InternalParallel bool
	}

	or struct {
		Commands         []Command
		Current          Command
		InternalParallel bool
	}

	pipe struct {
		From             Command
		To               Command
		InternalParallel bool
	}
)

func Parse(s string) ([]Command, error) {
	var currentCommand Command = &exe{}
	commands := []Command{currentCommand}

	tokens := Tokenize(s)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.Type {
		case TokenTypeName:
			currentCommand.PushName(t.Content)
		case TokenTypeOperator:
			if next, err := shift(&commands, currentCommand, t.Content); err != nil {
				return nil, err
			} else {
				currentCommand = next
			}
		}
	}

	return commands, nil
}

func shift(commands *[]Command, current Command, token string) (Command, error) {
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

func (e *exe) PushName(name string) {
	if len(e.Name) == 0 {
		e.Name = name
	} else {
		e.Args = append(e.Args, name)
	}
}

func (a *and) PushName(name string) {
	a.Current.PushName(name)
}

func (o *or) PushName(name string) {
	o.Current.PushName(name)
}

func (p *pipe) PushName(name string) {
	p.To.PushName(name)
}

func (e *exe) And() Command {
	next := &exe{}
	a := &and{Current: next, Commands: []Command{e, next}}
	return a
}

func (a *and) And() Command {
	a.Current = &exe{}
	a.Commands = append(a.Commands, a.Current)
	return a
}

func (o *or) And() Command {
	i := len(o.Commands) - 1
	next := &exe{}
	a := &and{Current: next, Commands: []Command{o, o.Commands[i], next}}
	o.Current = nil
	o.Commands = o.Commands[:i]
	return a
}

func (p *pipe) And() Command {
	next := &exe{}
	a := &and{Current: next, Commands: []Command{p, next}}
	return a
}

func (e *exe) Or() Command {
	next := &exe{}
	o := &or{Current: next, Commands: []Command{e, next}}
	return o
}

func (a *and) Or() Command {
	next := &exe{}
	o := &or{Current: next, Commands: []Command{a, next}}
	a.Current = nil
	return o
}

func (o *or) Or() Command {
	o.Current = &exe{}
	o.Commands = append(o.Commands, o.Current)
	return o
}

func (p *pipe) Or() Command {
	next := &exe{}
	return &or{Current: next, Commands: []Command{p, next}}
}

func (e *exe) Pipe() Command {
	return &pipe{From: e, To: &exe{}}
}

func (a *and) Pipe() Command {
	a.Current = &pipe{From: a.Current, To: &exe{}}
	a.Commands[len(a.Commands)-1] = a.Current
	return a
}

func (o *or) Pipe() Command {
	o.Current = &pipe{From: o.Current, To: &exe{}}
	o.Commands[len(o.Commands)-1] = o.Current
	return o
}

func (p *pipe) Pipe() Command {
	return &pipe{From: p, To: &exe{}}
}

func (a *and) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	for _, cmd := range a.Commands {
		if err := cmd.Run(stdin, stdout, stderr); err != nil {
			return err
		}
	}
	return nil
}

func (o *or) Run(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	for _, cmd := range o.Commands {
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
		fromErr = p.From.Run(stdin, stdoutWriter, stderr)
	}()

	go func() {
		defer wg.Done()
		toErr = p.To.Run(stdoutReader, stdout, stderr)
	}()

	wg.Wait()

	if fromErr != nil {
		return fromErr
	}
	return toErr
}

func (e *exe) Parallel() {
	e.InternalParallel = true
}

func (a *and) Parallel() {
	a.InternalParallel = true
}

func (o *or) Parallel() {
	o.InternalParallel = true
}

func (p *pipe) Parallel() {
	p.InternalParallel = true
}

func (e *exe) IsParallel() bool {
	return e.InternalParallel
}

func (a *and) IsParallel() bool {
	return a.InternalParallel
}

func (o *or) IsParallel() bool {
	return o.InternalParallel
}

func (p *pipe) IsParallel() bool {
	return p.InternalParallel
}

func (e *exe) GetLast() *exe {
	return e
}

func (a *and) GetLast() *exe {
	return a.Commands[len(a.Commands)-1].GetLast()
}

func (o *or) GetLast() *exe {
	return o.Commands[len(o.Commands)-1].GetLast()
}

func (p *pipe) GetLast() *exe {
	return p.To.GetLast()
}

func RunCommand(c Command, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if c.IsParallel() {
		go c.Run(stdin, stdout, stderr)
		return nil
	} else {
		return c.Run(stdin, stdout, stderr)
	}
}
