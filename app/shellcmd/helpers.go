package shellcmd

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"runtime"
	"strings"
)

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
