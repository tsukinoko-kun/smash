package extensions

import (
	"smash/internal/assert"
	"smash/internal/shell/parser"
)

type (
	Completion struct {
		Value        string
		DisplayValue string
		Description  string
	}
)

func GetCompletion(tokens []parser.Token, newArg bool) []Completion {
	assert.GreaterThan(len(tokens), 0, "tokens for GetCompletion must not be empty")
	name := tokens[0].Content
	var args []string
	if len(tokens) > 1 {
		for _, t := range tokens[1:] {
			args = append(args, t.Content)
		}
	}
	if newArg {
		args = append(args, "")
	}

	var completions []Completion
	for _, ew := range extensionRegister {
		if ew.active {
			for _, s := range ew.extension.Completions(name, args) {
				completions = append(completions, Completion{Value: s, DisplayValue: s, Description: ew.extension.Title()})
			}
		}
	}
	return completions
}
