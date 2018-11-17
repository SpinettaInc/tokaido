package utils

import (
	"os/exec"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
)

// CommandSubstitution - Execute a command and return the output value. No exit on stdErr
func CommandSubstitution(name string, args ...string) string {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = conf.GetConfig().Tokaido.Project.Path
	stdoutStderr, _ := cmd.CombinedOutput()

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}

// CommandSubSplitOutput - Execute a command and return the output value split into stdout and stderr. No exit on stdErr
func CommandSubSplitOutput(name string, args ...string) (string, error) {
	return CommandSubSplitOutputContext(conf.GetConfig().Tokaido.Project.Path, name, args...)
}

// CommandSubSplitOutputContext - Execute a command and return the output value split into stdout and stderr from a directory context. No exit on stdErr
func CommandSubSplitOutputContext(directory, name string, args ...string) (string, error) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		DebugErrOutput(err)
		return "", err
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr)), nil
}
