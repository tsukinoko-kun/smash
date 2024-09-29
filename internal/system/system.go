package system

import (
	"regexp"
	"strings"
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
			sb.WriteString(strings.Repeat(" ", longest+1))
			sb.WriteString(info)
		} else {
			break
		}
		i++
	}
	return sb.String()
}
