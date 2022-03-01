package token

import "onix/lex/scanner"

type Token int

const (
	_ = iota

	FILE_END
	NOBODY

	litBegin

	INT
	FLOAT
	STRING

	litEnd

	operatorBegin

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	operatorEnd

	keywordBegin

	BREAK
	CASE
	CHAN
	CONST
	CONTINUE
	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR
	FUNC
	GO
	GOTO
	IF
	IMPORT
	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN
	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR

	keywordEnd
)

var tokens = [...]string{

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK:       "break",
	CASE:        "case",
	CHAN:        "chan",
	CONST:       "const",
	CONTINUE:    "continue",
	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",
	FUNC:        "func",
	GO:          "go",
	GOTO:        "goto",
	IF:          "if",
	IMPORT:      "import",
	INTERFACE:   "interface",
	MAP:         "map",
	PACKAGE:     "package",
	RANGE:       "range",
	RETURN:      "return",
	SELECT:      "select",
	STRUCT:      "struct",
	SWITCH:      "switch",
	TYPE:        "type",
	VAR:         "var",
}

var tokenRev map[string]Token

func init() {
	tokenRev = make(map[string]Token)
	for k, v := range tokens {
		tokenRev[v] = Token(k)
	}
}

func Lookup(tok string) Token {
	for k, v := range tokenRev {
		if tok == k {
			return v
		}
	}
	return NOBODY
}

func IsKeyword(tok Token) bool {
	return tok > keywordBegin && tok < keywordEnd
}

func IsLiteral(s string) bool {
	if s == "" {
		return false
	}
	// literanl can not be keyword
	if IsKeyword(Lookup(s)) {
		return false
	}
	runes := []rune(s)
	if !scanner.IsLit(runes[0]) {
		return false
	}
	for _, ri := range runes {
		if !scanner.IsLit(ri) && !scanner.IsDecimal(ri) {
			return false
		}
	}
	return true
}
