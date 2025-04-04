package shellcmd

import "strings"

type RedirectOptions struct {
	File string
}

type CmdWithArgs struct {
	Name            string
	Args            []string
	RedirectOptions RedirectOptions
}

func ParseArgs(input string) *CmdWithArgs {
	var result []string
	var currentArg strings.Builder
	var i int
	inSingleQuotes := false
	inDoubleQuotes := false

	for i = 0; i < len(input); i++ {
		c := input[i]

		if c == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			continue
		}

		if c == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if (c == ' ' || c == '\t') && !inSingleQuotes && !inDoubleQuotes {
			if currentArg.Len() > 0 {
				result = append(result, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		if c == '\\' {
			if inSingleQuotes {
				currentArg.WriteByte('\\')
				continue
			}

			if i+1 < len(input) {
				nextChar := input[i+1]
				if inDoubleQuotes {
					if nextChar == '$' || nextChar == '"' || nextChar == '\\' || nextChar == '`' {
						i++
						currentArg.WriteByte(nextChar)
					} else {
						currentArg.WriteByte('\\')
					}
				} else {
					i++
					currentArg.WriteByte(nextChar)
				}
			} else {
				currentArg.WriteByte('\\')
			}
			continue
		}

		currentArg.WriteByte(c)
	}

	if currentArg.Len() > 0 {
		result = append(result, currentArg.String())
	}

	if len(result) == 0 {
		return nil
	}

	cmdWithArgs := &CmdWithArgs{
		Name: result[0],
	}

	prepareArgs(cmdWithArgs, result[1:])

	return cmdWithArgs
}

func prepareArgs(cmd *CmdWithArgs, args []string) {
	if len(args) == 0 {
		return
	}

	if len(args) == 1 {
		cmd.Args = []string{args[0]}
		return
	}

	redirectOptionsIndex := redirectOptionsIndex(args)
	if redirectOptionsIndex == -1 {
		cmd.Args = args
		return
	}

	cmd.Args = args[:redirectOptionsIndex]
	cmd.RedirectOptions.File = args[redirectOptionsIndex+1]
}

func redirectOptionsIndex(args []string) int {
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			return i
		}
	}

	return -1
}
