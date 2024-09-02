package session

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mika-sandbox/otsukai/logger"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/value"
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
	"os"
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
		logger.Errf("ssh error: username not found")
		return nil, re.EXECUTION_ERROR
	}

	host, err := getRemoteHost(opts.Remote)
	if err != nil {
		logger.Errf("ssh error: remote host not found")
		return nil, re.EXECUTION_ERROR
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
		logger.Errf("ssh error: %s", err)
		return nil, re.EXECUTION_ERROR
	}

	session, err := client.NewSession()
	if err != nil {
		logger.Errf("ssh error: %s", err)
		return nil, re.EXECUTION_ERROR
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
		logger.Errf("ssh error: %s", err)
		return re.EXECUTION_ERROR
	}

	ns.Stdout = &session.stdout
	if err = ns.Run(command); err != nil {
		logger.Errf("failed to run command: %s", err)
		return re.EXECUTION_ERROR
	}

	if stdout {
		fmt.Println(session.stdout.String())
	}

	return nil
}

func (session *RemoteSession) CopyToRemote(local string, remote string) error {
	client, err := scp.NewClientFromExistingSSH(session.client, &scp.ClientOption{})
	if err != nil {
		logger.Errf("ssh error: %s", err)
		return re.EXECUTION_ERROR
	}

	_, err = client.NewSession()
	if err != nil {
		logger.Errf("ssh error: %s", err)
	}

	stat, err := os.Stat(local)
	if err == nil {
		if stat.IsDir() {
			// copy recursively
			err = client.CopyDirToRemote(local, remote, &scp.DirTransferOption{
				Context:      context.Background(),
				PreserveProp: true,
			})
			if err != nil {
				logger.Errf("ssh error: %s", err)
				return re.EXECUTION_ERROR
			}

			return nil
		} else {
			// copy file
			err = client.CopyFileToRemote(local, remote, &scp.FileTransferOption{
				Context:      context.Background(),
				PreserveProp: true,
			})
			if err != nil {
				logger.Errf("ssh error: %s", err)
				return re.EXECUTION_ERROR
			}

			return nil
		}
	}

	if os.IsNotExist(err) {
		logger.Errf("local path not found: %s", err)
		return re.EXECUTION_ERROR
	}

	logger.Errf("unknown error")
	return re.EXECUTION_ERROR

}

func (session *RemoteSession) CopyToLocal(local string, remote string, isDir bool) error {
	client, err := scp.NewClientFromExistingSSH(session.client, &scp.ClientOption{})
	if err != nil {
		logger.Errf("ssh error: %s", err)
		return re.EXECUTION_ERROR
	}

	_, err = client.NewSession()
	if err != nil {
		logger.Errf("ssh error: %s", err)
	}

	defer client.Close()

	if isDir {
		// copy recursively
		err = client.CopyDirFromRemote(remote, local, &scp.DirTransferOption{
			Context:      context.Background(),
			PreserveProp: true,
		})
		if err != nil {
			logger.Errf("ssh error: %s", err)
			return re.EXECUTION_ERROR
		}

		return nil
	} else {
		// copy file
		err = client.CopyFileFromRemote(remote, local, &scp.FileTransferOption{
			Context:      context.Background(),
			PreserveProp: true,
		})
		if err != nil {
			logger.Errf("ssh error: %s", err)
			return re.EXECUTION_ERROR
		}

		return nil
	}
}

func (session *RemoteSession) Close() {
	if session.session != nil {
		session.session.Close()
	}

	if session.client != nil {
		session.client.Close()
	}
}
