package output

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type yamlPrinter struct {
	w io.Writer
}

func (p *yamlPrinter) Print(rows any) error {
	enc := yaml.NewEncoder(p.w)
	enc.SetIndent(2)
	if err := enc.Encode(rows); err != nil {
		return fmt.Errorf("encode yaml: %w", err)
	}

	return nil
}
