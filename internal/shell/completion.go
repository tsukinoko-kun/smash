package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Completion struct {
	Value       string
	DisplayText string
	Delete      int
}

func (c Completion) String() string {
	return fmt.Sprintf(" %s (%s) ", c.Value, c.DisplayText)
}

const completionLimit = 24

func GetCompletions(input string) []Completion {
	fields := strings.Fields(input)
	if input[len(input)-1] == ' ' {
		fields = append(fields, "")
	}
	switch len(fields) {
	case 0:
		return []Completion{}
	case 1:
		return getNameCompletion(input)
	default:
		return findWithPath(fields[len(fields)-1], false)
	}
}

func getNameCompletion(input string) []Completion {
	completions := findWithPath(input, true)

	for _, tool := range internalToolNames {
		if strings.HasPrefix(tool, input) {
			completions = append(completions, Completion{Value: tool, DisplayText: "shell internal tool", Delete: len(input)})
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
				completions = append(completions, Completion{Value: n, DisplayText: "executable in PATH", Delete: len(input)})
				if len(completions) >= completionLimit {
					break outer
				}
			}
		}
	}
	return completions
}

func findWithPath(partialPath string, executable bool) []Completion {
	var completions []Completion

	expandedPartialPath := os.ExpandEnv(partialPath)
	if strings.HasPrefix(expandedPartialPath, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			expandedPartialPath = home + expandedPartialPath[1:]
		}
	}
	if strings.HasSuffix(partialPath, "/") {
		if fi, err := os.Stat(expandedPartialPath); err != nil || !fi.IsDir() {
			return completions
		}

		// this is a directory, get its content

		entries, err := os.ReadDir(expandedPartialPath)
		if err != nil {
			return completions
		}
		for _, f := range entries {
			c := Completion{
				Value:  f.Name(),
				Delete: 0,
			}

			if f.IsDir() {
				c.DisplayText = "directory"
				c.Value += "/"
			} else {
				if executable {
					if fi, err := f.Info(); err != nil || fi.Mode()&0111 == 0 {
						continue
					}
					c.DisplayText = "executable file"
				} else {
					c.DisplayText = "file"
				}
			}

			completions = append(completions, c)
			if len(completions) >= completionLimit {
				return completions
			}
		}
	} else {
		// get matching paths inside the parent directory

		parentDir := filepath.Dir(expandedPartialPath)
		entries, err := os.ReadDir(parentDir)
		if err != nil {
			return completions
		}
		partialFile := filepath.Base(expandedPartialPath)

		for _, f := range entries {
			if !strings.HasPrefix(f.Name(), partialFile) {
				continue
			}

			c := Completion{
				Value:  f.Name(),
				Delete: len(partialFile),
			}

			if f.IsDir() {
				c.DisplayText = "directory"
				c.Value += "/"
			} else {
				if executable {
					if fi, err := f.Info(); err != nil || fi.Mode()&0111 == 0 {
						continue
					}
					c.DisplayText = "executable file"
					if partialPath[0] != '.' && partialPath[0] != '/' {
						c.Value = "./" + c.Value
					}
				} else {
					c.DisplayText = "file"
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
