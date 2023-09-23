package lexer

import (
	"strings"

	"github.com/hazed7/compiler/src/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.New(token.EQ, string(ch)+string(l.ch), l.line, l.column)
		} else {
			tok = token.New(token.ASSIGN, string(l.ch), l.line, l.column)
		}
	case '+':
		tok = token.New(token.PLUS, string(l.ch), l.line, l.column)
	case '-':
		tok = token.New(token.MINUS, string(l.ch), l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.New(token.NOT_EQ, string(ch)+string(l.ch), l.line, l.column)
		} else {
			tok = token.New(token.BANG, string(l.ch), l.line, l.column)
		}
	case '*':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token.New(token.MULTILINE_COMMENT_END, "*/", l.line, l.column-len(tok.Literal))
		} else {
			tok = token.New(token.ASTERISK, string(l.ch), l.line, l.column)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token.New(token.COMMENT, l.readComment(), l.line, l.column-len(tok.Literal))
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = token.New(token.MULTILINE_COMMENT_START, l.readComment(), l.line, l.column-len(tok.Literal))
		} else {
			tok = token.New(token.SLASH, string(l.ch), l.line, l.column)
		}
	case '<':
		tok = token.New(token.LT, string(l.ch), l.line, l.column)
	case '>':
		tok = token.New(token.GT, string(l.ch), l.line, l.column)
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), l.line, l.column)
	case ':':
		tok = token.New(token.COLON, string(l.ch), l.line, l.column)
	case ',':
		tok = token.New(token.COMMA, string(l.ch), l.line, l.column)
	case '(':
		tok = token.New(token.LPAREN, string(l.ch), l.line, l.column)
	case ')':
		tok = token.New(token.RPAREN, string(l.ch), l.line, l.column)
	case '{':
		tok = token.New(token.LBRACE, string(l.ch), l.line, l.column)
	case '}':
		tok = token.New(token.RBRACE, string(l.ch), l.line, l.column)
	case 't':
		if l.peekChar() == 'r' && l.peekCharAt(2) == 'u' && l.peekCharAt(3) == 'e' && !isLetter(l.peekCharAt(4)) {
			tok = token.New(token.TRUE, "true", l.line, l.column)
			// Consume the "true" characters
			l.readChar()
			l.readChar()
			l.readChar()
			l.readChar()
		} else {
			tok = token.New(token.IDENT, l.readIdentifier(), l.line, l.column-len(tok.Literal))
		}
	case 'f':
		if l.peekChar() == 'a' && l.peekCharAt(2) == 'l' && l.peekCharAt(3) == 's' && l.peekCharAt(4) == 'e' && !isLetter(l.peekCharAt(5)) {
			tok = token.New(token.FALSE, "false", l.line, l.column)
			// Consume the "false" characters
			l.readChar()
			l.readChar()
			l.readChar()
			l.readChar()
			l.readChar()
		} else {
			tok = token.New(token.IDENT, l.readIdentifier(), l.line, l.column-len(tok.Literal))
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok = token.New(token.TokenType(l.readIdentifier()), l.readIdentifier(), l.line, l.column-len(tok.Literal))
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.REAL
			} else {
				tok.Type = token.INT
			}
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			return tok
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch), l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		l.position = l.readPosition
		l.readPosition++
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) peekCharAt(offset int) byte {
	if l.position+offset >= len(l.input) {
		return 0
	}

	return l.input[l.position+offset]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
