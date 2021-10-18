package table

import (
	"fmt"
	"strings"
)

type Column struct {
	Name string
	Type DataType
}

func NewColumn(name string, dt string) (c *Column) {
	c = new(Column)
	c.Name = name
	switch strings.ToLower(dt) {
	case "int":
		c.Type = Int
	case "text":
		c.Type = Text
	}
	return
}

type Table struct {
	Columns []*Column
	Rows    [][]Cell
}

func NewTable(cols []*Column) (tb *Table) {
	tb = &Table{
		Columns: cols,
		Rows:    make([][]Cell, 0),
	}
	return
}

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
