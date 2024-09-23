package main

import (
	"fmt"
	"os"
	"os/signal"
	"smash/internal/gui"
	"smash/internal/shell"
)

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		for range sigint {
			<-sigint
			shell.Cancel()
		}
	}()

	defer close(sigint)

	for {
		if userInput, err := gui.RunPrompt(); err != nil {
			fmt.Printf("gui error: %v\n", err)
		} else {
			if userInput == "" {
				continue
			}
			_ = shell.AddToHistory(userInput)
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
