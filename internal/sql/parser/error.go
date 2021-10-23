package parser

import (
	"fmt"
	"mintsql/internal/sql/token"
)

const (
	InvalidStmtError = iota + 10
	InvalidNameError
	InvalidValueError
	InvalidTypeError
	MissingSymbolError
	MissingKeywordError
	ExcessiveTokenError
	ParseUnknownError
)

type SyntaxError struct {
	Code   int
	tk     *token.Token
	expect []fmt.Stringer
}

func (e *SyntaxError) Error() string {
	switch e.Code {
	case InvalidStmtError:
		return e.fmtSyntaxError("invalid statement")
	case InvalidNameError:
		return e.fmtSyntaxError("invalid name")
	case InvalidValueError:
		return e.fmtSyntaxError("invalid value")
	case InvalidTypeError:
		return e.fmtSyntaxError("invalid data type")
	case MissingSymbolError:
		return e.fmtSyntaxError("missing symbol")
	case MissingKeywordError:
		return e.fmtSyntaxError("missing keyword")
	case ExcessiveTokenError:
		return e.fmtSyntaxError("excessive token")
	case ParseUnknownError:
		fallthrough
	default:
		return "unknown error"
	}
}

func (e *SyntaxError) fmtSyntaxError(msg string) string {
	exp := fmtExpected(e.expect)
	if e.tk == nil {
		return fmt.Sprintf("syntax error: missing token%s", exp)
	}
	return fmt.Sprintf("syntax error @%s: %s%s, got %s", e.tk.Location, msg, exp, e.tk)
}

func Error(code int, get *token.Token, expect ...fmt.Stringer) *SyntaxError {
	return &SyntaxError{
		Code:   code,
		tk:     get,
		expect: expect,
	}
}

func fmtExpected(exp []fmt.Stringer) string {
	switch len(exp) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(", expected '%s'", exp[0].String())
	default:
		return fmt.Sprintf(", expected one of %s", exp)
	}
}
