package parser

import (
	"strings"
)

func parseEnv(input string) (string, string, bool) {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
