package lexer

import (
	"mintsql/internal/sql/token"
	"strings"
	"unicode/utf8"
)

const (
	NEWLINE rune = 10
	EOF     rune = 0
)

type LexFn func(*Lexer) LexFn

type cursor struct {
	start          int
	pos            int
	width          int
	lastLineLength int
}

type Lexer struct {
	input  string
	tokens chan *token.Token
	state  LexFn
	ok     bool

	cursor
	location token.Location
}

func (l *Lexer) emit(kind token.Kind) {
	newToken := &token.Token{
		Kind:     kind,
		Value:    l.input[l.start:l.pos],
		Location: l.location,
	}
	l.tokens <- newToken
	l.start = l.pos
}

func (l *Lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) remainder() string {
	return l.input[l.pos:]
}

func (l *Lexer) backup() {
	if l.cursor.pos > l.cursor.start {
		l.cursor.pos -= l.cursor.width
		l.location.Column--
		if res, _ := utf8.DecodeRuneInString(l.input[l.pos:]); res == NEWLINE {
			l.location.Line--
			l.location.Column = l.cursor.lastLineLength
		}
	}
}

func (l *Lexer) next() (res rune) {
	if l.cursor.pos >= len(l.input) {
		l.cursor.width = 0
		return EOF
	}
	res, l.cursor.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.cursor.pos += l.cursor.width
	if res == NEWLINE {
		l.location.Line++
		l.cursor.lastLineLength = l.location.Column
		l.location.Column = 1
		return
	}
	l.location.Column++
	return
}

func (l *Lexer) peek() rune {
	if l.cursor.pos >= len(l.input) {
		return EOF
	}
	c, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return c
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) acceptOneIf(pred func(rune) bool) rune {
	r := l.next()
	if pred(r) {
		return r
	}
	l.backup()
	return -1
}

func (l *Lexer) acceptManyIf(pred func(rune) bool) (len int) {
	r := l.next()
	for ; pred(r); r = l.next() {
		len++
	}
	l.backup()
	return
}

func (l *Lexer) acceptOneIn(domain string) rune {
	r := l.next()
	if strings.ContainsRune(domain, r) {
		return r
	}
	l.backup()
	return -1
}

func (l *Lexer) acceptManyIn(domain string) (len int) {
	r := l.next()
	for ; strings.ContainsRune(domain, r); r = l.next() {
		len++
	}
	l.backup()
	return
}

func (l *Lexer) err() LexFn {
	l.emit(token.KindError)
	return nil
}

func (l *Lexer) NextToken() *token.Token {
	t, ok := <-l.tokens
	if ok {
		return t
	}
	return nil
}

func (l *Lexer) try(funcs ...LexFn) {
	for _, f := range funcs {
		lex := *l
		lex.state = f
		go lex.Lex()
		if !l.ok {
			*l = lex
			break
		}
	}
}

func (l *Lexer) Lex() {
	defer func() {
		close(l.tokens)
	}()
	for ; l.state != nil; l.state = l.state(l) {
	}
}

func New(src string) *Lexer {
	return &Lexer{
		input:    src,
		tokens:   make(chan *token.Token, len(src)/2),
		state:    lexBegin,
		cursor:   cursor{},
		location: token.Location{Line: 1},
	}
}
