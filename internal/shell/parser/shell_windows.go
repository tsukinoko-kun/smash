//go:build windows

package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var InternalToolNames = [...]string{
	"exit",
	"echo",
	"printf",
	"time",
	"calc",
	"cd",
	"zu",
	"env",
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
	case "zu":
		return true, e.zu()
	case "smashfetch":
		return true, e.smashfetch(stdout)
	case "sleep":
		return true, e.sleep()
	case "env":
		return true, e.env(stdout)
	default:
		return false, nil
	}
}

func (e *exe) env(stdout io.Writer) error {
	sb := strings.Builder{}
	for _, e := range os.Environ() {
		sb.WriteString(e)
		sb.WriteRune('\n')
	}
	_, _ = fmt.Fprint(stdout, sb.String())
	return nil
}
