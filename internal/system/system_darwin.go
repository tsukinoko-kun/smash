//go:build darwin

package system

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"
	"smash/internal/color"
	"strings"
)

const (
	Name  = "macOS"
	Ascii = color.FgGreen + `                     ..'` + "\n" +
		color.FgGreen + `                 ,xNMM.` + "\n" +
		color.FgGreen + `               .OMMMMo` + "\n" +
		color.FgGreen + `               lMM"` + "\n" +
		color.FgGreen + `     .;loddo:.  .olloddol;.` + "\n" +
		color.FgGreen + `   cKMMMMMMMMMMNWMMMMMMMMMM0:` + "\n" +
		color.FgYellow + ` .KMMMMMMMMMMMMMMMMMMMMMMMWd.` + "\n" +
		color.FgYellow + ` XMMMMMMMMMMMMMMMMMMMMMMMX.` + "\n" +
		color.FgRed + `;MMMMMMMMMMMMMMMMMMMMMMMM:` + "\n" +
		color.FgRed + `:MMMMMMMMMMMMMMMMMMMMMMMM:` + "\n" +
		color.FgRed + `.MMMMMMMMMMMMMMMMMMMMMMMMX.` + "\n" +
		color.FgRed + ` kMMMMMMMMMMMMMMMMMMMMMMMMWd.` + "\n" +
		color.FgMagenta + ` 'XMMMMMMMMMMMMMMMMMMMMMMMMMMk` + "\n" +
		color.FgMagenta + `  'XMMMMMMMMMMMMMMMMMMMMMMMMK.` + "\n" +
		color.FgBlue + `    kMMMMMMMMMMMMMMMMMMMMMMd` + "\n" +
		color.FgBlue + `     ;KMMMMMMMWXXWMMMMMMMk.` + "\n" +
		color.FgBlue + `       "cooc*"    "*coo'"` + "\n"
	fetchInfoCount = 7
)

var (
	DefaultShell = "zsh"
)

func init() {
	if shell, ok := os.LookupEnv("SHELL"); ok {
		DefaultShell = filepath.Base(shell)
	}
}

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
	case 5:
		return color.FgYellow + "DE" + color.Reset + ": Aqua" + DefaultShell, true
	case 6:
		return color.FgYellow + "WM" + color.Reset + ": Quartz Compositor", true
	default:
		return "", false
	}
}
