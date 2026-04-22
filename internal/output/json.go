package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonPrinter struct {
	w io.Writer
}

func (p *jsonPrinter) Print(rows any) error {
	enc := json.NewEncoder(p.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rows); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
