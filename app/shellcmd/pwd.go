package shellcmd

import (
	"fmt"
	"os"
)

func cmdPwd(p runFunParams) string {
	currentDir, _ := os.Getwd()
	return fmt.Sprintf("%s", currentDir)
}
