package main

import (
	"bufio"
	"fmt"
	"os"
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
	} else {
		fmt.Fprintf(os.Stdout, "%s: not found\n", cmdToCheck)
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
