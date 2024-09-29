//go:build !windows

package parser

import "io"

var InternalToolNames = [...]string{
	"exit",
	"echo",
	"printf",
	"time",
	"calc",
	"cd",
}

func (e *exe) internal(stdin io.Reader, stdout io.Writer, stderr io.Writer) (bool, error) {
	switch e.Name {
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
	case "smashfetch":
		return true, e.smashfetch(stdout)
	case "sleep":
		return true, e.sleep()
	default:
		return false, nil
	}
}
