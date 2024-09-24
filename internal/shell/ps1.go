package shell

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"smash/internal/color"
	"smash/internal/env"
	"strings"
	"time"
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

	baseContext.reset()
	_ = baseContext.walkDir(wd)
	di.Dev = baseContext.string()

	return di
}

func (di dirInfo) Expand(s string) string {
	switch s {
	case "DEV":
		return di.Dev
	case "PWD":
		return di.Wd
	case "USER":
		return u.Username
	case "HOSTNAME":
		h, _ := os.Hostname()
		return h
	case "TIME":
		return time.Now().String()

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
