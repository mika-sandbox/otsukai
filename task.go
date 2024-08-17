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
