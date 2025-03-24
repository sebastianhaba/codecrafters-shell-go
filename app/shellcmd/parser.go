package shellcmd

import "strings"

func ParseArgs(input string) []string {
	var result []string
	var currentArg strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch c {
		case '\'':
			inQuotes = !inQuotes

			if !inQuotes && i+1 < len(input) && input[i+1] == '\'' {
				inQuotes = true
				i++ // Pomijamy następny apostrof
			}
		case ' ', '\t':
			if !inQuotes {
				if currentArg.Len() > 0 {
					result = append(result, currentArg.String())
					currentArg.Reset()
				}
			} else {
				currentArg.WriteByte(c)
			}
		default:
			currentArg.WriteByte(c)
		}
	}

	if currentArg.Len() > 0 {
		result = append(result, currentArg.String())
	}

	return result
}
