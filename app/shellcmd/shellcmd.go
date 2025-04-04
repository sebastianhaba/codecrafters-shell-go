package shellcmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

const (
	HomeDir = "~"
)

type runFunParams struct {
	Name string
	Path string
	Cmd  *CmdWithArgs
}

type runFunc func(p runFunParams) string

var supportedCmd = map[string]runFunc{}

type ShellCmd struct {
	Name string
	run  runFunc
	path string
}

func init() {
	supportedCmd["exit"] = cmdExit
	supportedCmd["echo"] = cmdEcho
	supportedCmd["type"] = cmdType
	supportedCmd["pwd"] = cmdPwd
	supportedCmd["cd"] = cmdCd
	supportedCmd["type"] = cmdType
}

func New(name string) *ShellCmd {
	cmdFunc, exists := supportedCmd[name]
	var extCmdPath string
	var extCmdExists bool
	if !exists {
		if extCmdExists, extCmdPath = isCmdInPath(name); extCmdExists {
			cmdFunc = cmdExternal
		} else {
			cmdFunc = cmdNotFound
		}
	}

	return &ShellCmd{
		Name: name,
		run:  cmdFunc,
		path: extCmdPath,
	}
}

func (s *ShellCmd) Run(cmd *CmdWithArgs) string {
	if s.run == nil {
		return ""
	}

	return s.run(runFunParams{
		Name: s.Name,
		Cmd:  cmd,
		Path: s.path,
	})
}

func cmdExit(p runFunParams) string {
	var exitCode int
	var err error
	if len(p.Cmd.Args) > 0 {
		exitCode, err = strconv.Atoi(p.Cmd.Args[0])
		if err != nil {
			exitCode = 0
		}
	}

	os.Exit(exitCode)
	return ""
}

func cmdEcho(p runFunParams) string {
	if p.Cmd.RedirectOptions.File != "" {
		file, err := os.Create(p.Cmd.RedirectOptions.File)
		if err != nil {
			return ""
		}

		defer file.Close()

		_, err = file.WriteString(strings.Join(p.Cmd.Args, " "))
		if err != nil {
			return ""
		}

		return ""
	}

	return fmt.Sprintf("%s", strings.Join(p.Cmd.Args, " "))
}

func cmdPwd(p runFunParams) string {
	currentDir, _ := os.Getwd()
	return fmt.Sprintf("%s", currentDir)
}

func cmdCd(p runFunParams) string {
	var path string
	if len(p.Cmd.Args) == 0 {
		path = ""
	}

	path = p.Cmd.Args[0]
	if path == HomeDir {
		path = getHomeDir()
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Sprintf("%s: %s: No such file or directory", p.Name, path)
	}

	return ""
}

func cmdType(p runFunParams) string {
	cmdToCheck := p.Cmd.Args[0]
	_, cmdIsSupported := supportedCmd[cmdToCheck]
	if cmdIsSupported {
		return fmt.Sprintf("%s is a shell builtin", cmdToCheck)
	}

	externalCmdResult := checkCmdInPath(cmdToCheck)
	return fmt.Sprintf("%s", externalCmdResult)
}

func cmdExternal(p runFunParams) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fileExt := ""

	if runtime.GOOS == "windows" {
		fileExt = ".exe"
	}

	cmd := exec.Command(p.Name+fileExt, p.Cmd.Args...)
	cmd.Dir = strings.Replace(p.Path, p.Name+fileExt, "", 1)

	if p.Cmd.RedirectOptions.File != "" {
		file, err := os.Create(p.Cmd.RedirectOptions.File)
		if err != nil {
			return ""
		}
		defer file.Close()
		cmd.Stdout = file
		cmd.Stderr = &stderr
	} else {
		cmd.Stdout = &stdout
	}

	cmd.Run()

	if p.Cmd.RedirectOptions.File != "" {
		return strings.TrimSpace(stderr.String())
	}

	return strings.TrimSpace(stdout.String())
}

func cmdNotFound(p runFunParams) string {
	return fmt.Sprintf("%s: command not found", p.Name)
}

func getHomeDir() string {
	homeDirectory := ""
	if runtime.GOOS == "windows" {
		usr, err := user.Current()
		if err == nil {
			homeDirectory = usr.HomeDir
		}
	} else {
		homeDirectory = os.Getenv("HOME")
	}

	return homeDirectory
}

func isCmdInPath(cmd string) (bool, string) {
	pathEnv := os.Getenv("PATH")
	separator := ":"
	fileExt := ""
	directorySeparator := "/"
	if runtime.GOOS == "windows" {
		separator = ";"
		fileExt = ".exe"
		directorySeparator = "\\"
	}
	directoriesToSearch := strings.Split(pathEnv, separator)
	for _, directory := range directoriesToSearch {
		path := directory + directorySeparator + cmd + fileExt
		if _, err := os.Stat(path); err == nil {
			return true, path
		}
	}

	return false, ""
}

func checkCmdInPath(cmd string) string {
	exists, path := isCmdInPath(cmd)
	if exists {
		return fmt.Sprintf("%s is %s", cmd, path)
	}

	return fmt.Sprintf("%s: not found", cmd)
}
