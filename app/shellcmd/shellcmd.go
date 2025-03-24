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
	Args []string
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

func (s *ShellCmd) Run(args []string) string {
	if s.run == nil {
		return ""
	}

	return s.run(runFunParams{
		Name: s.Name,
		Args: args,
		Path: s.path,
	})
}

func cmdExit(p runFunParams) string {
	var exitCode int
	var err error
	if len(p.Args) > 0 {
		exitCode, err = strconv.Atoi(p.Args[0])
		if err != nil {
			exitCode = 0
		}
	}

	os.Exit(exitCode)
	return ""
}

func cmdEcho(p runFunParams) string {
	return fmt.Sprintf("%s", strings.Join(p.Args, " "))
}

func cmdPwd(p runFunParams) string {
	currentDir, _ := os.Getwd()
	return fmt.Sprintf("%s", currentDir)
}

func cmdCd(p runFunParams) string {
	var path string
	if len(p.Args) == 0 {
		path = ""
	}

	path = p.Args[0]
	if path == HomeDir {
		path = getHomeDir()
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Sprintf("%s: %s: No such file or directory", p.Name, path)
	}

	return ""
}

func cmdType(p runFunParams) string {
	cmdToCheck := p.Args[0]
	_, cmdIsSupported := supportedCmd[cmdToCheck]
	if cmdIsSupported {
		return fmt.Sprintf("%s is a shell builtin", cmdToCheck)
	}

	externalCmdResult := checkCmdInPath(cmdToCheck)
	return fmt.Sprintf("%s", externalCmdResult)
}

func cmdExternal(p runFunParams) string {
	var stdout bytes.Buffer
	fileExt := ""

	if runtime.GOOS == "windows" {
		fileExt = ".exe"
	}

	cmd := exec.Command(p.Name+fileExt, p.Args...)
	cmd.Dir = strings.Replace(p.Path, p.Name+fileExt, "", 1)
	cmd.Stdout = &stdout
	cmd.Run()
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
