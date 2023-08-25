package gen

import (
	"fmt"
	"go/format"
)

func formatGo(in string) (string, error) {
	formatted, err := format.Source([]byte(in))
	if err != nil {
		return in, fmt.Errorf("failed formatting Go code: %w", err)
	}
	return string(formatted), nil
}