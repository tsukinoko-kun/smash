package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	Config      *smashConfig
	Alias       map[string][]string
	HistoryFile string
)

func init() {
	loadVars()
	HistoryFile = filepath.Join(getConfigDir(), "history.txt")
	Config = getConfigFile()
	if Config.Alias != nil {
		Alias = make(map[string][]string, len(Config.Alias))
		for key, val := range Config.Alias {
			switch v := val.(type) {
			case string:
				Alias[key] = []string{v}
			case []any:
				sArr := make([]string, len(v))
				for i, s := range v {
					if s, ok := s.(string); ok {
						sArr[i] = s
					} else {
						sArr[i] = fmt.Sprintf("%v", s)
					}
				}
				Alias[key] = sArr
			default:
				Alias[key] = []string{fmt.Sprintf("%v", val)}
			}
		}
	} else {
		Alias = make(map[string][]string)
	}
}

type (
	smashColor struct {
		CompletionText       string `toml:"completion_text"`
		CompletionSelectedBg string `toml:"completion_selected_bg"`
	}
	smashConfig struct {
		InteractivePrompt string         `toml:"ps1"`
		LogPrompt         string         `toml:"ps2"`
		Alias             map[string]any `toml:"alias"`
		Color             smashColor     `toml:"color"`
	}
)

func getConfigDir() string {
	var d string
	if xdgConfig, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		d = filepath.Join(xdgConfig, "smash")
	} else {
		xdgConfig = filepath.Join(GetUser().HomeDir, ".config")
		if _, err := os.Stat(xdgConfig); err == nil {
			d = filepath.Join(xdgConfig, "smash")
		} else {
			userConfig, err := os.UserConfigDir()
			if err != nil {
				panic(errors.Join(errors.New("failed to get user config dir"), err))
			}
			d = filepath.Join(userConfig, "smash")
		}
	}
	if err := os.MkdirAll(d, 0700); err != nil {
		panic(errors.Join(errors.New("failed to create config dir"), err))
	}
	return d
}

func getConfigFile() *smashConfig {
	p := filepath.Join(getConfigDir(), "config.toml")
	c := &smashConfig{
		Alias: make(map[string]any),
		Color: smashColor{
			CompletionText:       "8",
			CompletionSelectedBg: "4",
		},
	}
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			// create file
			func() {
				fmt.Printf("creating config file %s\n", p)
				f, err := os.Create(p)
				if err != nil {
					panic(errors.Join(errors.New("failed to create config file"), err))
				}
				defer f.Close()
				c.InteractivePrompt = "${Color.FgHiBlack}$USER@$PWD\t$DEV${Color.Reset}\n${Color.FgBlue}❯${Color.Reset} "
				c.LogPrompt = "${Color.FgHiBlack}$PWD${Color.Reset} "
				c.Alias["l"] = []string{"ls", "-l"}
				e := toml.NewEncoder(f)
				if err := e.Encode(c); err != nil {
					panic(errors.Join(errors.New("failed to write config file"), err))
				}
			}()
		} else {
			panic(errors.Join(errors.New("failed to stat config file"), err))
		}
	}

	f, err := os.Open(p)
	if err != nil {
		panic(errors.Join(errors.New("failed to open config file"), err))
	}
	defer f.Close()

	d := toml.NewDecoder(f)
	if _, err := d.Decode(c); err != nil {
		panic(errors.Join(errors.New("failed to parse config file"), err))
	}
	return c
}
