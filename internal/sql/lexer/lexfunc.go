package lexer

import (
	"mintsql/internal/sql/token"
	"strings"
	"unicode"
)

func LexBegin(l *Lexer) LexFn {
	for {
		switch nxt := l.Next(); {
		case nxt == NEWLINE:
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
			return LexKeyword(l)
		case nxt == '_':
			l.Backup()
			return LexIdentifier(l)
		case strings.ContainsRune(token.Symbols, nxt):
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
	if l.AcceptManyIn(digits) > 0 {
		l.AcceptOneIn(".")
		l.AcceptManyIn(digits)
	} else if l.AcceptOneIn(".") {
		if l.AcceptManyIn(digits) == 0 {
			return l.Err("bad numeric")
		}
	} else {
		return l.Err("bad numeric")
	}
	if !l.AcceptOneIn("eE") {
		l.Emit(token.KindNumeric)
		return LexBegin
	}
	l.AcceptOneIn("+-")
	if l.AcceptManyIn(digits) > 0 {
		l.Emit(token.KindNumeric)
		return LexBegin
	}
	return l.Err("bad numeric")
}

func LexString(l *Lexer) LexFn {
	if !l.AcceptOneIn("'\"") {
		return l.Err("bad string")
	}
	l.Ignore()
	for c := l.Next(); !unicode.IsControl(c); c = l.Next() {
		if c == '\'' || c == '"' {
			l.Backup()
			l.Emit(token.KindString)
			l.Next()
			l.Ignore()
			return LexBegin
		}
	}
	return l.Err("bad string")
}

func LexIdentifier(l *Lexer) LexFn {
	if c := l.Next(); c != '_' && !unicode.IsLetter(c) {
		return l.Err("bad identifier")
	}
	for c := l.Next(); c == '_' || unicode.IsLetter(c) || unicode.IsNumber(c); c = l.Next() {
	}
	l.Backup()
	l.Emit(token.KindIdentifier)
	return LexBegin
}

func LexSymbol(l *Lexer) LexFn {
	if l.AcceptOneIn(token.Symbols) {
		l.Emit(token.KindSymbol)
		return LexBegin
	}
	return l.Err("unrecognizable symbol")
}

func LexKeyword(l *Lexer) LexFn {
	c := l.Next()
	var match string
	candidates := token.Keywords

	for i := 0; len(candidates) > 0 && unicode.IsLetter(c); i++ {
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
		c = l.Next()
	}
	l.Backup()

	if match != "" {
		l.Emit(token.KindKeyword)
		return LexBegin
	}
	return LexIdentifier
}
