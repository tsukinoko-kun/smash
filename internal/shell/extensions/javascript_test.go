package extensions

import (
	"os"
	"slices"
	"testing"
)

func testCompletion(t *testing.T, ex Extension, name string, args []string, expected []string) {
	res := ex.Completions(name, args)
	if !slices.Equal(res, expected) {
		t.Errorf("completion of `%s %s` returned %v expected %v", name, args, res, expected)
	}
}

func TestJavascript_Completions(t *testing.T) {
	t.Parallel()
	js := &javascript{}

	f, err := os.CreateTemp("", "*.package.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"scripts": {"build": "echo build", "dev": "echo dev"}}`)
	_ = f.Close()
	js.packageJson = f.Name()
	defer os.Remove(js.packageJson)

	testCompletion(t, js, "npm", []string{"ru"}, []string{"run"})
	testCompletion(t, js, "npm", []string{""}, []string{"install", "run"})
	testCompletion(t, js, "npm", []string{"run", ""}, []string{"build", "dev"})
	testCompletion(t, js, "npm", []string{"run", "bui"}, []string{"build"})
}

func TestJavascript_CompletionsEmpty(t *testing.T) {
	t.Parallel()
	js := &javascript{}

	f, err := os.CreateTemp("", "*.package.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{}`)
	_ = f.Close()
	js.packageJson = f.Name()
	defer os.Remove(js.packageJson)

	testCompletion(t, js, "npm", []string{"ru"}, []string{"run"})
	testCompletion(t, js, "npm", []string{""}, []string{"install", "run"})
	testCompletion(t, js, "npm", []string{"run", ""}, []string{})
	testCompletion(t, js, "npm", []string{"run", "bui"}, []string{})
}
