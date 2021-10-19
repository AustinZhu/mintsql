package database

import "fmt"

const (
	DuplicateTableError = iota + 40
	NoSuchTableError
	NoSuchColumnError
	NoSuchTypeError
	IllegalValueError
	TypeCheckError
	UnknownError
)

type StoreError struct {
	Code      int
	tableName string
	hint      string
}

func (e *StoreError) Error() string {
	switch e.Code {
	case DuplicateTableError:
		return e.fmtStoreError("table '%s' already exists")
	case NoSuchTableError:
		return e.fmtStoreError("table '%s' does not exist")
	case IllegalValueError:
		return e.fmtStoreError("illegal value '%s'")
	case NoSuchTypeError:
		return e.fmtStoreError("unclassified type '%s'")
	case TypeCheckError:
		return e.fmtStoreError("'%s' type mismatch")
	case NoSuchColumnError:
		return e.fmtStoreError("column '%s' does not exist")
	case UnknownError:
		fallthrough
	default:
		return "unknown error"
	}
}

func (e *StoreError) fmtStoreError(tmpl string) string {
	matter := fmt.Sprintf(tmpl, e.hint)
	return fmt.Sprintf("store error @%s: %s", e.tableName, matter)
}

func Error(code int, table string, hint string) *StoreError {
	return &StoreError{
		Code:      code,
		tableName: table,
		hint:      hint,
	}
}
