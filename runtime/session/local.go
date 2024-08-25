package session

import "os/exec"

type LocalSession struct{}

func CreateLocalSession() (*LocalSession, error) {
	return &LocalSession{}, nil
}

func (session *LocalSession) Run(command string, stdout bool) error {
	cmd := exec.Command(command)
	return cmd.Run()
}

func (session *LocalSession) Close() {
	// NOTHING TO DO
}
