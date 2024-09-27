package extensions

import (
	"os"
	"os/exec"
	"smash/internal/color"
	"strings"
)

type git struct {
}

func (_ *git) PromptDisplay() string {
	cmd := exec.Command("git", "branch", "--show-current")
	if out, err := cmd.Output(); err == nil {
		return color.FgGreen + "\U000F02A2 " + strings.TrimSpace(string(out))
	} else {
		return color.FgGreen + "\U000F02A2"
	}
}

func (_ *git) FileMatch(f os.DirEntry, _ string) bool {
	return f.Name() == ".git"
}

func (_ *git) Completions(name string, args []string) []string {
	return nil
}

func (_ *git) Title() string {
	return "Git"
}

func init() {
	priorityRegister(&git{})
}
