package scan

import (
	"fmt"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/socialmachines/soma/file"
)

type ErrorHandler func(pos file.Position, msg string)

type Scanner struct {
	file *file.File   // source file handle
	dir  string       // directory portion of file.Name()
	src  []byte       // source
	err  ErrorHandler // error reporting; or nil

	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
}

func (s *Scanner) Init(file *file.File, src []byte, err ErrorHandler) {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}

	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0

	s.next()
}

func (s *Scanner) Scan() (pos file.Pos, tok Token, lit string) {
	s.skipWhitespace()

	pos = s.file.Pos(s.offset)

	switch ch := s.ch; {
	case isLetter(ch):
		tok, lit = s.scanName()
	case isBinary(ch):
		tok, lit = s.scanBinary()
	default:
		s.next()
		switch ch {
		case -1:
			tok, lit = EOF, "EOF"
		case '=':
			tok, lit = ASSIGN, "="
		case '{':
			tok, lit = LBRACE, "{"
		case '}':
			tok, lit = RBRACE, "}"
		case '[':
			tok, lit = LBRACK, "["
		case ']':
			tok, lit = RBRACK, "]"
		case '(':
			tok, lit = LPAREN, "("
		case ')':
			tok, lit = RPAREN, ")"
		case ';':
			tok, lit = CASCADE, ";"
		case ',':
			tok, lit = COMMA, ","
		case '.':
			tok, lit = PERIOD, "."
		default:
			s.error(s.file.Offset(pos), fmt.Sprintf("illegal character %#U", ch))
			tok, lit = ILLEGAL, string(ch)
		}
	}
	return
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func isLetter(ch rune) bool {
	return isUpper(ch) || isLower(ch)
}

func isUpper(ch rune) bool {
	return 'A' <= ch && ch <= 'Z' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isLower(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || ch >= 0x80 && unicode.IsLetter(ch)
}

func (s *Scanner) scanName() (Token, string) {
	scanIdent := func (s *Scanner) string {
		offs := s.offset
		for isLetter(s.ch) {
			s.next()
		}
		return string(s.src[offs:s.offset])
	}

	if isUpper(s.ch) {
		return GLOBAL, scanIdent(s)
	}

	lit := scanIdent(s)
	if s.ch == ':' {
		s.next()
		return KEYWORD, lit+":"
	} else {
		return IDENT, lit
	}
}

func isBinary(ch rune) bool {
	switch ch {
		case '!', '%', '*', '/', '+', '|', '&', '^', '-', '>', '<', '=', '?', '\\', '~':
		return true
	}
	return false
}

func (s *Scanner) scanBinary() (Token, string) {
	offs := s.offset
	for isBinary(s.ch) {
		s.next()
	}

	bin := string(s.src[offs:s.offset])
	if bin == "->" {
		return DEFINE, bin
	}
	return BINARY, bin
}

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}

		r, w := rune(s.src[s.rdOffset]), 1
		switch {
			case r == 0:
			s.error(s.offset, "illegal character NULL")
			case r >= 0x80:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = -1 // eof
	}
}

func (s *Scanner) error(offset int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offset)), msg)
	}
}
