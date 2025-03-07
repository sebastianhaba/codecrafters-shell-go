package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
			exit(args)
		} else if cmd == "echo" {
			echo(args)
		} else {
			commandNotFound(cmdWithArgs[0])
		}
	}
}

func commandNotFound(cmd string) {
	fmt.Fprintf(os.Stdout, "%s: command not found\n", strings.TrimSpace(cmd))
}

func exit(args []string) {
	exitCode, _ := strconv.Atoi(args[0]e)
	os.Exit(exitCode)
}

func echo(args []string) {
	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args[:], " "))
}
