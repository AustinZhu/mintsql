package lexer

import (
	"errors"
	"mintsql/internal/sql/token"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	NEWLINE rune = 10
	EOF     rune = 0
)

type LexFn func(*Lexer) LexFn

type Cursor struct {
	Start          int
	Pos            int
	Width          int
	lastLineLength int
}

type Lexer struct {
	Input  string
	Tokens chan *token.Token
	State  LexFn
	Error  error

	Cursor
	token.Location
}

func (l *Lexer) Emit(kind token.Kind) {
	newToken := &token.Token{
		Kind:     kind,
		Value:    l.Input[l.Start:l.Pos],
		Location: l.Location,
	}
	l.Tokens <- newToken
	l.Start = l.Pos
}

func (l *Lexer) Current() string {
	return l.Input[l.Start:l.Pos]
}

func (l *Lexer) Remainder() string {
	return l.Input[l.Pos:]
}

func (l *Lexer) Last() rune {
	res, _ := utf8.DecodeRuneInString(l.Input[l.Pos-l.Width : l.Pos])
	return res
}

func (l *Lexer) Backup() {
	if l.Cursor.Pos > l.Cursor.Start {
		l.Cursor.Pos -= l.Cursor.Width
		l.Location.Column--
		if res, _ := utf8.DecodeRuneInString(l.Input[l.Pos:]); res == NEWLINE {
			l.Location.Line--
			l.Location.Column = l.Cursor.lastLineLength
		}
	}
}

func (l *Lexer) Next() (res rune) {
	if l.Cursor.Pos >= len(l.Input) {
		l.Cursor.Width = 0
		return EOF
	}
	res, l.Cursor.Width = utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Cursor.Pos += l.Cursor.Width
	l.Location.Column++
	return res
}

func (l *Lexer) Peek() rune {
	if l.Cursor.Pos >= len(l.Input) {
		return EOF
	}
	c, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return c
}

func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

func (l *Lexer) AcceptOneIf(pred func(rune) bool) rune {
	r := l.Next()
	if pred(r) {
		return r
	}
	l.Backup()
	return -1
}

func (l *Lexer) AcceptManyIf(pred func(rune) bool) (len int) {
	r := l.Next()
	for ; pred(r); r = l.Next() {
		len++
	}
	l.Backup()
	return len
}

func (l *Lexer) AcceptOneIn(domain string) rune {
	r := l.Next()
	if strings.ContainsRune(domain, r) {
		return r
	}
	l.Backup()
	return -1
}

func (l *Lexer) AcceptManyIn(domain string) (len int) {
	r := l.Next()
	for ; strings.ContainsRune(domain, r); r = l.Next() {
		len++
	}
	l.Backup()
	return len
}

func (l *Lexer) Err(msg string) LexFn {
	l.Emit(token.KindError)
	l.Error = errors.New(msg)
	return nil
}

func (l *Lexer) NextToken() *token.Token {
	t, ok := <-l.Tokens
	if ok {
		return t
	}
	return nil
}

func (l *Lexer) Try(funcs ...LexFn) {
	for _, f := range funcs {
		lex := *l
		lex.State = f
		go lex.Run()
		if l.Error == nil {
			*l = lex
			break
		}
	}
}

func (l *Lexer) Run() {
	defer func() {
		l.Shutdown()
	}()
	for ; l.State != nil; l.State = l.State(l) {
	}
}

func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

func New(src string, init LexFn) *Lexer {
	return &Lexer{
		Input:    src,
		Tokens:   make(chan *token.Token, len(src)/2),
		State:    init,
		Cursor:   Cursor{},
		Location: token.Location{Line: 1},
	}
}

func NewFromFile(path string, init LexFn) *Lexer {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &Lexer{
		Input:    string(src),
		Tokens:   make(chan *token.Token, len(src)/2),
		State:    init,
		Cursor:   Cursor{},
		Location: token.Location{Line: 1},
	}
}
