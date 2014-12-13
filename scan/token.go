package scan

type Token int

const (
	ILLEGAL Token = iota
	EOF           // EOF

	BINARY  // '!', '%', '*', '/', '+', '|', '&', '^', '-', '>', '<', '=', '?', '\', '~':
	GLOBAL  // True, False, Nil
	IDENT   // aBool, not
	KEYWORD // ifTrue:

	ASSIGN // =
	DEFINE // ->

	LBRACE // {
	RBRACE // }
	LBRACK // [
	RBRACK // ]
	LPAREN // (
	RPAREN // )

	CASCADE // ;
	COMMA   // ,
	PERIOD  // .
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	BINARY:  "BINARY",
	GLOBAL:  "GLOBAL",
	IDENT:   "IDENT",
	KEYWORD: "KEYWORD",

	ASSIGN: ":=",
	DEFINE: "->",

	LBRACE: "{",
	RBRACE: "}",
	LBRACK: "[",
	RBRACK: "]",
	LPAREN: "(",
	RPAREN: ")",

	CASCADE: ";",
	COMMA:   ",",
	PERIOD:  ".",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	return s
}
