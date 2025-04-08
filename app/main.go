package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/shellcmd"
	"log"
	"strings"
)

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    shellcmd.AutoComplete(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil { // np. Ctrl+D
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		cmdWithArgs := shellcmd.ParseArgs(input)

		cmd := shellcmd.New(cmdWithArgs.Name)
		result := cmd.Run(cmdWithArgs)

		if result != "" {
			fmt.Println(result)
		}
	}
}
