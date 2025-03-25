package shellcmd

import "strings"

func ParseArgs(input string) []string {
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

	return result
}
