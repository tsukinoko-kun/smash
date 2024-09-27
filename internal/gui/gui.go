package gui

import (
	"fmt"
	"smash/internal/shell/extensions"
	"smash/internal/shell/history"
	"strings"

	"smash/internal/env"
	"smash/internal/shell"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	SelectedCompletionBg = lipgloss.Color(env.Config.Color.CompletionSelectedBg)
	CompletionTextColor  = lipgloss.Color(env.Config.Color.CompletionText)
)

type model struct {
	textInput    textinput.Model
	err          error
	lines        int
	prompt       string
	completions  []extensions.Completion
	selectedComp int
	showComp     bool
}

func initialModel() (*model, error) {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Prompt = "" // Remove the default '>' prompt

	prompt, err := shell.Ps1()
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	return &model{
		textInput: ti,
		err:       nil,
		lines:     0,
		prompt:    prompt,
	}, nil
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	trimedInputValue := strings.TrimSpace(m.textInput.Value())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			history.ResetHistoryIndex()
			if !m.showComp {
				if len(trimedInputValue) == 0 {
					return m, nil
				}
				m.showComp = true
				m.completions = shell.GetCompletions(m.textInput.Value())
				m.selectedComp = 0
			} else {
				if len(m.completions) == 0 {
					m.showComp = false
					m.selectedComp = 0
					return m, nil
				} else {
					m.selectedComp = (m.selectedComp + 1) % len(m.completions)
				}
			}
			return m, nil
		case tea.KeyShiftTab:
			if m.showComp {
				m.selectedComp = (m.selectedComp - 1 + len(m.completions)) % len(m.completions)
			}
			return m, nil
		case tea.KeyUp, tea.KeyDown:
			if m.showComp {
				if msg.Type == tea.KeyUp {
					m.selectedComp = (m.selectedComp - 1 + len(m.completions)) % len(m.completions)
				} else {
					m.selectedComp = (m.selectedComp + 1) % len(m.completions)
				}
				return m, nil
			} else {
				var hisContent string
				var hisOk bool
				if msg.Type == tea.KeyUp {
					hisContent, hisOk = history.HistoryBack()
				} else {
					hisContent, hisOk = history.HistoryForward()
				}
				if hisOk {
					m.textInput.SetValue(hisContent)
					m.textInput.CursorEnd()
				}
			}
		case tea.KeyEsc:
			if m.showComp {
				m.showComp = false
				return m, nil
			}
		case tea.KeyEnter:
			history.ResetHistoryIndex()
			if m.showComp {
				c := m.completions[m.selectedComp]
				m.textInput.SetValue(c.Value)
				m.textInput.CursorEnd()
				m.showComp = false
				return m, nil
			}
			if m.textInput.Value() != "" {
				return m, tea.Quit
			}
		case tea.KeyCtrlC:
			history.ResetHistoryIndex()
			m.textInput.SetValue("")
			return m, tea.Quit
		}
	}

	pre := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)
	trimedInputValue = strings.TrimSpace(m.textInput.Value())

	changed := pre != m.textInput.Value()

	if changed {
		if m.showComp {
			if len(trimedInputValue) == 0 {
				m.showComp = false
				m.completions = nil
			} else {
				newCompletions := shell.GetCompletions(m.textInput.Value())
				if !completionsEqual(m.completions, newCompletions) {
					m.completions = newCompletions
					m.selectedComp = 0
				}
			}
		}
		history.ResetHistoryIndex()
	}

	return m, cmd
}

// Helper function to compare two slices of completions
func completionsEqual(a, b []extensions.Completion) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	view := lipgloss.NewStyle().Render(m.prompt + m.textInput.View())

	if m.showComp {
		view += "\n"
		for i, comp := range m.completions {
			style := lipgloss.NewStyle().Foreground(CompletionTextColor)
			if i == m.selectedComp {
				style = style.Background(SelectedCompletionBg)
			}
			view += style.Render(comp.String()) + "\n"
		}
	}

	m.lines = strings.Count(view, "\n")

	return view
}

func RunPrompt() (string, error) {
	SelectedCompletionBg = lipgloss.Color(env.Config.Color.CompletionSelectedBg)
	CompletionTextColor = lipgloss.Color(env.Config.Color.CompletionText)
	history.ResetHistoryIndex()

	initialModel, err := initialModel()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(initialModel)

	m, err := p.Run()
	if err != nil {
		return "", err
	}

	// Clear only the lines we used
	fmt.Printf("\033[%dA", m.(*model).lines) // Move cursor up
	fmt.Printf("\033[J")                     // Clear from cursor to end of screen

	return strings.TrimSpace(m.(*model).textInput.Value()), nil
}
