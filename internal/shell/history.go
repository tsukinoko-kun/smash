package shell

import (
	"bufio"
	"fmt"
	"os"
	"smash/internal/env"
)

const maxHistory = 512

var (
	history      []string
	currentIndex int
)

func loadHistory() error {
	file, err := os.OpenFile(env.HistoryFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
		if len(history) > maxHistory {
			history = history[1:]
		}
	}

	currentIndex = len(history)
	return scanner.Err()
}

// HistoryBack goes one step back in the history and returns that entry
func HistoryBack() (string, bool) {
	if len(history) == 0 {
		if err := loadHistory(); err != nil {
			return "", false
		}
	}

	if currentIndex > 0 {
		currentIndex--
		return history[currentIndex], true
	}
	return "", false
}

// HistoryForward goes one step forward in the history and returns that entry
func HistoryForward() (string, bool) {
	if len(history) == 0 {
		if err := loadHistory(); err != nil {
			return "", false
		}
	}

	if currentIndex < len(history)-1 {
		currentIndex++
		return history[currentIndex], true
	}
	currentIndex = len(history)
	return "", true
}

// AddToHistory adds a new entry to the history
func AddToHistory(entry string) error {
	if len(history) == 0 {
		if err := loadHistory(); err != nil {
			return err
		}
	}

	if len(history) != 0 {
		latestHistoryEntry := history[len(history)-1]
		if latestHistoryEntry == entry {
			return nil
		}
	}

	history = append(history, entry)
	if len(history) > maxHistory {
		history = history[1:]
	}
	currentIndex = len(history)

	// Append to file
	file, err := os.OpenFile(env.HistoryFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintln(file, entry); err != nil {
		return fmt.Errorf("failed to write to history file: %w", err)
	}

	return nil
}

func ResetHistoryIndex() {
	currentIndex = len(history)
}
