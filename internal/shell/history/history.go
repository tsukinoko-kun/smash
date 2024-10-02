package history

import (
	"bufio"
	"fmt"
	"os"
	"smash/internal/env"
)

const (
	historyHwm = 128
	historyLwm = 64
)

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
	}

	if len(history) > historyHwm {
		history = history[historyLwm:]
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
	// ensure the latest history is loaded
	if len(history) == 0 {
		if err := loadHistory(); err != nil {
			return err
		}
	}

	// check if the latest history entry is the same as the new entry
	if len(history) != 0 {
		latestHistoryEntry := history[len(history)-1]
		if latestHistoryEntry == entry {
			return nil
		}
	}

	// append new entry to history in memory
	history = append(history, entry)

	// ensure history length does not exceed the high water mark
	if len(history) > historyHwm {
		history = history[historyLwm:]

		// write whole history to history file
		file, err := os.Create(env.HistoryFile)
		if err != nil {
			return fmt.Errorf("failed to open history file: %w", err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		defer writer.Flush()
		for _, entry := range history {
			if _, err := writer.WriteString(entry + "\n"); err != nil {
				return fmt.Errorf("failed to write to history file: %w", err)
			}
		}
	} else {
		// append new entry to history file
		file, err := os.OpenFile(env.HistoryFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open history file: %w", err)
		}
		defer file.Close()
		if _, err := file.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}
	}
	currentIndex = len(history)

	return nil
}

func ResetHistoryIndex() {
	currentIndex = len(history)
}
