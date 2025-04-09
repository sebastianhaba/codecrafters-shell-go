package shellcmd

import (
	"os"
	"strconv"
)

func cmdExit(p runFunParams) string {
	var exitCode int
	var err error
	if len(p.Cmd.Args) > 0 {
		exitCode, err = strconv.Atoi(p.Cmd.Args[0])
		if err != nil {
			exitCode = 0
		}
	}

	os.Exit(exitCode)
	return ""
}
