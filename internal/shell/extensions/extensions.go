package extensions

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var extensionRegister []*extensionWrapper

type (
	Extension interface {
		PromptDisplay() string
		FileMatch(f os.DirEntry, p string) bool
	}
	extensionWrapper struct {
		extension Extension
		active    bool
	}
)

func register(e Extension) {
	extensionRegister = append(extensionRegister, &extensionWrapper{extension: e, active: true})
}

func priorityRegister(e Extension) {
	extensionRegister = append([]*extensionWrapper{{extension: e, active: true}}, extensionRegister...)
}

func SetContext(root string) error {
	// deactivate all extensions
	for _, ew := range extensionRegister {
		ew.active = false
	}

	// activate extensions that match the current context
	for {
		entries, err := os.ReadDir(root)
		if err != nil {
			return errors.Join(errors.New("failed to read directory"), err)
		}
		for _, f := range entries {
			p := filepath.Join(root, f.Name())
			for _, ew := range extensionRegister {
				if ew.active {
					continue
				}
				if ew.extension.FileMatch(f, p) {
					ew.active = true
					continue
				}
			}
		}

		nextRoot := filepath.Dir(root)
		if nextRoot == root {
			break
		}
		root = nextRoot
	}
	return nil
}

func PromptDisplay() string {
	sb := strings.Builder{}
	for _, ew := range extensionRegister {
		if ew.active {
			sb.WriteString(ew.extension.PromptDisplay())
			sb.WriteString(" ")
		}
	}
	return sb.String()
}
