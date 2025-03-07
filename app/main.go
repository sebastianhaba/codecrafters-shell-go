package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		var cmd string
		var arg string

		numArgs, _ := fmt.Fscanln(os.Stdin, &cmd, &arg)
		if cmd == "exit" {
			if numArgs == 2 {
				exitCode, _ := strconv.Atoi(arg)
				os.Exit(exitCode)
			}
		}

		fmt.Fprintf(os.Stdout, "%s: command not found\n", strings.TrimSpace(cmd))
	}
}
