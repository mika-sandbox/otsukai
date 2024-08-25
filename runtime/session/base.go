package session

type ISession interface {
	Run(command string, stdout bool) error
	Close()
}
