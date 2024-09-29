//go:build windows

package system

import (
	"os"
	"runtime/debug"
	"smash/internal/color"
	"strings"
)

const (
	Ascii = color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		"\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n"
	fetchInfoCount = 8
)

func fetchInfo(i int) (string, bool) {
	switch i {
	case 0:
		return color.FgBlue + username + color.Reset + "@" + color.FgBlue + hostname + color.Reset, true
	case 1:
		return color.Reset + hr, true
	case 2:
		return color.FgYellow + "OS" + color.Reset + ": Windows", true
	case 3:
		if info, ok := debug.ReadBuildInfo(); ok {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash " + info.Main.Version, true
		} else {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash", true
		}
	case 4:
		return color.FgYellow + "Default Shell" + color.Reset + ": cmd.exe", true
	case 5:
		return color.FgYellow + "Desktop Environment" + color.Reset + ": " + "Windows Desktop Environment", true
	case 6:
		return color.FgYellow + "Window Manager" + color.Reset + ": " + "Desktop Window Manager", true
	case 7:
		sb := strings.Builder{}
		sb.WriteString(color.FgYellow)
		sb.WriteString("Terminal")
		sb.WriteString(color.Reset)
		sb.WriteString(": ")
		tOk := false
		if termProgram, ok := os.LookupEnv("TERM_PROGRAM"); ok {
			sb.WriteString(termProgram)
			sb.WriteString(" ")
			tOk = true
		}
		if term, ok := os.LookupEnv("TERM"); ok {
			sb.WriteString(term)
			tOk = true
		}
		return sb.String(), tOk
	default:
		return "", false
	}
}
