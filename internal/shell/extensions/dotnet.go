package extensions

import (
	"os"
	"smash/internal/color"
	"strings"
)

type dotnet struct {
}

const dotnetPromptDisplay = color.FgMagenta + "\U000F0AAE"

func (_ *dotnet) PromptDisplay() string {
	return dotnetPromptDisplay
}

func (_ *dotnet) FileMatch(f os.DirEntry, _ string) bool {
	n := f.Name()
	return strings.HasSuffix(n, ".sln") || strings.HasSuffix(n, ".csproj") || strings.HasSuffix(n, ".fsproj")
}

func init() {
	register(&dotnet{})
}
