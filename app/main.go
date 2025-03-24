package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/shellcmd"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		cmdWithArgs := shellcmd.ParseArgs(input)
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
