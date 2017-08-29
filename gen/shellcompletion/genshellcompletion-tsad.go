package main

import (
	"fmt"

	"github.com/juliengk/go-utils/filedir"
	"github.com/kassisol/tsa/daemon"
)

func main() {
	scPath := "/tmp/tsad/shellcompletion"
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := filedir.CreateDirIfNotExist(scPath, true, 0755); err != nil {
		fmt.Println(err)
	}

	cmd := daemon.NewCommand()
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
	}
}
