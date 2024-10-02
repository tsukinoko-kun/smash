package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// find all protobuf files
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".proto" {
			fmt.Fprintf(os.Stderr, "compiling %s\n", path)
			// compile the protobuf file
			if err := run("protoc", "--go_out=.", path); err != nil {
				return err
			}
		}
		return nil
	})
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
