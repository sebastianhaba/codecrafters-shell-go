package shellcmd

import (
	"fmt"
	"os"
	"strings"
)

func cmdEcho(p runFunParams) string {
	if p.Cmd.RedirectOptions.StdOutputFile != "" {
		var file *os.File
		var err error
		if p.Cmd.RedirectOptions.StdOutputAppend {
			file, err = os.OpenFile(p.Cmd.RedirectOptions.StdOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return ""
			}

			defer file.Close()
			_, err = file.WriteString("\n" + strings.Join(p.Cmd.Args, " "))

		} else {
			file, err = os.Create(p.Cmd.RedirectOptions.StdOutputFile)
			if err != nil {
				return ""
			}

			defer file.Close()
			_, err = file.WriteString(strings.Join(p.Cmd.Args, " "))
		}

		if err != nil {
			return ""
		}

		return ""
	}

	if p.Cmd.RedirectOptions.StdErrorFile != "" {
		file, err := os.Create(p.Cmd.RedirectOptions.StdErrorFile)
		if err != nil {
			return ""
		}

		defer file.Close()

		output := strings.Join(p.Cmd.Args, " ")
		return output
	}

	return fmt.Sprintf("%s", strings.Join(p.Cmd.Args, " "))
}
