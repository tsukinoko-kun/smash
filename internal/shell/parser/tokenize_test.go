package parser_test

import (
	"smash/internal/shell/parser"
	"testing"
)

func checkTokens(t *testing.T, tokens []parser.Token, expected []parser.Token) {
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, token := range tokens {
		if token.Content != expected[i].Content {
			t.Fatalf("expected %s, got %s", expected[i].Content, token.Content)
		}
		if token.Type != expected[i].Type {
			t.Fatalf("expected %d, got %d", expected[i].Type, token.Type)
		}
		if token.Start != expected[i].Start {
			t.Fatalf("expected %d, got %d", expected[i].Start, token.Start)
		}
		if token.End != expected[i].End {
			t.Fatalf("expected %d, got %d", expected[i].End, token.End)
		}
	}
}

func TestTokenize1(t *testing.T) {
	t.Parallel()
	tokens := parser.Tokenize(`ls -la && echo "foo 'bar'"`)

	expected := []parser.Token{
		{Content: "ls", Type: parser.TokenTypeName, Start: 0, End: 1},
		{Content: "-la", Type: parser.TokenTypeName, Start: 3, End: 5},
		{Content: "&&", Type: parser.TokenTypeOperator, Start: 7, End: 8},
		{Content: "echo", Type: parser.TokenTypeName, Start: 10, End: 13},
		{Content: `foo 'bar'`, Type: parser.TokenTypeName, Start: 15, End: 25},
	}

	checkTokens(t, tokens, expected)
}

func TestTokenize2(t *testing.T) {
	t.Parallel()
	tokens := parser.Tokenize(`  ls  -la  &&  echo  "foo `)

	expected := []parser.Token{
		{Content: "ls", Type: parser.TokenTypeName, Start: 2, End: 3},
		{Content: "-la", Type: parser.TokenTypeName, Start: 6, End: 8},
		{Content: "&&", Type: parser.TokenTypeOperator, Start: 11, End: 12},
		{Content: "echo", Type: parser.TokenTypeName, Start: 15, End: 18},
		{Content: "foo ", Type: parser.TokenTypeName, Start: 21, End: 25},
	}

	checkTokens(t, tokens, expected)
}

func TestTokenize3(t *testing.T) {
	t.Parallel()
	tokens := parser.Tokenize(`FOO=foobar echo $FOO`)

	expected := []parser.Token{
		{Content: "FOO=foobar", Type: parser.TokenTypeEnvironment, Start: 0, End: 9},
		{Content: "echo", Type: parser.TokenTypeName, Start: 11, End: 14},
		{Content: "foobar", Type: parser.TokenTypeName, Start: 16, End: 19},
	}

	checkTokens(t, tokens, expected)
}
