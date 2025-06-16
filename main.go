package main

import (
	"github.com/fun7257/sgv/cmd"
)

var (
	goVersion string
	commit    string
)

func main() {
	cmd.SetBuildInfo(goVersion, commit)
	cmd.Execute()
}
