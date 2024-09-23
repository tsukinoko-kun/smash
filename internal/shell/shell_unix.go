//go:build !windows

package shell

import "io"

var internalToolNames = [...]string{
	"exit",
	"echo",
	"printf",
	"time",
	"calc",
	"cd",
}

func (e *exe) internal(stdin io.Reader, stdout io.Writer, stderr io.Writer) (bool, error) {
	switch e.name {
	case "exit":
		return true, e.exit()
	case "echo":
		return true, e.echo(stdout)
	case "printf":
		return true, e.printf(stdout)
	case "time":
		return true, e.time(stdout, stderr)
	case "calc":
		return true, e.calc(stdin, stdout)
	case "cd":
		return true, e.cd()
	default:
		return false, nil
	}
}
