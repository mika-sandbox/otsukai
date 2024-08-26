package session

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"otsukai/runtime/value"
	"strings"
	"time"
)

type RemoteSession struct {
	client  *ssh.Client
	session *ssh.Session
	stdout  bytes.Buffer
}

type CreateRemoteSessionOpts struct {
	Remote  value.IValueObject
	Timeout *time.Duration
}

func getUserName(remote value.IValueObject) (*string, error) {
	hash, err := remote.ToHashObject()
	if err != nil {
		return nil, err
	}

	user, err := hash["user"].ToString()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getRemoteHost(remote value.IValueObject) (*string, error) {
	hash, err := remote.ToHashObject()
	if err != nil {
		return nil, err
	}

	host, err := hash["host"].ToString()
	if err != nil {
		return nil, err
	}

	if strings.Contains(*host, ":") {
		return host, nil
	}

	hostWithPort := *host + ":22"

	return &hostWithPort, nil
}

func CreateRemoteSession(opts *CreateRemoteSessionOpts) (*RemoteSession, error) {
	user, err := getUserName(opts.Remote)
	if err != nil {
		return nil, err
	}

	host, err := getRemoteHost(opts.Remote)
	if err != nil {
		return nil, err
	}

	timeout := 10 * time.Second
	if opts.Timeout != nil {
		timeout = *opts.Timeout
	}

	cfg := &ssh.ClientConfig{
		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	client, err := ssh.Dial("tcp", *host, cfg)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	var stdout bytes.Buffer
	session.Stdout = &stdout

	return &RemoteSession{
		client,
		session,
		stdout,
	}, nil
}

func (session *RemoteSession) Run(command string, stdout bool) error {
	ns, err := session.client.NewSession()
	if err != nil {
		return err
	}

	ns.Stdout = &session.stdout
	if err = ns.Run(command); err != nil {
		return err
	}

	if stdout {
		fmt.Println(session.stdout.String())
	}

	return nil
}

func (session *RemoteSession) Close() {
	if session.session != nil {
		session.session.Close()
	}

	if session.client != nil {
		session.client.Close()
	}
}
