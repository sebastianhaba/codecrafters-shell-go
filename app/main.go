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

		cmd := shellcmd.New(cmdWithArgs.Name)
		result := cmd.Run(cmdWithArgs)

		if result != "" {
			fmt.Fprintf(os.Stdout, "%s\n", result)
		}
	}
}
