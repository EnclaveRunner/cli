package output

import "io"

// Format represents the output rendering format.
type Format int

const (
	FormatTable Format = iota
	FormatJSON
	FormatYAML
)

// ParseFormat converts a string to a Format. Defaults to FormatTable.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatTable
	}
}

// Column describes one column in table output.
type Column struct {
	Header  string
	Extract func(row any) string
	// MinWidth is the minimum column width. 0 means use header length.
	MinWidth int
}

// Printer renders resource slices to an io.Writer.
type Printer interface {
	Print(rows any) error
}

// New returns the appropriate Printer for the requested format.
func New(format Format, columns []Column, w io.Writer) Printer {
	switch format {
	case FormatJSON:
		return &jsonPrinter{w: w}
	case FormatYAML:
		return &yamlPrinter{w: w}
	case FormatTable:
		return &tablePrinter{columns: columns, w: w}
	default:
		return &tablePrinter{columns: columns, w: w}
	}
}
