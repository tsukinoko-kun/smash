package parser

import (
	"os"
	"path/filepath"
	"smash/internal/assert"
	"smash/internal/env"
	"strings"
)

type (
	TokenType    uint8
	TokenBuilder struct {
		Content strings.Builder
		Start   int
		End     int
	}
	Token struct {
		Type    TokenType
		Content string
		Start   int
		End     int
	}
)

const (
	TokenTypeName = iota
	TokenTypeOperator
)

func newToken(start int) *TokenBuilder {
	return &TokenBuilder{
		Content: strings.Builder{},
		Start:   start,
		End:     -1,
	}
}

func Tokenize(input string) []Token {
	var tokens []Token
	tb := newToken(0)
	quotation := rune(0)

	next := func(end int, nextStart int) {
		tb.End = end
		for _, s := range expand(tb.Content.String(), quotation) {
			t := Token{Content: s, Start: tb.Start, End: tb.End}
			for _, op := range operators {
				if t.Content == op {
					t.Type = TokenTypeOperator
					break
				}
			}
			tokens = append(tokens, t)
		}
		tb = newToken(nextStart)
	}

	for i, char := range []rune(input) {
		switch char {
		case ' ', '\t':
			if quotation != 0 {
				tb.Content.WriteRune(char)
			} else if tb.Content.Len() != 0 {
				next(i-1, i+1)
			} else {
				tb.Start++
			}
		case '}':
			if quotation == '{' {
				quotation = rune(0)
			}
			tb.Content.WriteRune(char)
		case '{':
			if quotation == 0 {
				quotation = char
			}
			tb.Content.WriteRune(char)
		case '"', '\'', '`':
			if quotation == char {
				next(i, i+1)
				quotation = rune(0)
			} else if quotation == 0 {
				quotation = char
			} else {
				tb.Content.WriteRune(char)
			}
		default:
			tb.Content.WriteRune(char)
		}
	}

	if tb.Content.Len() != 0 {
		tb.End = len(input) - 1
		for _, s := range expand(tb.Content.String(), quotation) {
			t := Token{Content: s, Start: tb.Start, End: tb.End}
			for _, op := range operators {
				if t.Content == op {
					t.Type = TokenTypeOperator
					break
				}
			}
			tokens = append(tokens, t)
		}
	}

	// check tokens
	for _, t := range tokens {
		assert.GreaterThan(t.Start, -1, "token start must be greater than -1")
		assert.SmallerThan(t.End, len(input), "token end must be smaller than input length")
	}

	return tokens
}

func expand(s string, quotation rune) []string {
	switch quotation {
	case '\'':
		return []string{s}
	case '"':
		return []string{expandVariables(expandTilde(s))}
	default:
		return expandGlobs(expandVariables(expandTilde(s)))
	}
}

func expandTilde(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '~' {
		ss := env.GetUser().HomeDir + s[1:]
		if _, err := os.Stat(ss); err == nil {
			s = ss
		}
	}
	return s
}

func expandVariables(s string) string {
	return os.ExpandEnv(s)
}

func expandGlobs(s string) []string {
	if matches, err := filepath.Glob(s); err == nil && len(matches) > 0 {
		return matches
	}
	return []string{s}
}

var operators = [...]string{"&&", "&", ";", "||", "|"}
