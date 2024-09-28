//go:build linux

package system

import (
	"bufio"
	"os"
	"path/filepath"
	"smash/internal/color"
	"strings"
)

var (
	Name         = "Linux"
	id           string
	idLike       string
	Ascii        string
	DefaultShell = "sh"
)

const (
	asciiLogoAlpine = color.FgBlue + `   /\ /\` + "\n" +
		color.FgBlue + `  // \  \` + "\n" +
		color.FgBlue + ` //   \  \` + "\n" +
		color.FgBlue + `///    \  \` + "\n" +
		color.FgBlue + `//      \  \` + "\n" +
		color.FgBlue + `         \` + "\n"
	asciiLogoArch = color.FgBlue + `      /\` + "\n" +
		color.FgBlue + `     /  \` + "\n" +
		color.FgBlue + `    /    \` + "\n" +
		color.FgBlue + `   /      \` + "\n" +
		color.FgBlue + `  /   ,,   \` + "\n" +
		color.FgBlue + ` /   |  |   \` + "\n" +
		color.FgBlue + `/_-''    ''-_\` + "\n"
	asciiLogoBsd = color.FgRed + "             ,        ,\n" +
		color.FgRed + "            /(        )`\n" +
		color.FgRed + "            \\ \\___   / |\n" +
		color.FgRed + "            /- _  `-/  '\n" +
		color.FgRed + "           (/\\/ \\ \\   /\\\n" +
		color.FgRed + "           / /   | `    \\\n" +
		color.FgRed + "           O O   ) /    |\n" +
		color.FgRed + "           `-^--'`<     '\n" +
		color.FgRed + "          (_.)  _  )   /\n" +
		color.FgRed + "           `.___/`    /\n" +
		color.FgRed + "             `-----' /\n" +
		color.FgRed + "<----.     __ / __   \\\n" +
		color.FgRed + "<----|====O)))==) \\) /====|\n" +
		color.FgRed + "<----'    `--' `.__,' \\\n" +
		color.FgRed + "             |        |\n" +
		color.FgRed + "              \\       /       /\\\n" +
		color.FgRed + "         ______( (_  / \\______/\n" +
		color.FgRed + "       ,'  ,-----'   |\n" +
		color.FgRed + "       `--{__________)\n"

	asciiLogoCentos = color.FgGreen + ` ____` + color.FgYellow + `^` + color.FgMagenta + `____` + "\n" +
		color.FgGreen + ` |\  ` + color.FgYellow + `|` + color.FgMagenta + `  /|` + "\n" +
		color.FgGreen + ` | \ ` + color.FgYellow + `|` + color.FgMagenta + ` / |` + "\n" +
		color.FgMagenta + `<---- ` + color.FgBlue + `---->` + "\n" +
		color.FgBlue + ` | / ` + color.FgGreen + `|` + color.FgYellow + ` \ |` + "\n" +
		color.FgBlue + ` |/__` + color.FgGreen + `|` + color.FgYellow + `__\|` + "\n" +
		color.FgGreen + `     v` + "\n"

	asciiLogoDebian = color.FgRed + `  _____` + "\n" +
		color.FgRed + ` /  __ \` + "\n" +
		color.FgRed + `|  /    |` + "\n" +
		color.FgRed + `|  \___-` + "\n" +
		color.FgRed + `-_` + "\n" +
		color.FgRed + `  --_` + "\n"

	asciiLogoElementary = color.FgBlue + `  _______` + "\n" +
		color.FgBlue + ` / ____  \\` + "\n" +
		color.FgBlue + `/  |  /  /\\` + "\n" +
		color.FgBlue + `|__\\ /  / |` + "\n" +
		color.FgBlue + `\\   /__/  /` + "\n" +
		color.FgBlue + ` \\_______/` + "\n"

	asciiLogoFedora = color.FgBlue + `        ,'''''.` + "\n" +
		color.FgBlue + `       |   ,.  |` + "\n" +
		color.FgBlue + `       |  |  '_'` + "\n" +
		color.FgBlue + `  ,....|  |..` + "\n" +
		color.FgBlue + `.'  ,_;|   ..'` + "\n" +
		color.FgBlue + `|  |   |  |` + "\n" +
		color.FgBlue + `|  ',_,'  |` + "\n" +
		color.FgBlue + ` '.     ,'` + "\n" +
		color.FgBlue + `   '''''` + "\n"

	asciiLogoFreebsd = color.FgRed + `/\,-'''''-,/\` + "\n" +
		color.FgRed + `\_)       (_/` + "\n" +
		color.FgRed + `|           |` + "\n" +
		color.FgRed + `|           |` + "\n" +
		color.FgRed + ` ;         ;` + "\n" +
		color.FgRed + `  '-_____-'` + "\n"

	asciiLogoGentoo = color.FgCyan + ` _-----_` + "\n" +
		color.FgCyan + `(       \` + "\n" +
		color.FgCyan + `\    0   \` + "\n" +
		color.FgCyan + ` \        )` + "\n" +
		color.FgCyan + ` /      _/` + "\n" +
		color.FgCyan + `(     _-` + "\n" +
		color.FgCyan + `\____-` + "\n"

	asciiLogoManjaro = color.FgGreen + "||||||||| ||||\n" +
		color.FgGreen + "||||||||| ||||\n" +
		color.FgGreen + "||||      ||||\n" +
		color.FgGreen + "|||| |||| ||||\n" +
		color.FgGreen + "|||| |||| ||||\n" +
		color.FgGreen + "|||| |||| ||||\n" +
		color.FgGreen + "|||| |||| ||||\n"

	asciiLogoMint = color.FgGreen + ` __________` + "\n" +
		color.FgGreen + `|_          \` + "\n" +
		color.FgGreen + `  | | _____ |` + "\n" +
		color.FgGreen + `  | | | | | |` + "\n" +
		color.FgGreen + `  | | | | | |` + "\n" +
		color.FgGreen + `  | \_____/ |` + "\n" +
		color.FgGreen + `  \_________/` + "\n"

	asciiLogoNixos = color.FgBlue + `  \\  \\ //` + "\n" +
		color.FgBlue + ` ==\\__\\/ //` + "\n" +
		color.FgBlue + `   //   \\//` + "\n" +
		color.FgBlue + `==//     //==` + "\n" +
		color.FgBlue + ` //\\___//` + "\n" +
		color.FgBlue + `// /\\  \\==` + "\n" +
		color.FgBlue + `  // \\  \\` + "\n"

	asciiLogoOpenbsd = color.FgYellow + `      _____` + "\n" +
		color.FgYellow + `    \\-     -/` + "\n" +
		color.FgYellow + ` \\_/         \\` + "\n" +
		color.FgYellow + ` |        O O |` + "\n" +
		color.FgYellow + ` |_  <   )  3 )` + "\n" +
		color.FgYellow + ` / \\         /` + "\n" +
		color.FgYellow + `    /-_____-\\` + "\n"

	asciiLogoSuse = color.FgGreen + `  _______` + "\n" +
		color.FgGreen + `__|   __ \` + "\n" +
		color.FgGreen + `     / .\ \` + "\n" +
		color.FgGreen + `     \__/ |` + "\n" +
		color.FgGreen + `   _______|` + "\n" +
		color.FgGreen + `   \_______` + "\n" +
		color.FgGreen + `__________/` + "\n"
	asciiLogoOpensuse = asciLogoSuse

	asciiLogoPopos = color.FgCyan + `______` + "\n" +
		color.FgCyan + `\   _ \        __` + "\n" +
		color.FgCyan + ` \ \ \ \      / /` + "\n" +
		color.FgCyan + `  \ \_\ \    / /` + "\n" +
		color.FgCyan + `   \  ___\  /_/` + "\n" +
		color.FgCyan + `    \ \    _` + "\n" +
		color.FgCyan + `   __\_\__(_)_` + "\n" +
		color.FgCyan + "  (___________)\n"

	asciiLogoRaspbian = color.FgGreen + `   .~~.   .~~.` + "\n" +
		color.FgGreen + `  '. \ ' ' / .'` + "\n" +
		color.FgRed + `   .~ .~~~..~.` + "\n" +
		color.FgRed + `  : .~.'~'.~. :` + "\n" +
		color.FgRed + ` ~ (   ) (   ) ~` + "\n" +
		color.FgRed + `( : '~'.~.'~' : )` + "\n" +
		color.FgRed + ` ~ .~ (   ) ~. ~` + "\n" +
		color.FgRed + `  (  : '~' :  )` + "\n" +
		color.FgRed + `   '~ .~~~. ~'` + "\n" +
		color.FgRed + `       '~'` + "\n"

	asciiLogoRedstar = color.FgRed + `                    ..` + "\n" +
		color.FgRed + `                  .oK0l` + "\n" +
		color.FgRed + `                 :0KKKKd.` + "\n" +
		color.FgRed + `               .xKO0KKKKd` + "\n" +
		color.FgRed + `              ,Od' .d0000l` + "\n" +
		color.FgRed + `             .c;.   .'''...           ..'.` + "\n" +
		color.FgRed + `.,:cloddxxxkkkkOOOOkkkkkkkkxxxxxxxxxkkkx:` + "\n" +
		color.FgRed + `;kOOOOOOOkxOkc'...',;;;;,,,'',;;:cllc:,.` + "\n" +
		color.FgRed + ` .okkkkd,.lko  .......',;:cllc:;,,'''''.` + "\n" +
		color.FgRed + `   .cdo. :xd' cd:.  ..';'',,,'',,;;;,'.` + "\n" +
		color.FgRed + `      . .ddl.;doooc'..;oc;'..';::;,'.` + "\n" +
		color.FgRed + `        coo;.oooolllllllcccc:'.  .` + "\n" +
		color.FgRed + `       .ool''lllllccccccc:::::;.` + "\n" +
		color.FgRed + `       ;lll. .':cccc:::::::;;;;'` + "\n" +
		color.FgRed + `       :lcc:'',..';::::;;;;;;;,,.` + "\n" +
		color.FgRed + `       :cccc::::;...';;;;;,,,,,,.` + "\n" +
		color.FgRed + `       ,::::::;;;,'.  ..',,,,'''.` + "\n" +
		color.FgRed + `        ........          ......` + "\n"

	asciiLogoSteamdeck = color.FgWhite + ` __` + "\n" +
		color.FgWhite + `   \` + "\n" +
		color.FgBlue + `##  ` + color.FgWhite + `\` + "\n" +
		color.FgBlue + `##  ` + color.FgWhite + `/` + "\n" +
		color.FgWhite + ` __/` + "\n"

	asciiLogoUbuntu = color.FgRed + "              ..-::::::-.`\n" +
		color.FgRed + "         `.:+++++++++++" + color.FgWhite + "ooo" + color.FgRed + "++:.`\n" +
		color.FgRed + "       ./+++++++++++++" + color.FgWhite + "sMMMNdyo" + color.FgRed + "+/.\n" +
		color.FgRed + "     .++++++++++++++++" + color.FgWhite + "oyhmMMMMms" + color.FgRed + "++.\n" +
		color.FgRed + "   `/+++++++++" + color.FgWhite + "osyhddddhys" + color.FgRed + "+" + color.FgWhite + "osdMMMh" + color.FgRed + "++/`\n" +
		color.FgRed + "  `+++++++++" + color.FgWhite + "ydMMMMNNNMMMMNds" + color.FgRed + "+" + color.FgWhite + "oyyo" + color.FgRed + "++++`\n" +
		color.FgRed + "  +++++++++" + color.FgWhite + "dMMNhso" + color.FgRed + "++++" + color.FgWhite + "oydNMMmo" + color.FgRed + "++++++++`\n" +
		color.FgRed + " :+" + color.FgWhite + "odmy" + color.FgRed + "+++" + color.FgWhite + "ooysoohmNMMNmyoohMMNs" + color.FgRed + "+++++++:\n" +
		color.FgRed + " ++" + color.FgWhite + "dMMm" + color.FgRed + "+" + color.FgWhite + "oNMd" + color.FgRed + "++" + color.FgWhite + "yMMMmhhmMMNs+yMMNo" + color.FgRed + "+++++++\n" +
		color.FgRed + "`++" + color.FgWhite + "NMMy" + color.FgRed + "+" + color.FgWhite + "hMMd" + color.FgRed + "+" + color.FgWhite + "oMMMs" + color.FgRed + "++++" + color.FgWhite + "sMMN" + color.FgRed + "++" + color.FgWhite + "NMMs" + color.FgRed + "+++++++.\n" +
		color.FgRed + "`++" + color.FgWhite + "NMMy" + color.FgRed + "+" + color.FgWhite + "hMMd" + color.FgRed + "+" + color.FgWhite + "oMMMo" + color.FgRed + "++++" + color.FgWhite + "sMMN" + color.FgRed + "++" + color.FgWhite + "mMMs" + color.FgRed + "+++++++.\n" +
		color.FgRed + " ++" + color.FgWhite + "dMMd" + color.FgRed + "+" + color.FgWhite + "oNMm" + color.FgRed + "++" + color.FgWhite + "yMMNdhhdMMMs" + color.FgRed + "+y" + color.FgWhite + "MMNo" + color.FgRed + "+++++++\n" +
		color.FgRed + " :+" + color.FgWhite + "odmy" + color.FgRed + "++" + color.FgWhite + "oo" + color.FgRed + "+" + color.FgWhite + "ss" + color.FgRed + "+" + color.FgWhite + "ohNMMMMmho" + color.FgRed + "+" + color.FgWhite + "yMMMs" + color.FgRed + "+++++++:\n" +
		color.FgRed + "  +++++++++" + color.FgWhite + "hMMmhs+ooo+oshNMMms" + color.FgRed + "++++++++\n" +
		color.FgRed + "  `++++++++" + color.FgWhite + "oymMMMMNmmNMMMMmy+oys" + color.FgRed + "+++++`\n" +
		color.FgRed + "   `/+++++++++" + color.FgWhite + "oyhdmmmmdhso+sdMMMs" + color.FgRed + "++/\n" +
		color.FgRed + "     ./+++++++++++++++" + color.FgWhite + "oyhdNMMMms" + color.FgRed + "++.\n" +
		color.FgRed + "       ./+++++++++++++" + color.FgWhite + "hMMMNdyo" + color.FgRed + "+/.\n" +
		color.FgRed + "         `.:+++++++++++" + color.FgWhite + "sso" + color.FgRed + "++:.`\n" +
		color.FgRed + "              ..-:::::-..\n"
)

