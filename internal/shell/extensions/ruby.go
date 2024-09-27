package extensions

import (
	"bufio"
	"fmt"
	"os"
	"smash/internal/color"
)

type ruby struct {
	gemfile string
}

const rubyPromptDisplay = color.FgRed + "\uE739"

func (ru *ruby) PromptDisplay() string {
	f, err := os.Open(ru.gemfile)
	if err != nil {
		return rubyPromptDisplay
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		var version string
		if _, err := fmt.Sscanf(line, `ruby "%s"`, &version); err == nil {
			return rubyPromptDisplay + " " + version
		}
	}
	return rubyPromptDisplay
}

func (ru *ruby) FileMatch(f os.DirEntry, p string) bool {
	if f.Name() == "Gemfile" {
		ru.gemfile = p
		return true
	}
	return false
}

func (_ *ruby) Completions(name string, args []string) []string {
	return nil
}

func (_ *ruby) Title() string {
	return "Ruby"
}

func init() {
	register(&ruby{})
}
