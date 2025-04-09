// Package shellcmd provides an implementation of basic shell commands and a mechanism
// for executing commands in a shell-like environment.
// The package includes built-in commands such as cd, pwd, echo, exit,
// and allows execution of external programs available in the system path.
package shellcmd

// HomeDir represents the path to the user's home directory denoted by the "~" symbol.
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

// ShellCmd represents a single shell command that can be executed.
// It contains the command name, execution function, and path to the program (for external commands).
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

// New creates a new ShellCmd instance based on the provided name.
// The function checks if the command is one of the built-in shell commands.
// If not, it checks if the command exists in the system PATH.
// If the command is not found, it will use the "command not found" handler function.
//
// Example usage:
//
//	cmd := shellcmd.New("echo")
//	output := cmd.Run(args)
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

// Run executes the shell command with the provided arguments.
// Returns the result of the command execution as a string.
// If the execution function is not defined, returns an empty string.
//
// Parameters:
//   - cmd: A CmdWithArgs structure containing the command name, arguments,
//     and output redirection options.
//
// Returns:
//   - The result of the command execution as a string.
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
