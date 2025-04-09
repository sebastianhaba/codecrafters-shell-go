package shellcmd

import "fmt"

func cmdType(p runFunParams) string {
	cmdToCheck := p.Cmd.Args[0]
	_, cmdIsSupported := supportedCmd[cmdToCheck]
	if cmdIsSupported {
		return fmt.Sprintf("%s is a shell builtin", cmdToCheck)
	}

	externalCmdResult := checkCmdInPath(cmdToCheck)
	return fmt.Sprintf("%s", externalCmdResult)
}
