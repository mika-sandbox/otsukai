package parser

import (
	"errors"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"regexp"
	"slices"
)

type (
	Identifier struct {
		Name string `parser:"@Identifier"`
	}

	Pair struct {
		Identifier HashIdentifier `parser:"@@"`
		Value      Value          `parser:"@@"`
	}

	Literal struct {
		String  *string  `parser:"@String  |"`
		Number  *float64 `parser:"@Number  |"`
		Boolean *bool    `parser:"('true' | 'false') |"`
		Null    bool     `parser:"'nil'"`
	}

	HashIdentifier struct {
		Identifier string `parser:"@Identifier ':'"`
	}

	HashSymbol struct {
		Identifier string `parser:"':' @Identifier"`
	}

	HashObject struct {
		OpenBracket  string `parser:"'{'"`
		Pairs        []Pair `parser:"@@ (',' @@)+"`
		CloseBracket string `parser:"'}'"`
	}

	Value struct {
		Literal    *Literal    `parser:"@@ |"`
		HashSymbol *HashSymbol `parser:"@@ |"`
		Hash       *HashObject `parser:"@@"`
	}

	IdentifierNameExpression struct {
		Identifier Identifier `parser:"@@"`
	}

	MemberAccessExpressionTExpression struct {
		MemberAccessExpression *MemberAccessExpression        `parser:"@@ |"`
		InvocationExpression   *InvocationExpressionWithParen `parser:"@@"`
	}

	MemberAccessExpression struct {
		Expression MemberAccessExpressionTExpression `parser:"@@"`
		Dot        string                            `parser:"'.'"`
		Identifier Identifier                        `parser:"@@"`
	}

	LambdaExpression struct {
		DoKeyword  string      `parser:"'do'"`
		NewLine    string      `parser:"@NewLine?"`
		Statements []Statement `parser:"@@*"`
		EndKeyword string      `parser:"'end'"`
	}

	Argument struct {
		Identifier *string    `parser:"(@Identifier ':')?"`
		Expression Expression `parser:"@@"`
	}

	ArgumentList struct {
		Argument         []Argument        `parser:"(@@ (',' @@)*)?"`
		LambdaExpression *LambdaExpression `parser:"@@?"`
	}

	InvocationExpression struct {
		Expression   ExpressionWithIdentifier `parser:"@@"`
		ArgumentList ArgumentList             `parser:"@@"`
	}

	InvocationExpressionWithParen struct {
		Expression   ExpressionWithIdentifier `parser:"@@"`
		OpenParen    string                   `parser:"'('"`
		ArgumentList ArgumentList             `parser:"@@"`
		CloseParen   string                   `parser:"')'"`
	}

	ValueExpression struct {
		Value Value `parser:"@@"`
	}

	ExpressionWithIdentifier struct {
		IdentifierNameExpression *IdentifierNameExpression `parser:"@@"`
		// Expression               *Expression               `parser:"@@"`
	}

	Expression struct {
		// MemberAccessExpression        *MemberAccessExpression        `parser:"@@ |"`
		LambdaExpression              *LambdaExpression              `parser:"@@ |"`
		InvocationExpression          *InvocationExpression          `parser:"@@ |"`
		InvocationExpressionWithParen *InvocationExpressionWithParen `parser:"@@ |"`
		IfExpression                  *IfStatementOrExpression       `parser:"@@ |"`
		ValueExpression               *ValueExpression               `parser:"@@"`
	}

	BlockStatement struct {
		OpenBracket  string      `parser:"'{'"`
		Statements   []Statement `parser:"@@"`
		CloseBracket string      `parser:"'}'"`
	}

	IfStatementOrExpressionConditionExpression struct {
		InvocationExpression          *InvocationExpression          `parser:"@@ |"`
		InvocationExpressionWithParen *InvocationExpressionWithParen `parser:"@@ |"`
		ValueExpression               *ValueExpression               `parser:"@@"`
	}

	IfStatementOrExpression struct {
		IfKeyword  string                                     `parser:"'if'"`
		OpenParen  *string                                    `parser:"'('?"`
		Condition  IfStatementOrExpressionConditionExpression `parser:"@@"`
		CloseParen *string                                    `parser:"')'?"`
		NewParen   *string                                    `parser:"@NewLine?"`
		Statements []Statement                                `parser:"@@*"`
		EndKeyword string                                     `parser:"'end'"`
	}

	ExpressionStatement struct {
		Expression Expression `parser:"@@"`
	}

	StatementInternal struct {
		BlockStatement      *BlockStatement          `parser:"@@ |"`
		IfStatement         *IfStatementOrExpression `parser:"@@ |"`
		ExpressionStatement *ExpressionStatement     `parser:"@@"`
	}

	Statement struct {
		Statement StatementInternal `parser:"@@"`
		NewLine   string            `parser:"@NewLine"`
	}

	CompilationUnit struct {
		Statements []Statement `parser:"@@*"`
	}
)

type Entry struct {
	Statements []Statement `parser:"@@*"`
}

var RubyLikeLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Identifier", `[A-Za-z_][A-Za-z0-9_]*`},
	{"HashIdentifier", `[A-Za-z_][A-Za-z0-9_]+:`},
	{"HashKeyword", `[A-Za-z0-9_]+`},
	{`String`, `'[^']*'|"[^"]*"`},
	{"Number", `(?:\d*\.)?\d+`},
	{"Comment", `#.*(\n)?`},
	{"Whitespace", `[ \t]+`},
	{"NewLine", `[\r\n]+`},

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
	participle.UseLookahead(2),
)

var Keywords = []string{
	"do",
	"if",
	"end",
}

func (v *Identifier) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()

	if token.EOF() {
		return errors.New("unexpected eof")
	}

	if slices.Contains(Keywords, token.Value) {
		return participle.NextMatch
	}

	identifier, err := regexp.Compile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if err != nil {
		return err
	}

	matches := identifier.FindAllString(token.Value, -1)
	if len(matches) != 1 {
		return participle.NextMatch
	}

	next := lex.Next()
	if next.EOF() {
		return errors.New("unexpected eof")
	}

	*v = Identifier{
		Name: next.Value,
	}

	return nil
}
