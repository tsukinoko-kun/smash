//go:build windows

package system

import (
	"os"
	"os/user"
	"runtime/debug"
	"smash/internal/color"
	"strings"
)

const (
	Name  = "Windows"
	Ascii = color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		"\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n"
	DefaultShell   = "cmd.exe"
	fetchInfoCount = 5
)

var hr string

func fetchInfo(i int) (string, bool) {
	switch i {
	case 0:
		username := "unknown"
		hostname := "unknown"
		if u, err := user.Current(); err == nil {
			username = u.Username
		}
		if h, err := os.Hostname(); err == nil {
			hostname = h
		}
		x := color.FgMagenta + username + color.Reset + "@" + color.FgMagenta + hostname + color.Reset
		hr = strings.Repeat("-", len(username)+len(hostname)+1)
		return x, true
	case 1:
		return color.Reset + hr, true
	case 2:
		return color.FgYellow + "OS" + color.Reset + ": " + Name, true
	case 3:
		if info, ok := debug.ReadBuildInfo(); ok {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash " + info.Main.Version, true
		} else {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash", true
		}
	case 4:
		return color.FgYellow + "Default Shell" + color.Reset + ": " + DefaultShell, true
	default:
		return "", false
	}
}
