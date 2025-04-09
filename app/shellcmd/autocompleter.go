package shellcmd

import (
	"sort"
	"strings"
)

var autoCompleteTabs = ""

type ShellAutoCompleter struct{}

// Do performs auto-completion based on the entered text.
// The method analyzes the partially entered text and returns possible completions.
//
// Parameters:
//   - line: Array of runes representing the current input line
//   - pos: Current cursor position in the line
//
// Returns:
//   - newLine: Array of possible completions
//   - length: Length of the completion
func (c ShellAutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	commands := make([]string, 0)
	for key := range supportedCmd {
		commands = append(commands, key)
	}

	for _, cmd := range getCmdsInPathEnv() {
		commands = append(commands, cmd)
	}

	input := string(line[:pos])
	input = strings.Replace(input, "\a", "", -1)

	possibleMatches := make([]string, 0)

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) {
			cmd = strings.Replace(cmd, input, "", 1)
			if contains(possibleMatches, cmd) {
				continue
			}

			possibleMatches = append(possibleMatches, cmd)
		}
	}

	if len(possibleMatches) == 1 {
		autoCompleteTabs = ""
		return [][]rune{[]rune(possibleMatches[0] + " ")}, len(possibleMatches[0])
	}

	if len(possibleMatches) > 1 {
		if sameSize(possibleMatches) {
			if autoCompleteTabs == input {
				sort.Strings(possibleMatches)
				autocomplete := ""

				for index, cmd := range possibleMatches {
					if index == 0 {
						autocomplete += input + cmd
					} else {
						autocomplete += "  " + input + cmd
					}
				}

				autocomplete = "\n" + autocomplete + "\n$ " + input

				return [][]rune{[]rune(autocomplete)}, len(autocomplete)
			} else {
				autoCompleteTabs = input
				return [][]rune{[]rune("\a")}, 0
			}
		} else {
			autocomplete := getBestMatch(possibleMatches)
			return [][]rune{[]rune(autocomplete)}, len(autocomplete)
		}
	}

	return [][]rune{[]rune("\a")}, 0
}

// AutoComplete returns a ShellAutoCompleter object that provides command
// auto-completion functionality for the shell.
//
// Example usage:
//
//	completer := shellcmd.AutoComplete()
//
// Use completer for auto-completion of entered commands
func AutoComplete() ShellAutoCompleter {
	return ShellAutoCompleter{}
}
