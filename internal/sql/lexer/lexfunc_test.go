package lexer

import (
	"mintsql/internal/sql/token"
	"testing"
)

func TestLexNumeric(t *testing.T) {
	succKinds := []token.Kind{token.KindNumeric, token.KindEof}
	failKinds := []token.Kind{token.KindError}
	tests := []struct {
		input        string
		expectValues []string
		expectKinds  []token.Kind
	}{
		{
			input:        "123",
			expectValues: []string{"123", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123.",
			expectValues: []string{"123.", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123.4",
			expectValues: []string{"123.4", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123e4",
			expectValues: []string{"123e4", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123.e+4",
			expectValues: []string{"123.e+4", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123.4e5",
			expectValues: []string{"123.4e5", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123.4e+5",
			expectValues: []string{"123.4e+5", ""},
			expectKinds:  succKinds,
		},
		{
			input:        ".123",
			expectValues: []string{".123", ""},
			expectKinds:  succKinds,
		},
		{
			input:        ".123e4",
			expectValues: []string{".123e4", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "abc",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
		{
			input:        "123.4e",
			expectValues: []string{"123.4e"},
			expectKinds:  failKinds,
		},
		{
			input:        ".e-3",
			expectValues: []string{"."},
			expectKinds:  failKinds,
		},
		{
			input:        "123.e+",
			expectValues: []string{"123.e+"},
			expectKinds:  failKinds,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexNumeric)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValues[i] || tk.Kind != test.expectKinds[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValues[i], test.expectKinds[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}

func TestLexString(t *testing.T) {
	succKinds := []token.Kind{token.KindString, token.KindEof}
	failKinds := []token.Kind{token.KindError}
	tests := []struct {
		input        string
		expectValues []string
		expectKinds  []token.Kind
	}{
		{
			input:        "'abc'",
			expectValues: []string{"abc", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "\"abc\"",
			expectValues: []string{"abc", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "''",
			expectValues: []string{"", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "\"\"",
			expectValues: []string{"", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "\"ab\"\"c\"",
			expectValues: []string{"ab\"\"c", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "\"ab\\'c\"",
			expectValues: []string{"ab\\'c", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "'",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
		{
			input:        "\"",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
		{
			input:        "abc",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexString)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValues[i] || tk.Kind != test.expectKinds[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValues[i], test.expectKinds[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}

func TestLexKeyword(t *testing.T) {
	succKinds := []token.Kind{token.KindKeyword, token.KindEof}
	failKinds := []token.Kind{token.KindError}
	tests := []struct {
		input        string
		expectValues []string
		expectKinds  []token.Kind
	}{
		{
			input:        "SELECT",
			expectValues: []string{"SELECT", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "from",
			expectValues: []string{"from", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "INSERT",
			expectValues: []string{"INSERT", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "INTO",
			expectValues: []string{"INTO", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "INT",
			expectValues: []string{"INT", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "TAB",
			expectValues: []string{"TAB"},
			expectKinds:  failKinds,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexKeyword)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValues[i] || tk.Kind != test.expectKinds[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValues[i], test.expectKinds[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}

func TestLexSymbol(t *testing.T) {
	succKinds := []token.Kind{token.KindSymbol, token.KindEof}
	failKinds := []token.Kind{token.KindError}
	tests := []struct {
		input        string
		expectValues []string
		expectKinds  []token.Kind
	}{
		{
			input:        "*",
			expectValues: []string{"*", ""},
			expectKinds:  succKinds,
		},
		{
			input:        ")",
			expectValues: []string{")", ""},
			expectKinds:  succKinds,
		},
		{
			input:        ";",
			expectValues: []string{";", ""},
			expectKinds:  succKinds,
		},
		{
			input:        ".",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexSymbol)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValues[i] || tk.Kind != test.expectKinds[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValues[i], test.expectKinds[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}

func TestLexIdentifier(t *testing.T) {
	succKinds := []token.Kind{token.KindIdentifier, token.KindEof}
	failKinds := []token.Kind{token.KindError}
	tests := []struct {
		input        string
		expectValues []string
		expectKinds  []token.Kind
	}{
		{
			input:        "_abc",
			expectValues: []string{"_abc", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "abc",
			expectValues: []string{"abc", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "__",
			expectValues: []string{"__", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "a1b2c3",
			expectValues: []string{"a1b2c3", ""},
			expectKinds:  succKinds,
		},
		{
			input:        "123a",
			expectValues: []string{""},
			expectKinds:  failKinds,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexIdentifier)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValues[i] || tk.Kind != test.expectKinds[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValues[i], test.expectKinds[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}

func TestLexBegin(t *testing.T) {
	tests := []struct {
		input       string
		expectValue []string
		expectKind  []token.Kind
	}{
		{
			input:       "SELECT name, id FROM users;",
			expectValue: []string{"SELECT", "name", ",", "id", "FROM", "users", ";", ""},
			expectKind: []token.Kind{
				token.KindKeyword,
				token.KindIdentifier,
				token.KindSymbol,
				token.KindIdentifier,
				token.KindKeyword,
				token.KindIdentifier,
				token.KindSymbol,
				token.KindEof,
			},
		},
		{
			input:       "INSERT INTO users VALUES (2, 'Kate');",
			expectValue: []string{"INSERT", "INTO", "users", "VALUES", "(", "2", ",", "Kate", ")", ";", ""},
			expectKind: []token.Kind{
				token.KindKeyword,
				token.KindKeyword,
				token.KindIdentifier,
				token.KindKeyword,
				token.KindSymbol,
				token.KindNumeric,
				token.KindSymbol,
				token.KindString,
				token.KindSymbol,
				token.KindSymbol,
				token.KindEof,
			},
		},
		{
			input:       "CREATE TABLE users (id INT, name TEXT);",
			expectValue: []string{"CREATE", "TABLE", "users", "(", "id", "INT", ",", "name", "TEXT", ")", ";", ""},
			expectKind: []token.Kind{
				token.KindKeyword,
				token.KindKeyword,
				token.KindIdentifier,
				token.KindSymbol,
				token.KindIdentifier,
				token.KindKeyword,
				token.KindSymbol,
				token.KindIdentifier,
				token.KindKeyword,
				token.KindSymbol,
				token.KindSymbol,
				token.KindEof,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := withState(test.input, lexBegin)
			go lexer.Lex()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				if tk.Value != test.expectValue[i] || tk.Kind != test.expectKind[i] {
					t.Errorf(
						"expected '%s' of kind %d, got '%v' of kind %d",
						test.expectValue[i], test.expectKind[i], tk.Value, tk.Kind,
					)
				}
				tk = lexer.NextToken()
			}
		})
	}
}
