package parser

import (
	"fmt"
	"mintsql/internal/sql/token"
)

const (
	ParseStmtError = iota + 20
	ParseExprError
	ParseUnknownError
)

func Error(t *token.Token, msg string, expect ...interface{}) error {
	if t == nil {
		return fmt.Errorf("error: missing token, expected %s", expect)
	}
	return fmt.Errorf("%s: error: %s, expected %s, got '%s'", t.Location, msg, expect, t.Value)
}
