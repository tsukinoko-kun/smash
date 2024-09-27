package parser_test

import (
	"smash/internal/shell/parser"
	"testing"
)

func TestTokenize1(t *testing.T) {
	tokens := parser.Tokenize(`ls -la && echo "foo 'bar'"`)
	if len(tokens) != 5 {
		t.Fatal("expected 5 tokens, got", len(tokens))
	}

	expected := []parser.Token{
		{Content: "ls", Type: parser.TokenTypeName, Start: 0, End: 1},
		{Content: "-la", Type: parser.TokenTypeName, Start: 3, End: 5},
		{Content: "&&", Type: parser.TokenTypeOperator, Start: 7, End: 8},
		{Content: "echo", Type: parser.TokenTypeName, Start: 10, End: 13},
		{Content: `foo 'bar'`, Type: parser.TokenTypeName, Start: 15, End: 25},
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

func TestTokenize2(t *testing.T) {
	tokens := parser.Tokenize(`  ls  -la  &&  echo  "foo `)
	if len(tokens) != 5 {
		t.Fatal("expected 5 tokens, got", len(tokens))
	}

	expected := []parser.Token{
		{Content: "ls", Type: parser.TokenTypeName, Start: 2, End: 3},
		{Content: "-la", Type: parser.TokenTypeName, Start: 6, End: 8},
		{Content: "&&", Type: parser.TokenTypeOperator, Start: 11, End: 12},
		{Content: "echo", Type: parser.TokenTypeName, Start: 15, End: 18},
		{Content: "foo ", Type: parser.TokenTypeName, Start: 21, End: 25},
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
