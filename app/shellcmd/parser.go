package shellcmd

import "strings"

func ParseArgs(input string) []string {
	var result []string
	var currentArg strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNext := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		if escapeNext {
			if inDoubleQuotes && (c == '\\' || c == '$' || c == '"' || c == '\n') {
				currentArg.WriteByte(c)
			} else if !inDoubleQuotes {
				currentArg.WriteByte(c)
			} else {
				currentArg.WriteByte('\\')
				currentArg.WriteByte(c)
			}
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

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

		currentArg.WriteByte(c)
	}

	if currentArg.Len() > 0 {
		result = append(result, currentArg.String())
	}

	return result
}
