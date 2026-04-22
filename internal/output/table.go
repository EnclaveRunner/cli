package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"cli/internal/styles"
)

type tablePrinter struct {
	columns []Column
	w       io.Writer
}

func (p *tablePrinter) Print(rows any) error {
	items := toSlice(rows)
	if len(items) == 0 {
		fmt.Fprintln(p.w, styles.MutedStyle.Render("No results."))
		return nil
	}

	// Compute column widths: max of header length, MinWidth, and all cell values.
	widths := make([]int, len(p.columns))
	cells := make([][]string, len(items))

	for i, col := range p.columns {
		w := len(col.Header)
		if col.MinWidth > w {
			w = col.MinWidth
		}
		widths[i] = w
	}

	for r, row := range items {
		cells[r] = make([]string, len(p.columns))
		for c, col := range p.columns {
			val := col.Extract(row)
			// Strip ANSI for width calculation.
			plain := stripAnsi(val)
			if len(plain) > widths[c] {
				widths[c] = len(plain)
			}
			cells[r][c] = val
		}
	}

	// Render header.
	headerCells := make([]string, len(p.columns))
	for i, col := range p.columns {
		padded := pad(col.Header, widths[i])
		headerCells[i] = styles.HeaderStyle.Render(padded)
	}
	fmt.Fprintln(p.w, strings.Join(headerCells, ""))

	// Render rows.
	for _, row := range cells {
		rowCells := make([]string, len(p.columns))
		for i, cell := range row {
			plain := stripAnsi(cell)
			// Pad with plain spaces so the column aligns, then wrap with
			// a single-space margin on each side (no lipgloss padding, which
			// would mis-count width when cell already contains ANSI codes).
			padding := widths[i] - len([]rune(plain))
			if padding < 0 {
				padding = 0
			}
			rowCells[i] = " " + cell + strings.Repeat(" ", padding) + " "
		}
		fmt.Fprintln(p.w, strings.Join(rowCells, ""))
	}

	return nil
}

// toSlice converts any slice value to []any using reflection.
func toSlice(v any) []any {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return nil
	}
	out := make([]any, rv.Len())
	for i := range rv.Len() {
		out[i] = rv.Index(i).Interface()
	}
	return out
}

// stripAnsi removes ANSI escape codes for width measurement.
func stripAnsi(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
