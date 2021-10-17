package table

import (
	"fmt"
	"strings"
)

type Result struct {
	Columns []Column
	Rows    [][]Cell
}

func (r Result) String() string {
	sb := new(strings.Builder)
	for _, col := range r.Columns {
		sb.WriteString(fmt.Sprintf("| %s ", col.Name))
	}
	sb.WriteString("|\n")

	width := sb.Len()
	for i := 0; i < width; i++ {
		sb.WriteString("=")
	}

	for _, result := range r.Rows {
		sb.WriteString("\n")
		sb.WriteString("|")
		for i, cell := range result {
			typ := r.Columns[i].Type
			s := ""
			switch typ {
			case Int:
				s = fmt.Sprintf("%d", cell.AsInt())
			case Text:
				s = cell.AsText()
			}
			sb.WriteString(fmt.Sprintf(" %s | ", s))
		}
	}
	return sb.String()
}
