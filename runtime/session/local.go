package session

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"os/exec"
	"otsukai"
	re "otsukai/runtime/errors"
)

type LocalSession struct{}

func CreateLocalSession() (*LocalSession, error) {
	return &LocalSession{}, nil
}

func (session *LocalSession) Run(command string, stdout bool) error {
	args, err := shellwords.Parse(command)
	if err != nil {
		otsukai.Errf("failed to parse provided command: %s", err)
		return re.EXECUTION_ERROR
	}

	var cmd *exec.Cmd
	switch len(args) {
	case 0:
		return nil

	case 1:
		cmd = exec.Command(args[0])
		break

	default:
		cmd = exec.Command(args[0], args[1:]...)
	}

	out, err := cmd.Output()
	if err != nil {
		otsukai.Errf("failed to execute command: %s", err)
		return re.EXECUTION_ERROR
	}

	if stdout {
		fmt.Println(string(out))
	}

	return nil
}

func (session *LocalSession) CopyToRemote(local string, remote string) error {
	return nil
}

func (session *LocalSession) CopyToLocal(local string, remote string, isDir bool) error {
	return nil
}

func (session *LocalSession) Close() {
	// NOTHING TO DO
}
