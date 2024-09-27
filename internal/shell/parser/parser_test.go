package parser_test

import (
	"smash/internal/shell/parser"
	"smash/internal/tess"
	"testing"
)

func TestParseEmptyString(t *testing.T) {
	t.Parallel()
	cs, err := parser.Parse("")
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cs))
	}
	c := cs[0]
	name := tess.UnsafeCast[string](tess.Field(c, "Name"))
	if name != "" {
		t.Fatalf("expected empty name, got %s", name)
	}
}

func TestOneCommandNoArgs(t *testing.T) {
	t.Parallel()
	cs, err := parser.Parse("ls")
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cs))
	}

	c := cs[0]
	name := tess.UnsafeCast[string](tess.Field(c, "Name"))
	if name != "ls" {
		t.Fatalf("expected ls, got %s", name)
	}

	args := tess.UnsafeCast[[]string](tess.Field(c, "Args"))
	if len(args) != 0 {
		t.Fatalf("expected 0 args, got %d", len(args))
	}
}

func TestOneCommandOneArg(t *testing.T) {
	t.Parallel()
	cs, err := parser.Parse("ls -l")
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cs))
	}

	c := cs[0]
	name := tess.UnsafeCast[string](tess.Field(c, "Name"))
	if name != "ls" {
		t.Fatalf("expected ls, got %s", name)
	}

	args := tess.UnsafeCast[[]string](tess.Field(c, "Args"))
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}

	if args[0] != "-l" {
		t.Fatalf("expected -l, got %s", args[0])
	}
}

func TestOneCommandManyArgs(t *testing.T) {
	t.Parallel()
	cs, err := parser.Parse("ls -l -a -h")
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cs))
	}

	c := cs[0]
	name := tess.UnsafeCast[string](tess.Field(c, "Name"))
	if name != "ls" {
		t.Fatalf("expected ls, got %s", name)
	}

	args := tess.UnsafeCast[[]string](tess.Field(c, "Args"))
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}

	if args[0] != "-l" {
		t.Fatalf("expected -l, got %s", args[0])
	}

	if args[1] != "-a" {
		t.Fatalf("expected -a, got %s", args[1])
	}

	if args[2] != "-h" {
		t.Fatalf("expected -h, got %s", args[2])
	}
}

func TestAndCommands(t *testing.T) {
	t.Parallel()
	cs, err := parser.Parse("cd /bin && ls -la")
	if err != nil {
		t.Fatal(err)
	}

	if len(cs) != 1 {
		t.Fatalf("expected 1 commands, got %d", len(cs))
	}

	andCommands := tess.UnsafeCast[[]parser.Command](tess.Field(cs[0], "Commands"))
	if len(andCommands) != 2 {
		t.Fatalf("expected 2 && commands, got %d", len(andCommands))
	}

	c1 := andCommands[0]
	name := tess.UnsafeCast[string](tess.Field(c1, "Name"))
	if name != "cd" {
		t.Fatalf("expected cd, got %s", name)
	}
	args1 := tess.UnsafeCast[[]string](tess.Field(c1, "Args"))
	if len(args1) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args1))
	}
	if args1[0] != "/bin" {
		t.Fatalf("expected /bin, got %s", args1[0])
	}

	c2 := andCommands[1]
	name = tess.UnsafeCast[string](tess.Field(c2, "Name"))
	if name != "ls" {
		t.Fatalf("expected ls, got %s", name)
	}
	args2 := tess.UnsafeCast[[]string](tess.Field(c2, "Args"))
	if len(args2) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args2))
	}
	if args2[0] != "-la" {
		t.Fatalf("expected -la, got %s", args2[0])
	}
}
