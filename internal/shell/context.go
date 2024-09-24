package shell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"smash/internal/color"
	"strings"
)

type (
	context struct {
		check []contextCheck
	}

	contextCheck interface {
		String() string
		Ok() bool
		Reset()
		Process(entry []os.DirEntry)
	}

	contextCheckName struct {
		display string
		names   []string
		ok      bool
	}

	contextCheckSuffix struct {
		display string
		names   []string
		ok      bool
	}
)

func (c *contextCheckName) String() string {
	return c.display
}

func (c *contextCheckName) Ok() bool {
	return c.ok
}

func (c *contextCheckName) Reset() {
	c.ok = false
}

func (c *contextCheckName) Process(entry []os.DirEntry) {
	if c.ok {
		return
	}
	for _, f := range entry {
		if slices.Contains(c.names, f.Name()) {
			c.ok = true
			return
		}
	}
}

func (c *contextCheckSuffix) String() string {
	return c.display
}

func (c *contextCheckSuffix) Ok() bool {
	return c.ok
}

func (c *contextCheckSuffix) Reset() {
	c.ok = false
}

func (c *contextCheckSuffix) Process(entry []os.DirEntry) {
	if c.ok {
		return
	}
	for _, f := range entry {
		for _, n := range c.names {
			if strings.HasSuffix(f.Name(), n) {
				c.ok = true
				return
			}
		}
	}
}

var baseContext context

func init() {
	baseContext = context{
		check: []contextCheck{
			&contextCheckName{
				display: color.FgBlue + "\uE627",
				names:   []string{"go.mod", "go.sum"},
			},
			&contextCheckName{
				display: color.FgYellow + "\U000F031E",
				names:   []string{"package.json"},
			},
			&contextCheckName{
				display: color.FgMagenta + "\uF005",
				names:   []string{"gleam.toml"},
			},
			&contextCheckSuffix{
				display: color.FgCyan + "\U000F0868",
				names:   []string{"docker-compose.yml", "docker-compose.yaml", "Dockerfile"},
			},
			&contextCheckSuffix{
				display: color.FgMagenta + "\U000F0AAE",
				names:   []string{".sln", ".csproj", ".fsproj"},
			},
			&contextCheckName{
				display: color.FgRed + "\uF323",
				names:   []string{"Cargo.toml", "Cargo.lock"},
			},
			&contextCheckName{
				display: color.FgRed + "\uE739",
				names:   []string{"Gemfile", "config.ru", ".ruby-version"},
			},
			&contextCheckSuffix{
				display: color.FgRed + "\uE737",
				names:   []string{".sbt"},
			},
		},
	}
}

func (c *context) reset() {
	// reset all contextCheck::ok
	for i := range c.check {
		c.check[i].Reset()
	}
}

func (c *context) string() string {
	sb := strings.Builder{}

	cmd := exec.Command("git", "branch", "--show-current")
	if out, err := cmd.Output(); err == nil {
		sb.WriteString(color.FgGreen + "\U000F02A2 " + strings.TrimSpace(string(out)) + " ")
	}

	for _, ch := range c.check {
		if ch.Ok() {
			sb.WriteString(ch.String())
			sb.WriteString(" ")
		}
	}

	sb.WriteString(color.Reset)
	return sb.String()
}

func (c *context) walkDir(root string) error {
	for {
		entries, err := os.ReadDir(root)
		if err != nil {
			return errors.Join(errors.New("failed to read directory"), err)
		}
		for _, ch := range c.check {
			ch.Process(entries)
		}

		nextRoot := filepath.Dir(root)
		if nextRoot == root {
			break
		}
		root = nextRoot
	}
	return nil
}
