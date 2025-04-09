package shellcmd

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func cmdExternal(p runFunParams) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fileExt := ""

	if runtime.GOOS == "windows" {
		fileExt = ".exe"
	}

	cmd := exec.Command(p.Name+fileExt, p.Cmd.Args...)
	cmd.Dir = strings.Replace(p.Path, p.Name+fileExt, "", 1)

	if p.Cmd.RedirectOptions.StdOutputFile != "" {
		var file *os.File
		var err error
		if p.Cmd.RedirectOptions.StdOutputAppend {
			file, err = os.OpenFile(p.Cmd.RedirectOptions.StdOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(p.Cmd.RedirectOptions.StdOutputFile)
		}
		if err != nil {
			return ""
		}
		defer file.Close()

		if p.Cmd.RedirectOptions.StdOutputAppend {
			stat, _ := os.Stat(p.Cmd.RedirectOptions.StdOutputFile)
			if stat.Size() > 0 {
				file.WriteString("\n")
			}
		}

		cmd.Stdout = file
	} else {
		cmd.Stdout = &stdout
	}

	if p.Cmd.RedirectOptions.StdErrorFile != "" {
		var file *os.File
		var err error
		if p.Cmd.RedirectOptions.StdErrorAppend {
			file, err = os.OpenFile(p.Cmd.RedirectOptions.StdErrorFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(p.Cmd.RedirectOptions.StdErrorFile)
		}
		if err != nil {
			return ""
		}
		defer file.Close()
		cmd.Stderr = file
	} else {
		cmd.Stderr = &stderr
	}

	cmd.Run()

	if p.Cmd.RedirectOptions.StdOutputFile != "" {
		return strings.TrimSpace(stderr.String())
	}

	return strings.TrimSpace(stdout.String())
}
