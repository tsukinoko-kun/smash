package shell

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"smash/internal/shell/parser"
)

func Run(s string) error {
	defer func() {
		if r := recover(); r != nil {
			issueUrl, _ := url.Parse("https://github.com/tsukinoko-kun/smash/issues/new")
			q := issueUrl.Query()
			q.Add("title", "panic on command execution")
			q.Add("body", fmt.Sprintf("command:\n```shell\n%s\n```\n\nerror:\n```\n%v\n```\n\nOS: %s  \nArch: %s  \nStack:\n```\n%s\n```\n", s, r, runtime.GOOS, runtime.GOARCH, string(debug.Stack())))
			issueUrl.RawQuery = q.Encode()
			_, _ = fmt.Fprintf(os.Stderr, "smash: %v\nplease report this: %s\n", r, issueUrl.String())
		}
	}()
	Cancel()
	commands, err := parser.Parse(s)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to parse: %s", s), err)
	}

	for _, c := range commands {
		_ = parser.RunCommand(c, os.Stdin, os.Stdout, os.Stderr)
	}

	return nil
}

func Cancel() {
	parser.Cancel()
}
