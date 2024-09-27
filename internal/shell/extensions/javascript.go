package extensions

import (
	"encoding/json"
	"os"
	"smash/internal/color"
	"strings"
)

type (
	javascript struct {
		packageJson string
	}
	packageJson struct {
		Scripts map[string]string `json:"scripts"`
	}
)

const javascriptPromptDisplay = color.FgYellow + "\uE781"

func (_ *javascript) PromptDisplay() string {
	return javascriptPromptDisplay
}

func (js *javascript) FileMatch(f os.DirEntry, p string) bool {
	x := f.Name() == "package.json"
	if x {
		js.packageJson = p
	}
	return x
}

var (
	npm0 = []string{"install", "run"}
)

func (js *javascript) Completions(name string, args []string) []string {
	if name == "npm" || name == "pnpm" {
		switch len(args) {
		case 1:
			var cs []string
			for _, s := range npm0 {
				if strings.HasPrefix(s, args[0]) {
					cs = append(cs, s)
				}
			}
			return cs
		case 2:
			switch args[0] {
			case "run":
				pj, ok := js.parsePackageJson()
				if !ok {
					return nil
				}
				var scripts []string
				for s := range pj.Scripts {
					if strings.HasPrefix(s, args[1]) {
						scripts = append(scripts, s)
					}
				}
				return scripts
			}
		}
	}
	return nil
}

func (js *javascript) parsePackageJson() (packageJson, bool) {
	f, err := os.Open(js.packageJson)
	if err != nil {
		return packageJson{}, false
	}
	defer f.Close()

	var pj packageJson
	err = json.NewDecoder(f).Decode(&pj)
	return pj, err == nil
}

func (_ *javascript) Title() string {
	return "JavaScript"
}

func init() {
	register(&javascript{})
}
