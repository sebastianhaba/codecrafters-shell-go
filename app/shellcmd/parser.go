package shellcmd

import "strings"

type RedirectOptions struct {
	StdOutputFile   string
	StdErrorFile    string
	StdOutputAppend bool
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

	redirectStdOutputIndex, appendStdOutput := redirectStdOutputIndex(args)
	redirectStdErrorIndex := redirectStdErrorIndex(args)

	if redirectStdOutputIndex == -1 && redirectStdErrorIndex == -1 {
		cmd.Args = args
		return
	}

	argIndex := redirectStdOutputIndex
	if redirectStdOutputIndex == -1 {
		argIndex = redirectStdErrorIndex
	} else if redirectStdErrorIndex == -1 {
		argIndex = redirectStdOutputIndex
	} else {
		if redirectStdOutputIndex < redirectStdErrorIndex {
			argIndex = redirectStdOutputIndex
		} else {
			argIndex = redirectStdErrorIndex
		}
	}

	cmd.Args = args[:argIndex]

	if redirectStdOutputIndex != -1 {
		cmd.RedirectOptions.StdOutputFile = args[redirectStdOutputIndex+1]
		cmd.RedirectOptions.StdOutputAppend = appendStdOutput
	}

	if redirectStdErrorIndex != -1 {
		cmd.RedirectOptions.StdErrorFile = args[redirectStdErrorIndex+1]
	}
}

func redirectStdOutputIndex(args []string) (int, bool) {
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			return i, false
		}

		if arg == ">>" || arg == "1>>" {
			return i, true
		}
	}

	return -1, false
}

func redirectStdErrorIndex(args []string) int {
	for i, arg := range args {
		if arg == "2>" {
			return i
		}
	}

	return -1
}
