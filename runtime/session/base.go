package session

type ISession interface {
	Run(command string, stdout bool) error
	CopyToRemote(local string, remote string) error
	CopyToLocal(local string, remote string, isDir bool) error
	Close()
}
