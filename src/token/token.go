package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // e.g., x, foo
	INT   = "INT"
	REAL  = "REAL"
	BOOL  = "BOOLEAN"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	WHILE    = "WHILE"

	// Brackets
	LPAREN = "("
	RPAREN = "}"
	LBRACE = "{"
	RBRACE = "}"

	// Comparison operators
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LTE    = "<="
	GTE    = ">="

	// Comments
	COMMENT                 = "//"
	MULTILINE_COMMENT_START = "/*"
	MULTILINE_COMMENT_END   = "*/"
)

func New(tokenType TokenType, literal string, line, column int) Token {
	return Token{Type: tokenType, Literal: literal, Line: line, Column: column}
}
