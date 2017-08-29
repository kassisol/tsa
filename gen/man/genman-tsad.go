package main

import (
	"fmt"

	"github.com/juliengk/go-utils/filedir"
	"github.com/kassisol/tsa/daemon"
	"github.com/spf13/cobra/doc"
)

func main() {
	manPath := "/tmp/tsad/man"
	man8 := fmt.Sprintf("%s/man8", manPath)

	if err := filedir.CreateDirIfNotExist(man8, true, 0755); err != nil {
		fmt.Println(err)
	}

	header := &doc.GenManHeader{
		Title:   "TSAD",
		Section: "8",
		Source:  "Harbormaster",
	}
	opts := doc.GenManTreeOptions{
		Header:           header,
		Path:             man8,
		CommandSeparator: "-",
	}

	cmd := daemon.NewCommand()
	cmd.DisableAutoGenTag = true

	if err := doc.GenManTreeFromOpts(cmd, opts); err != nil {
		fmt.Println(err)
	}
}
