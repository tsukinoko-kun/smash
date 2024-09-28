package system

import (
	"runtime/debug"
	"smash/internal/color"
	"strings"
)

const (
	fetchInfoCount = 3
)

func SmashFetch() string {
	lines := strings.Split(Ascii, "\n")
	longest := 0
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	var sb strings.Builder
	i := 0
	for _, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
		for j := len(line); j < longest; j++ {
			sb.WriteString(" ")
		}
		if info, ok := fetchInfo(i); ok {
			sb.WriteString(info)
			i++
		}
	}
	for i < fetchInfoCount {
		if info, ok := fetchInfo(i); ok {
			for j := 0; j <= longest; j++ {
				sb.WriteString(" ")
			}
			sb.WriteString(info)
			sb.WriteString("\n")
		} else {
			break
		}
		i++
	}
	return sb.String()
}

func fetchInfo(i int) (string, bool) {
	switch i {
	case 0:
		return color.FgYellow + "OS" + color.Reset + ": " + Name, true
	case 1:
		if info, ok := debug.ReadBuildInfo(); ok {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash " + info.Main.Version, true
		} else {
			return color.FgYellow + "Interactive Shell" + color.Reset + ": smash", true
		}
	case 2:
		return color.FgYellow + "Default Shell" + color.Reset + ": " + DefaultShell, true
	default:
		return "", false
	}
}
