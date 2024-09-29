//go:build darwin

package system

import (
	"os"
	"os/exec"
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
	fetchInfoCount = 9
)

var (
	DefaultShell = "zsh"
)

func init() {
	if shell, ok := os.LookupEnv("SHELL"); ok {
		DefaultShell = filepath.Base(shell)
	}
}

func fetchInfo(i int) (string, bool) {
	switch i {
	case 0:
		return color.FgGreen + username + color.Reset + "@" + color.FgGreen + hostname + color.Reset, true
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
		return color.FgYellow + "Desktop Environment" + color.Reset + ": " + "Aqua", true
	case 6:
		return color.FgYellow + "Window Manager" + color.Reset + ": " + "Quartz Compositor", true
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
	case 8:
		sb := strings.Builder{}
		sb.WriteString(color.FgYellow)
		sb.WriteString("CPU")
		sb.WriteString(color.Reset)
		sb.WriteString(": ")
		cOk := false
		if cpu, err := sysctl("machdep.cpu.brand_string"); err == nil {
			sb.WriteString(cpu)
			sb.WriteString(" ")
			cOk = true
		}
		if cores, err := sysctl("hw.logicalcpu_max"); err == nil {
			sb.WriteString(cores)
			sb.WriteString("-Core")
			cOk = true
		}
		return sb.String(), cOk
	default:
		return "", false
	}
}

func sysctl(name string) (string, error) {
	out, err := exec.Command("sysctl", "-n", name).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
