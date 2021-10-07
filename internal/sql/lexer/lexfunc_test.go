package lexer

import (
	"fmt"
	"mintsql/internal/sql/token"
	"testing"
)

func TestLexDigits(t *testing.T) {
	tests := []struct {
		input       string
		expectValue []string
		expectKind  []token.Kind
	}{
		{
			input:       "123",
			expectValue: []string{"123", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123.",
			expectValue: []string{"123.", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123.4",
			expectValue: []string{"123.4", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123e4",
			expectValue: []string{"123e4", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123.e+4",
			expectValue: []string{"123.e+4", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123.4e5",
			expectValue: []string{"123.4e5", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "123.4e+5",
			expectValue: []string{"123.4e+5", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       ".123",
			expectValue: []string{".123", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       ".123e4",
			expectValue: []string{".123e4", "EOF"},
			expectKind:  []token.Kind{token.KindNumeric, token.KindEof},
		},
		{
			input:       "abc",
			expectValue: []string{""},
			expectKind:  []token.Kind{token.KindError},
		},
		{
			input:       "123.4e",
			expectValue: []string{"123.4e"},
			expectKind:  []token.Kind{token.KindError},
		},
		{
			input:       "123.4.e",
			expectValue: []string{"123.4", "."},
			expectKind:  []token.Kind{token.KindNumeric, token.KindError},
		},
		{
			input:       ".e-3",
			expectValue: []string{"."},
			expectKind:  []token.Kind{token.KindError},
		},
		{
			input:       "123.e+",
			expectValue: []string{"123.e+"},
			expectKind:  []token.Kind{token.KindError},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := New(test.input, LexNumeric)
			go lexer.Run()
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

func TestLexBegin(t *testing.T) {
	tests := []struct {
		input       string
		expectValue []string
		expectKind  []token.Kind
	}{
		{
			input:       "SELECT name, id FROM users;",
			expectValue: []string{"SELECT", "name", ",", "id", "FROM", "users", ";", "EOF"},
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
			expectValue: []string{"INSERT", "INTO", "users", "VALUES", "(", "2", ",", "Kate", ")", ";", "EOF"},
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
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := New(test.input, LexBegin)
			go lexer.Run()
			tk := lexer.NextToken()
			for i := 0; tk != nil; i++ {
				fmt.Println(tk)
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
