package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"smash/internal/shell/extensions"
	"smash/internal/shell/parser"
	"strings"
)

const completionLimit = 24

func isEmptyOrWhitespace(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(" \t\n\r", r) {
			return false
		}
	}
	return true
}

func quoteArgIfNeeded(arg string) string {
	if strings.ContainsAny(arg, " \t\n\r\\\"'`\v\f") {
		return fmt.Sprintf("%q", arg)
	}
	return arg
}

func GetCompletions(input string) []extensions.Completion {
	if isEmptyOrWhitespace(input) {
		return nil
	}

	tokens := parser.Tokenize(input)
	l := len(tokens)

	if l == 0 {
		return nil
	}

	if tokens[l-1].Type == parser.TokenTypeOperator {
		return nil
	}

	// find start of name series
	start := len(tokens) - 1
	for start > 0 {
		if tokens[start-1].Type == parser.TokenTypeName {
			start--
		} else {
			break
		}
	}

	tokensForCompletion := tokens[start:]

	newArg := input[len(input)-1] == ' '

	switch len(tokensForCompletion) {
	case 0:
		return nil
	case 1:
		t := tokensForCompletion[0]
		cs := extensions.GetCompletion(tokensForCompletion, newArg)
		if len(cs) != 0 {
			if newArg {
				for i, c := range cs {
					c.Value = input + quoteArgIfNeeded(c.Value)
					cs[i] = c
				}
			} else {
				for i, c := range cs {
					c.Value = input[:t.Start] + quoteArgIfNeeded(c.Value)
					cs[i] = c
				}
			}
			return cs
		}
		if newArg {
			cs = findWithPath("", false)
			for i, c := range cs {
				c.Value = input + quoteArgIfNeeded(c.Value)
				cs[i] = c
			}
		} else {
			cs = getNameCompletion(t.Content)
			for i, c := range cs {
				c.Value = input[:t.Start] + quoteArgIfNeeded(c.Value)
				cs[i] = c
			}
		}
		return cs
	default:
		t := tokensForCompletion[len(tokensForCompletion)-1]
		cs := extensions.GetCompletion(tokensForCompletion, newArg)
		if len(cs) != 0 {
			if newArg {
				for i, c := range cs {
					c.Value = input + quoteArgIfNeeded(c.Value)
					cs[i] = c
				}
			} else {
				for i, c := range cs {
					c.Value = input[:t.Start] + quoteArgIfNeeded(c.Value)
					cs[i] = c
				}
			}
			return cs
		}
		if newArg {
			cs = findWithPath("", false)
			for i, c := range cs {
				c.Value = input + quoteArgIfNeeded(c.Value)
				cs[i] = c
			}
		} else {
			cs = findWithPath(t.Content, false)
			for i, c := range cs {
				c.Value = input[:t.Start] + quoteArgIfNeeded(c.Value)
				cs[i] = c
			}
		}
		return cs
	}
}

func getNameCompletion(input string) []extensions.Completion {
	completions := findWithPath(input, true)

	for _, tool := range parser.InternalToolNames {
		if strings.HasPrefix(tool, input) {
			completions = append(completions, extensions.Completion{DisplayValue: tool, Value: tool, Description: "shell internal tool"})
			if len(completions) >= completionLimit {
				return completions
			}
		}
	}
	if path, ok := os.LookupEnv("PATH"); ok {
	outer:
		for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
			d, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, entry := range d {
				if entry.IsDir() {
					continue
				}
				if fileInfo, err := entry.Info(); err != nil || fileInfo.Mode()&0111 == 0 {
					continue
				}
				n := entry.Name()
				if !strings.HasPrefix(n, input) {
					continue
				}
				completions = append(completions, extensions.Completion{DisplayValue: n, Value: n, Description: "executable in PATH"})
				if len(completions) >= completionLimit {
					break outer
				}
			}
		}
	}
	return completions
}

func findWithPath(partialPath string, executable bool) []extensions.Completion {
	var completions []extensions.Completion

	if fi, err := os.Stat(partialPath); err == nil && fi.IsDir() && strings.HasSuffix(partialPath, "/") {
		// this is a directory, get its content

		entries, err := os.ReadDir(partialPath)
		if err != nil {
			return completions
		}
		for _, f := range entries {
			c := extensions.Completion{
				DisplayValue: f.Name(),
				Value:        filepath.Join(partialPath, f.Name()),
			}

			if f.IsDir() {
				c.Description = "directory"
				c.Value += "/"
			} else {
				if executable {
					if fi, err := f.Info(); err != nil || fi.Mode()&0111 == 0 {
						continue
					}
					c.Description = "executable file"
				} else {
					c.Description = "file"
				}
			}

			completions = append(completions, c)
			if len(completions) >= completionLimit {
				return completions
			}
		}
	} else {
		// get matching paths inside the parent directory

		parentDir := filepath.Dir(partialPath)
		entries, err := os.ReadDir(parentDir)
		if err != nil {
			return completions
		}
		var partialFile string
		if len(partialPath) != 0 {
			partialFile = filepath.Base(partialPath)
		} else {
			partialFile = ""
		}

		for _, f := range entries {
			if !strings.HasPrefix(f.Name(), partialFile) {
				continue
			}

			c := extensions.Completion{
				DisplayValue: f.Name(),
				Value:        filepath.Join(parentDir, f.Name()),
			}

			if f.IsDir() {
				c.Description = "directory"
				c.Value += "/"
			} else {
				if executable {
					if fi, err := f.Info(); err != nil || fi.Mode()&0111 == 0 {
						continue
					}
					c.Description = "executable file"
					if partialPath[0] != '.' && partialPath[0] != '/' {
						c.Value = quoteArgIfNeeded("./" + c.Value)
					}
				} else {
					c.Description = "file"
				}
			}

			completions = append(completions, c)
			if len(completions) >= completionLimit {
				return completions
			}
		}
	}

	return completions
}
