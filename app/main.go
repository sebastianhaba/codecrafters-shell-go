package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/shellcmd"
	"os"
	"os/exec"
	"strings"
)

var supportedCommands = []string{"exit", "echo", "type", "pwd", "cd"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		cmdWithArgs := parseArgs(input)
		cmdName := cmdWithArgs[0]
		args := cmdWithArgs[1:]

		if len(cmdWithArgs) == 0 {
			return
		}

		cmd := shellcmd.New(cmdName)
		result := cmd.Run(args)

		if result != "" {
			fmt.Fprintf(os.Stdout, "%s\n", result)
		}
	}
}

func parseArgs(input string) []string {
	var result []string
	var currentArg strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch c {
		case '\'':
			inQuotes = !inQuotes

			if !inQuotes && i+1 < len(input) && input[i+1] == '\'' {
				inQuotes = true
				i++ // Pomijamy następny apostrof
			}
		case ' ', '\t':
			if !inQuotes {
				if currentArg.Len() > 0 {
					result = append(result, currentArg.String())
					currentArg.Reset()
				}
			} else {
				currentArg.WriteByte(c)
			}
		default:
			currentArg.WriteByte(c)
		}
	}

	if currentArg.Len() > 0 {
		result = append(result, currentArg.String())
	}

	return result
}

func cmdType(args []string) {
	cmdToCheck := args[0]
	supportedCmd := contains(supportedCommands, cmdToCheck)
	if supportedCmd {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", cmdToCheck)
		return
	}

	pathResult := checkCmdInPath(cmdToCheck)
	fmt.Fprintf(os.Stdout, "%s", pathResult)
}

func cmdCd(args []string) {
	path := args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", path)
	}
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}

	return false
}

func isCmdInPath(cmd string) (bool, string) {
	pathEnv := os.Getenv("PATH")
	directoriesToSearch := strings.Split(pathEnv, ":")
	for _, directory := range directoriesToSearch {
		path := directory + "/" + cmd
		if _, err := os.Stat(path); err == nil {
			return true, path
		}
	}

	return false, ""
}

func checkCmdInPath(cmd string) string {
	exists, path := isCmdInPath(cmd)
	if exists {
		return fmt.Sprintf("%s is %s\n", cmd, path)
	}

	return fmt.Sprintf("%s: not found\n", cmd)
}

func execCmdInPath(cmdName, cmdPath string, args []string) {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = strings.Replace(cmdPath, cmdName, "", 1)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
