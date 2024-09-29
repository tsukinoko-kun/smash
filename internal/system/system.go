package system

import (
	"os"
	"os/user"
	"regexp"
	"strings"
)

var (
	ansiRe   = regexp.MustCompile(`\x1b\[[0-9]*m`)
	username string
	hostname string
	hr       string
)

func init() {
	if u, err := user.Current(); err == nil {
		username = u.Username
	}
	if h, err := os.Hostname(); err == nil {
		hostname = h
	}
	hr = strings.Repeat("-", len(username)+len(hostname)+1)
}

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
			sb.WriteString(strings.Repeat(" ", longest+1))
			sb.WriteString(info)
		} else {
			break
		}
		i++
	}
	return sb.String()
}
