package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var supportedCommands = []string{"exit", "echo", "type"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		cmdWithArgs := strings.Split(input, " ")

		if len(cmdWithArgs) == 0 {
			return
		}

		cmd := cmdWithArgs[0]
		args := cmdWithArgs[1:]

		if cmd == "exit" {
			cmdExit(args)
		} else if cmd == "echo" {
			cmdEcho(args)
		} else if cmd == "type" {
			cmdType(args)
		} else if cmdInPath, cmdPath := isCmdInPath(cmd); cmdInPath {
			execCmdInPath(cmd, cmdPath, args)
		} else {
			cmdNotFound(cmdWithArgs[0])
		}
	}
}

func cmdNotFound(cmd string) {
	fmt.Fprintf(os.Stdout, "%s: command not found\n", strings.TrimSpace(cmd))
}

func cmdExit(args []string) {
	exitCode, _ := strconv.Atoi(args[0])
	os.Exit(exitCode)
}

func cmdEcho(args []string) {
	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args[:], " "))
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
