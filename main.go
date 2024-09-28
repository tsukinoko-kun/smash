package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"smash/internal/env"
	"smash/internal/gui"
	"smash/internal/shell"
	"smash/internal/shell/history"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			issueUrl, _ := url.Parse("https://github.com/tsukinoko-kun/smash/issues/new")
			q := issueUrl.Query()
			q.Add("title", "runtime error")
			q.Add("body", fmt.Sprintf("Runtime error caught in main: %s\n\nExplain what you did:\n\nOS: %s  \nArch: %s  \nStack:\n```\n%s\n```\n", r, runtime.GOOS, runtime.GOARCH, string(debug.Stack())))
			issueUrl.RawQuery = q.Encode()
			_, _ = fmt.Fprintf(os.Stderr, "smash: runtime error: %s\nplease report this: %s\n", r, issueUrl.String())
			<-time.After(15 * time.Second)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		for range sigint {
			<-sigint
			shell.Cancel()
		}
	}()

	defer close(sigint)

	if len(env.Config.OnStart) != 0 {
		for _, cmd := range env.Config.OnStart {
			if err := shell.Run(cmd); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
		<-time.After(400 * time.Millisecond)
	}

	for {
		if userInput, err := gui.RunPrompt(); err != nil {
			fmt.Printf("gui error: %v\n", err)
		} else {
			if userInput == "" {
				continue
			}
			_ = history.AddToHistory(userInput)
			ps2, err := shell.Ps2()
			if err != nil {
				fmt.Printf("ps2 error: %v\n", err)
			} else {
				fmt.Println(ps2 + userInput)
			}
			if err := shell.Run(userInput); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
