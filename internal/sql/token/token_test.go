package token

import (
	"testing"
)

func TestToken_Equals(t *testing.T) {
	tests := []struct {
		name     string
		tk1      *Token
		tk2      *Token
		expected bool
	}{
		{
			name:     "Symbol",
			tk1:      NewSymbol(SEMICOLON),
			tk2:      NewSymbol(SEMICOLON),
			expected: true,
		},
		{
			name: "Keyword",
			tk1:  NewKeyword(SELECT),
			tk2: &Token{
				Value: "SELECT",
				Kind:  KindKeyword,
			},
			expected: true,
		},
		{
			name:     "Nil",
			tk1:      NewKeyword(SELECT),
			tk2:      nil,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.tk1.Equals(test.tk2) != test.expected {
				t.Error("token not match")
			}
			return
		})
	}
}

func TestIsEnd(t *testing.T) {
	tests := []struct {
		name     string
		tk       *Token
		expected bool
	}{
		{
			name:     "Nil",
			tk:       nil,
			expected: true,
		},
		{
			name: "EOF",
			tk: &Token{
				Value: "",
				Kind:  KindEof,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if IsEnd(test.tk) != test.expected {
				t.Error("not ending token")
			}
			return
		})
	}
}

func TestToken_String(t *testing.T) {
	tests := []struct {
		name     string
		tk       *Token
		expected string
	}{
		{
			name: "EOF",
			tk: &Token{
				Value: "",
				Kind:  KindEof,
			},
			expected: "EOF",
		},
		{
			name: "Error",
			tk: &Token{
				Value: "abd;",
				Kind:  KindError,
			},
			expected: "abd;",
		},
		{
			name:     "Error",
			tk:       NewKeyword(SELECT),
			expected: "'select'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.tk.String() != test.expected {
				t.Error("error stringify token")
			}
			return
		})
	}
}

func TestLocation_String(t *testing.T) {
	loc := &Location{
		Line:   10,
		Column: 20,
	}
	t.Run("Location", func(t *testing.T) {
		if loc.String() != "10:20" {
			t.Error("error stringify location")
		}
		return
	})
}

func TestIsKind(t *testing.T) {

}

func TestNewKeyword(t *testing.T) {

}

func TestNewStream(t *testing.T) {

}

func TestStream_Add(t *testing.T) {

}
