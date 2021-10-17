package lexer

import (
	"mintsql/internal/sql/token"
	"strings"
	"unicode"
)

func lexBegin(l *Lexer) LexFn {
	for {
		switch nxt := l.peek(); {
		case nxt == NEWLINE:
			l.location.Line++
			l.next()
			l.ignore()
		case unicode.IsSpace(nxt):
			l.next()
			l.ignore()
		case nxt == '.' || unicode.IsDigit(nxt):
			return lexNumeric(l)
		case nxt == '\'' || nxt == '"':
			return lexString(l)
		case unicode.IsLetter(nxt):
			return lexKeyword(l)
		case nxt == '_':
			return lexIdentifier(l)
		case strings.ContainsRune(token.Symbols, nxt):
			return lexSymbol(l)
		case nxt == EOF:
			l.next()
			l.emit(token.KindEof)
			return nil
		}
	}
}

func lexNumeric(l *Lexer) LexFn {
	digits := "1234567890"
	if l.acceptManyIn(digits) > 0 {
		l.acceptOneIn(".")
		l.acceptManyIn(digits)
	} else if l.acceptOneIn(".") > 0 {
		if l.acceptManyIn(digits) == 0 {
			return l.err("bad numeric")
		}
	} else {
		return l.err("bad numeric")
	}

	if l.acceptOneIn("eE") < 0 {
		l.emit(token.KindNumeric)
		return lexBegin
	}
	l.acceptOneIn("+-")
	if l.acceptManyIn(digits) > 0 {
		l.emit(token.KindNumeric)
		return lexBegin
	}
	return l.err("bad numeric")
}

func lexString(l *Lexer) LexFn {
	if l.acceptOneIn("'\"") < 0 {
		return l.err("bad string")
	}
	l.ignore()
	for c := l.next(); !unicode.IsControl(c); c = l.next() {
		if c == '\\' {
			if r := l.peek(); r == '\'' || r == '"' {
				l.next()
				continue
			}
		}
		if c == '\'' || c == '"' {
			if r := l.peek(); r == c {
				l.next()
				continue
			}
			l.backup()
			l.emit(token.KindString)
			l.next()
			l.ignore()
			return lexBegin
		}
	}
	return l.err("bad string")
}

func lexIdentifier(l *Lexer) LexFn {
	if l.acceptOneIf(func(r rune) bool { return r == '_' || unicode.IsLetter(r) }) < 0 {
		return l.err("bad identifier")
	}
	l.acceptManyIf(func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) })
	l.emit(token.KindIdentifier)
	return lexBegin
}

func lexSymbol(l *Lexer) LexFn {
	if l.acceptOneIn(token.Symbols) > 0 {
		l.emit(token.KindSymbol)
		return lexBegin
	}
	return l.err("unrecognizable symbol")
}

func lexKeyword(l *Lexer) LexFn {
	c := l.next()
	var match string
	candidates := token.Keywords

	for i := 0; unicode.IsLetter(c); i++ {
		newCandidates := make([]string, 0)
		anyMatch := false
		for _, can := range candidates {
			if rune(can[i]) == unicode.ToLower(c) {
				anyMatch = true
				if i == len(can)-1 {
					match = can
					continue
				}
				newCandidates = append(newCandidates, can)
			}
		}
		if !anyMatch {
			break
		}
		candidates = newCandidates
		c = l.next()
	}
	l.backup()

	if match != "" {
		l.emit(token.KindKeyword)
		return lexBegin
	}
	return lexIdentifier
}
