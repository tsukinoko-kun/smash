package shell

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"smash/internal/color"
	"smash/internal/env"
	"strings"
	"sync"
)

var u *user.User

func init() {
	var err error
	u, err = user.Current()
	if err != nil {
		panic(err)
	}
}

type dirInfo struct {
	Wd  string
	Dev string
}

func getDirInfo(wd string) dirInfo {
	di := dirInfo{
		Wd: wd,
	}

	wg := sync.WaitGroup{}
	dev := [4]string{}
	wg.Add(len(dev))
	go func() {
		defer wg.Done()
		dev[0] = dirInfoGit(wd)
	}()
	go func() {
		defer wg.Done()
		dev[1] = dirInfoDocker(wd)
	}()
	go func() {
		defer wg.Done()
		dev[2] = dirInfoGo(wd)
	}()
	go func() {
		defer wg.Done()
		dev[3] = dirInfoJavascript(wd)
	}()
	wg.Wait()

	sb := strings.Builder{}
	for _, d := range dev {
		if d != "" {
			sb.WriteString(d)
			sb.WriteString(" ")
		}
	}
	di.Dev = sb.String()
	return di
}

func dirInfoGit(wd string) string {
	// run git to get current branch
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = wd
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return color.FgGreen + "\U000F02A2 " + strings.TrimSpace(string(out)) + color.Reset
}

func dirInfoDocker(wd string) string {
	found := false
	dirEntries, err := os.ReadDir(wd)
	if err != nil {
		return ""
	}
	for _, d := range dirEntries {
		if d.IsDir() {
			continue
		}
		if strings.HasSuffix(d.Name(), "Dockerfile") ||
			strings.HasSuffix(d.Name(), "docker-compose.yml") ||
			strings.HasSuffix(d.Name(), "docker-compose.toml") {
			found = true
			break
		}
	}
	if found {
		return color.FgCyan + "\U000F0868" + color.Reset
	} else {
		return ""
	}
}

func dirInfoGo(wd string) string {
	goModPath := filepath.Join(wd, "go.mod")
	f, err := os.Open(goModPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	// try to find versino e.g. go 1.23.1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			return color.FgBlue + "\uE627 " + line[3:] + color.Reset
		}
	}
	return ""
}

func dirInfoJavascript(wd string) string {
	if _, err := os.Stat(filepath.Join(wd, "package.json")); err == nil {
		return color.FgYellow + "\U000F031E" + color.Reset
	}
	return ""
}

func (di dirInfo) Expand(s string) string {
	switch s {
	case "DEV":
		return di.Dev
	case "PWD":
		return di.Wd
	case "USER":
		return u.Username

	// Color Reset
	case "Color.Reset":
		return color.Reset

	// Color Foreground
	case "Color.FgBlack":
		return color.FgBlack
	case "Color.FgRed":
		return color.FgRed
	case "Color.FgGreen":
		return color.FgGreen
	case "Color.FgYellow":
		return color.FgYellow
	case "Color.FgBlue":
		return color.FgBlue
	case "Color.FgMagenta":
		return color.FgMagenta
	case "Color.FgCyan":
		return color.FgCyan
	case "Color.FgWhite":
		return color.FgWhite

	case "Color.FgHiBlack":
		return color.FgHiBlack
	case "Color.FgHiRed":
		return color.FgHiRed
	case "Color.FgHiGreen":
		return color.FgHiGreen
	case "Color.FgHiYellow":
		return color.FgHiYellow
	case "Color.FgHiBlue":
		return color.FgHiBlue
	case "Color.FgHiMagenta":
		return color.FgHiMagenta
	case "Color.FgHiCyan":
		return color.FgHiCyan
	case "Color.FgHiWhite":
		return color.FgHiWhite

	// Color Background
	case "Color.BgBlack":
		return color.BgBlack
	case "Color.BgRed":
		return color.BgRed
	case "Color.BgGreen":
		return color.BgGreen
	case "Color.BgYellow":
		return color.BgYellow
	case "Color.BgBlue":
		return color.BgBlue
	case "Color.BgMagenta":
		return color.BgMagenta
	case "Color.BgCyan":
		return color.BgCyan
	case "Color.BgWhite":
		return color.BgWhite

	case "Color.BgHiBlack":
		return color.BgHiBlack
	case "Color.BgHiRed":
		return color.BgHiRed
	case "Color.BgHiGreen":
		return color.BgHiGreen
	case "Color.BgHiYellow":
		return color.BgHiYellow
	case "Color.BgHiBlue":
		return color.BgHiBlue
	case "Color.BgHiMagenta":
		return color.BgHiMagenta
	case "Color.BgHiCyan":
		return color.BgHiCyan
	case "Color.BgHiWhite":
		return color.BgHiWhite

	}
	return os.Getenv(s)
}

func Ps1() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Join(errors.New("could not get working directory"), err)
	}
	di := getDirInfo(wd)
	if strings.HasPrefix(wd, u.HomeDir) {
		di.Wd = filepath.Join("~", wd[len(u.HomeDir):])
	}
	return os.Expand(env.Config.InteractivePrompt, di.Expand), nil
}

func Ps2() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Join(errors.New("could not get working directory"), err)
	}
	di := getDirInfo(wd)
	if strings.HasPrefix(wd, u.HomeDir) {
		di.Wd = filepath.Join("~", wd[len(u.HomeDir):])
	}
	return os.Expand(env.Config.LogPrompt, di.Expand), nil
}
