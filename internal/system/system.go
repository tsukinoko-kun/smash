package system

import (
	"os"
	"os/user"
	"regexp"
	"runtime/debug"
	"smash/internal/color"
	"strings"
)

const (
	fetchInfoCount = 5
)

var (
	ansiRe = regexp.MustCompile(`\x1b\[[0-9]*m`)
)

func visLen(s string) int {
	return len(ansiRe.ReplaceAllString(s, ""))
}

func SmashFetch() string {
	lines := strings.Split(Ascii, "\n")
	longest := 0
	for _, line := range lines {
		l := visLen(line)
		if l > longest {
			longest = l
		}
	}
	var sb strings.Builder
	i := 0
	for _, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
		r := longest - visLen(line) + 1
		sb.WriteString(strings.Repeat(" ", r))
		if info, ok := fetchInfo(i); ok {
			sb.WriteString(info)
			i++
		}
	}
	for i < fetchInfoCount {
		if info, ok := fetchInfo(i); ok {
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat(" ", len(hr)))
			sb.WriteString(info)
		} else {
			break
		}
		i++
	}
	return sb.String()
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
	default:
		return "", false
	}
}