func init() {
	if shell, ok := os.LookupEnv("SHELL"); ok {
		DefaultShell = filepath.Base(shell)
	}

	file, err := os.Open("/etc/os-release")
	if err == nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			Name = strings.Trim(line[12:], `"`)
		} else if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(line[3:], `"`)
		} else if strings.HasPrefix(line, "ID_LIKE=") {
			idLike = strings.Trim(line[8:], `"`)
		}
	}

	setAsciiArt(idLike)
	setAsciiArt(id)
}

func setAsciiArt(id string) {
	switch id {
	case "alpine":
		Ascii = asciiLogoAlpine
	case "arch":
		Ascii = asciiLogoArch
	case "bsd":
		Ascii = asciiLogoBsd
	case "centos":
		Ascii = asciiLogoCentos
	case "debian":
		Ascii = asciiLogoDebian
	case "elementary":
		Ascii = asciiLogoElementary
	case "fedora":
		Ascii = asciiLogoFedora
	case "freebsd":
		Ascii = asciiLogoFreebsd
	case "gentoo":
		Ascii = asciiLogoGentoo
	case "manjaro":
		Ascii = asciiLogoManjaro
	case "mint":
		Ascii = asciiLogoMint
	case "nixos":
		Ascii = asciiLogoNixos
	case "openbsd":
		Ascii = asciiLogoOpenbsd
	case "opensuse":
		Ascii = asciiLogoOpensuse
	case "popos":
		Ascii = asciiLogoPopos
	case "raspbian":
		Ascii = asciiLogoRaspbian
	case "redstar":
		Ascii = asciiLogoRedstar
	case "steamdeck":
		Ascii = asciiLogoSteamdeck
	case "suse":
		Ascii = asciiLogoSuse
	case "ubuntu":
		Ascii = asciiLogoUbuntu
	}
}
