package app

import (
	"fmt"
	"strings"

	"lab11/internal/model"
	"lab11/internal/service"
)

// Run parses the given XML file and prints each operation to stdout.
func Run(filePath string) error {
	ops, err := service.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %q: %w", filePath, err)
	}

	fmt.Printf("Parsed %d operation(s) from %s\n", len(ops), filePath)
	fmt.Println(strings.Repeat("-", 72))

	for i, op := range ops {
		fmt.Printf("[%d] table=%-40s type=%-20s txInd=%s\n    ts=%s  position=%s\n",
			i+1, op.Table, op.Type, op.TxInd, op.Ts, op.Position)

		for _, col := range op.Columns {
			fmt.Printf("    [%2d] %-25s  before=%-30s  after=%s\n",
				col.Index, col.Name, colVal(col.BeforeValue), colVal(col.AfterValue))
		}
		fmt.Println()
	}

	return nil
}

// colVal returns a display string for a column value element.
func colVal(cv *model.ColValue) string {
	if cv == nil {
		return "-"
	}
	if cv.IsNull == "true" {
		return "NULL"
	}
	v := strings.TrimSpace(cv.Value)
	if v == "" {
		return "(empty)"
	}
	return v
}
