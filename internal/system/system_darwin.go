//go:build darwin

package system

import "smash/internal/color"

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
)
