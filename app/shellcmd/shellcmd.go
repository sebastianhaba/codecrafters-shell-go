package shellcmd

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"sort"
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

var autoCompleteTabs = ""

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

	/*
		supportedCmd["xyz_foo"] = cmdExit
		supportedCmd["xyz_foo_bar"] = cmdExit
		supportedCmd["xyz_foo_bar_baz"] = cmdExit

	*/

}

type ShellAutoCompleter struct{}

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

	// Brak dopasowania
	return [][]rune{[]rune("\a")}, 0
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

func AutoComplete() ShellAutoCompleter {
	return ShellAutoCompleter{}
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
		/*
			_, err = file.WriteString(output)
			if err != nil {
				return ""
			}
		*/

		return output
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

func getCmdsInPathEnv() []string {
	pathEnv := os.Getenv("PATH")
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	directoriesToSearch := strings.Split(pathEnv, separator)
	var result []string
	for _, directory := range directoriesToSearch {
		files, _ := os.ReadDir(directory)
		for _, file := range files {
			result = append(result, file.Name())
		}
	}

	return result
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sameSize(s []string) bool {
	maxLen := len(s[0])
	for _, a := range s {
		if len(a) != maxLen {
			return false
		}
	}

	return true
}

func getBestMatch(s []string) string {
	maxLen := math.MaxInt32
	bestMatch := ""
	for _, a := range s {
		if len(a) < maxLen {
			bestMatch = a
			maxLen = len(a)
		}
	}

	return bestMatch
}
