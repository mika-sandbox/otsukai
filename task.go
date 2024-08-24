package otsukai

type Task struct {
	Name        string
	Statements  []Statement
	BeforeHooks []Hook
	AfterHooks  []Hook
}

type Hook struct {
	Statements []Statement
}

func CreateTask(name string, statements []Statement) Task {
	return Task{
		Name:        name,
		Statements:  statements,
		BeforeHooks: make([]Hook, 0),
		AfterHooks:  make([]Hook, 0),
	}
}

func CreateHook(statements []Statement) Hook {
	return Hook{
		Statements: statements,
	}
}
