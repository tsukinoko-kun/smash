//go:build windows

package system

import "smash/internal/color"

const (
	Name  = "Windows"
	Ascii = color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		color.FgRed + "lllllll  " + color.FgGreen + "lllllll\n" +
		"\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n" +
		color.FgBlue + "lllllll  " + color.FgYellow + "lllllll\n"
)
