package lexer

import (
	"fmt"
	"mintsql/internal/sql/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

type LexFn func(*Lexer) LexFn

type Cursor struct {
	Start int
	Pos   int
	Width int
}

type Lexer struct {
	Name   string
	Input  string
	Tokens chan token.Token
	State  LexFn

	Cursor
	token.Location
}

func (l *Lexer) Emit(kind token.Kind) {
	l.Tokens <- token.Token{
		Kind:     kind,
		Value:    l.Input[l.Start:l.Pos],
		Location: l.Location,
	}
	l.Start = l.Pos
}

func (l *Lexer) Current() string {
	return l.Input[l.Start:l.Pos]
}

func (l *Lexer) Remainder() string {
	return l.Input[l.Pos:]
}

func (l *Lexer) Backup() {
	l.Pos -= l.Width
	l.Column -= l.Width
}

func (l *Lexer) Next() rune {
	if l.Pos >= len(l.Input) {
		l.Width = 0
		return token.EOF
	}

	res, width := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width = width
	l.Pos += l.Width
	l.Column++
	if res == token.Newline {
		l.Line++
		l.Column = 0
	}
	return res
}

func (l *Lexer) Peek() rune {
	c := l.Next()
	l.Backup()
	return c
}

func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

func (l *Lexer) IsWhiteSpace() bool {
	c, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return unicode.IsSpace(c)
}

func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptMany(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

func (l *Lexer) SkipWhiteSpace() {
	for {
		c := l.Next()
		if c == token.EOF {
			l.Emit(token.KindEof)
			break
		}
		if !unicode.IsSpace(c) {
			l.Backup()
			break
		}
	}
}

func (l *Lexer) Errorf(format string, args ...interface{}) LexFn {
	l.Tokens <- token.Token{
		Kind:     token.KindError,
		Value:    fmt.Sprintf(format, args...),
		Location: l.Location,
	}
	return nil
}

func (l *Lexer) NextToken() token.Token {
	for {
		select {
		case t := <-l.Tokens:
			return t
		default:
			l.State = l.State(l)
		}
	}
}

func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

func (l *Lexer) Run(init LexFn) {
	defer l.Shutdown()
	for st := init; st != nil; {
		st = st(l)
	}
}
