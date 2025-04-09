package shellcmd

import (
	"fmt"
	"os"
)

func cmdCd(p runFunParams) string {
	var path string
	if len(p.Cmd.Args) == 0 {
		path = ""
	}

	path = p.Cmd.Args[0]
	if path == HomeDir {
		path = getHomeDir()
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Sprintf("%s: %s: No such file or directory", p.Name, path)
	}

	return ""
}
