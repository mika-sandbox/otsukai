package otsukai

import (
	"errors"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"slices"
)

type (
	Boolean                  bool
	WithoutKeywordIdentifier struct {
		Identifier string `parse:"@Identifier"`
	}

	// identifier:
	HashIdentifier struct {
		Identifier   string   `parser:"@Identifier"`
		ColonKeyword struct{} `parser:"':'"`
	}

	// :keyword
	HashKeyword struct {
		ColonKeyword struct{} `parser:"':'"`
		Identifier   string   `parser:"@Identifier"`
	}

	// { object: any_value ... }
	HashObject struct {
		OpenBracket  struct{}          `parser:"'{'"`
		Hashes       []*HashExpression `parser:"@@ (',' @@)*"`
		CloseBracket struct{}          `parser:"'}'"`
	}

	// any values
	Literal struct {
		Number  *float64 `parser:"@Number |"`
		String  *string  `parser:"@String |"`
		Boolean *Boolean `parser:"@('true' | 'false') |"`
		Null    bool     `parser:"'nil'"`
	}

	Value struct {
		Literal *Literal     `parser:"@@ |"`
		Keyword *HashKeyword `parser:"@@ |"`
		Object  *HashObject  `parser:"@@"`
	}

	Expression struct {
		Literal                         *Value                           `parser:"@@ |"`
		Hash                            *HashExpression                  `parser:"@@ |"`
		Keyword                         *string                          `parser:"@HashKeyword |"`
		InvocationWithoutArgsExpression *InvocationWithoutArgsExpression `parser:"@@ |"`
		InvocationWithArgsExpression    *InvocationWithArgsExpression    `parser:"@@"`
	}

	DoExpressionWithoutStatements struct {
		DoKeyword  struct{} `parser:"'do'"`
		EndKeyword struct{} `parser:"'end'"`
	}

	// do ... end
	DoExpressionWithStatements struct {
		DoKeyword  struct{}    `parser:"'do'"`
		Statements []Statement `parser:"@@*"`
		EndKeyword struct{}    `parser:"'end'"`
	}

	DoExpression struct {
		WithoutStatements *DoExpressionWithoutStatements `parser:"@@ |"`
		WithStatements    *DoExpressionWithStatements    `parser:"@@"`
	}

	// method(...)
	// method ...
	// method ... do ... end
	InvocationWithArgsExpression struct {
		Identifier       WithoutKeywordIdentifier `parser:"@@"`
		OpenParen        *struct{}                `parser:"'('?"`
		Arguments        []Expression             `parser:"@@ (',' @@)*"`
		CloseParen       *struct{}                `parser:"')'?"`
		LambdaExpression *DoExpression            `parser:"@@?"`
	}

	// method()
	// method
	// method do ... end
	InvocationWithoutArgsExpression struct {
		Identifier       WithoutKeywordIdentifier `parser:"@@"`
		OpenParen        struct{}                 `parser:"'('"`
		CloseParen       struct{}                 `parser:"')'"`
		LambdaExpression *DoExpression            `parser:"@@?"`
	}

	// hash: ...
	HashExpression struct {
		Identifier HashIdentifier `parser:"@@"`
		Value      Value          `parser:"@@"`
	}

	// set key: xxx
	SetStatement struct {
		SetKeyword struct{}       `parser:"'set'"`
		Expression HashExpression `parser:"@@"`
	}

	// task :identifier do ... end
	TaskStatement struct {
		TaskKeyword struct{}     `parser:"'task'"`
		Identifier  HashKeyword  `parser:"@@"`
		DoStatement DoExpression `parser:"@@"`
	}

	// hook after: :identifier do ... end
	HookStatement struct {
		HookKeyword struct{}       `parser:"'hook'"`
		Parameters  HashExpression `parser:"@@"`
		DoStatement DoExpression   `parser:"@@"`
	}

	// statement
	ExpressionStatement struct {
		DoExpression                    *DoExpression                    `parser:"@@ |"`
		InvocationWithoutArgsExpression *InvocationWithoutArgsExpression `parser:"@@ |"`
		InvocationWithArgsExpression    *InvocationWithArgsExpression    `parser:"@@"`
	}

	// if cond ... end
	IfStatement struct {
		IfKeyword  struct{}    `parser:"'if'"`
		OpenParen  *struct{}   `parser:"'('?"`
		Condition  Expression  `parser:"@@"`
		CloseParen *struct{}   `parser:"')'?"`
		Statements []Statement `parser:"@@"`
		EndKeyword struct{}    `parser:"'end'"`
	}

	// any statement
	Statement struct {
		SetStatement        *SetStatement        `parser:"@@ |"`
		TaskStatement       *TaskStatement       `parser:"@@ |"`
		HookStatement       *HookStatement       `parser:"@@ |"`
		IfStatement         *IfStatement         `parser:"@@ |"`
		ExpressionStatement *ExpressionStatement `parser:"@@ "`
	}
)

func (v Value) ToValueObject() (IValueObject, error) {
	if v.Keyword != nil {
		return &StringValueObject{val: v.Keyword.Identifier}, nil
	}

	if v.Object != nil {
		items := map[string]IValueObject{}

		for _, hash := range v.Object.Hashes {
			items[hash.Identifier.Identifier], _ = hash.Value.ToValueObject()
		}

		return &HashValueObject{val: items}, nil
	}

	if v.Literal != nil {
		if v.Literal.String != nil {
			return &StringValueObject{val: *v.Literal.String}, nil
		}

		if v.Literal.Number != nil {
			return &Float64ValueObject{val: *v.Literal.Number}, nil
		}

		if v.Literal.Boolean != nil {
			return &BooleanValueObject{val: *v.Literal.Boolean == true}, nil
		}

		if v.Literal.Null {
			return nil, nil
		}
	}

	return nil, errors.New("invalid value")
}

type Entry struct {
	Statements []Statement `parser:"@@*"`
}

const (
	KIND_SET_STATEMENT = iota
	KIND_TASK_STATEMENT
	KIND_IF_STATEMENT
)

var RubyLikeLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Identifier", `[A-Za-z_][A-Za-z0-9_]*`},
	{"HashIdentifier", `[A-Za-z_][A-Za-z0-9_]+:`},
	{"HashKeyword", `[A-Za-z0-9_]+`},
	{`String`, `'[^']*'|"[^"]*"`},
	{"Number", `(?:\d*\.)?\d+`},
	{"Comment", `#.*(\n)?`},
	{"Whitespace", `[ \t\n\r]+`},

	// indirect references for lexer
	{"COLON", ":"},
	{"OPEN_PAREN", `\(`},
	{"CLOSE_PAREN", `\)`},
	{"OPEN_BRACKET", "{"},
	{"CLOSE_BRACKET", "}"},
	{"COMMA", ","},
})

var Parser = participle.MustBuild[Entry](
	participle.Lexer(RubyLikeLexer),
	participle.Elide("Comment", "Whitespace"),
	participle.Unquote("String"),
	participle.UseLookahead(2),
)

var Keywords = []string{
	"do",
	"if",
	"end",
}

func (v *WithoutKeywordIdentifier) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()

	if token.EOF() {
		return errors.New("unexpected eof")
	}

	if slices.Contains(Keywords, token.Value) {
		return participle.NextMatch
	}

	next := lex.Next()
	if next.EOF() {
		return errors.New("unexpected eof")
	}

	*v = WithoutKeywordIdentifier{
		Identifier: next.Value,
	}

	return nil
}
