package extensions

import (
	"bufio"
	"os"
	"smash/internal/color"
	"strings"
)

type golang struct {
	mod string
}

const goPromptDisplay = color.FgBlue + "\uE627"

func (g *golang) PromptDisplay() string {
	f, err := os.Open(g.mod)
	if err != nil {
		return goPromptDisplay
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			return goPromptDisplay + " " + strings.TrimSpace(line[3:])
		}
	}
	return goPromptDisplay
}

func (g *golang) FileMatch(f os.DirEntry, p string) bool {
	if f.Name() == "go.mod" {
		g.mod = p
		return true
	}
	return false
}

func init() {
	register(&golang{})
}
