package lexer

import (
	"mintsql/internal/sql/token"
	"strconv"
	"strings"
	"unicode"
)

const symbols = ";*,()"

var keywords = []string{
	string(token.As),
	string(token.From),
	string(token.Create),
	string(token.Insert),
	string(token.Int),
	string(token.Select),
	string(token.Into),
	string(token.Table),
	string(token.Text),
	string(token.Values),
}

func LexBegin(l *Lexer) LexFn {
	for {
		switch nxt := l.Next(); {
		case nxt == Newline:
			l.Location.Line++
			l.Ignore()
		case unicode.IsSpace(nxt):
			l.Ignore()
		case nxt == '.' || unicode.IsDigit(nxt):
			l.Backup()
			return LexNumeric(l)
		case nxt == '\'' || nxt == '"':
			l.Backup()
			return LexString(l)
		case unicode.IsLetter(nxt):
			l.Backup()
			return LexKeyword(l) // -> LexIdentifier
		case nxt == '_':
			l.Backup()
			return LexIdentifier(l)
		case strings.ContainsRune(symbols, nxt):
			l.Backup()
			return LexSymbol(l)
		case nxt == EOF:
			l.Emit(token.KindEof)
			return nil
		}
	}
}

func LexNumeric(l *Lexer) LexFn {
	digits := "1234567890"
	if l.AcceptMany(digits) > 0 {
		l.AcceptOne(".")
		l.AcceptMany(digits)
	} else if l.AcceptOne(".") {
		if l.AcceptMany(digits) == 0 {
			return l.Error()
		}
	} else {
		return l.Error()
	}
	if !l.AcceptOne("eE") {
		l.Emit(token.KindNumeric)
		return LexBegin
	}
	l.AcceptOne("+-")
	if l.AcceptMany(digits) > 0 {
		l.Emit(token.KindNumeric)
		return LexBegin
	}
	return l.Error()
}

func LexString(l *Lexer) LexFn {
	if !l.AcceptOne("'\"") {
		return l.Error()
	}
	l.Ignore()
	for c := l.Next(); !unicode.IsControl(c); c = l.Next() {
		if c == '\'' || c == '"' {
			l.Backup()
			if _, err := strconv.Unquote(l.Current()); err != nil {
				break
			}
			l.Emit(token.KindString)
			l.Next()
			l.Ignore()
			return LexBegin
		}
	}
	return l.Error()
}

func LexIdentifier(l *Lexer) LexFn {
	if c := l.Next(); c != '_' && !unicode.IsLetter(c) {
		return l.Error()
	}
	for c := l.Next(); c == '_' || unicode.IsLetter(c) || unicode.IsNumber(c); c = l.Next() {
	}
	l.Backup()
	l.Emit(token.KindIdentifier)
	return LexBegin
}

func LexSymbol(l *Lexer) LexFn {
	if l.AcceptOne(symbols) {
		l.Emit(token.KindSymbol)
		return LexBegin
	}
	return l.Error()
}

func LexKeyword(l *Lexer) LexFn {
	matches := keywords
	cnt := 0
	for len(matches) > 0 {
		c := l.Next()
		if !unicode.IsLetter(c) {
			break
		}
		newMatches := make([]string, 0)
		for _, m := range matches {
			if rune(m[cnt]) == unicode.ToLower(c) {
				newMatches = append(newMatches, m)
			}
			// TODO longest match
			if m == strings.ToLower(l.Current()) && unicode.IsSpace(l.Peek()) {
				l.Emit(token.KindKeyword)
				return LexBegin
			}
		}
		cnt++
		matches = newMatches
	}
	l.Backup()
	return LexIdentifier
}
