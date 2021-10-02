package lexer

import "fmt"

type location struct {
	line uint
	col  uint
}

type cursor struct {
	pointer uint
	loc     location
}

type lexer func(string, cursor) (*token, cursor, bool)

var lexers = []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}

func lex(source string) ([]*token, error) {
	var tokens []*token
	var cur cursor
	for cur.pointer < uint(len(source)) {
		var success bool
		for _, l := range lexers {
			var tk *token
			var nextCur cursor
			if tk, nextCur, success = l(source, cur); success {
				cur = nextCur
				if tk != nil {
					tokens = append(tokens, tk)
				}
				break
			}
		}
		if !success {
			hint := ""
			if cnt := len(tokens); cnt > 0 {
				hint = " after " + tokens[cnt].value
			}
			return nil, fmt.Errorf("unable to lex token%s at %d:%d", hint, cur.loc.line, cur.loc.col)
		}
	}
	return tokens, nil
}

//digits
//digits.[digits][e[+-]digits]
//[digits].digits[e[+-]digits]
//digitse[+-]digits
func lexNumeric(source string, init cursor) (*token, cursor, bool) {
	cur := init
	var foundPeriod bool
	var foundExp bool
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		cur.loc.col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExp := c == 'e'

		if cur.pointer == init.pointer {
			if !isDigit && !isPeriod {
				return nil, init, false
			}
			foundPeriod = isPeriod
			continue
		}
		if isPeriod {
			if foundPeriod {
				return nil, init, false
			}
			foundPeriod = true
			continue
		}
		if isExp {
			if foundExp {
				return nil, init, false
			}
			foundPeriod = true
			foundExp = true
			if cur.pointer == uint(len(source)-1) {
				return nil, init, false
			}
			nextC := source[cur.pointer+1]
			if nextC == '-' || nextC == '+' {
				cur.pointer++
				cur.loc.col++
			}
			continue
		}
		if !isDigit {
			break
		}
	}

	if cur.pointer == init.pointer {
		return nil, init, false
	}
	return &token{
		value: source[init.pointer:cur.pointer],
		loc:   init.loc,
		kind:  kindNumeric,
	}, cur, true
}

func lexCharDelimited(source string, init cursor, delimiter byte) (*token, cursor, bool) {
	cur := init
	// end of source
	if init.pointer >= uint(len(source)) {
		return nil, init, false
	}
	// first char is not delimiter
	if source[init.pointer] != delimiter {
		return nil, init, false
	}

	cur.loc.col++
	cur.pointer++
	var value []byte
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		if c == delimiter {
			// no char after delimiter or next char is not a delimiter
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				return &token{
					value: string(value),
					loc:   init.loc,
					kind:  kindString,
				}, cur, true
			} else {
				value = append(value, delimiter)
				cur.pointer++
				cur.loc.col++
			}
		}
		value = append(value, c)
		cur.loc.col++
	}
	return nil, init, false
}

func lexString(source string, init cursor) (*token, cursor, bool) {
	return lexCharDelimited(source, init, '\'')
}

func lexSymbol(source string, init cursor) (*token, cursor, bool) {
	c := source[init.pointer]
	cur := init
	cur.pointer++
	cur.loc.col++

	switch c {
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	symbols := []symbol{
		symComma,
		symLeftParen,
		symRightParen,
		symAsterisk,
		symSemicolon,
	}
	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	match := longestMatch(source, init, options)
	if match == "" {
		return nil, init, false
	}
	cur.pointer = init.pointer + uint(len(match))
	cur.loc.col = init.loc.col + uint(len(match))

	return &token{
		value: match,
		loc:   init.loc,
		kind:  kindSymbol,
	}, cur, true
}

func longestMatch(source string, init cursor, options []string) string {
	return ""
}