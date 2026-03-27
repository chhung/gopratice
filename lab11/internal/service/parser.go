package service

import (
	"encoding/xml"
	"io"
	"os"
	"strings"

	"lab11/internal/model"
)

// ParseFile reads a CDC XML file that contains multiple <operation> elements
// with no root wrapper, and returns all parsed operations.
func ParseFile(filePath string) ([]model.Operation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// The files have no root element, so wrap them to produce valid XML.
	wrapped := io.MultiReader(
		strings.NewReader("<root>"),
		f,
		strings.NewReader("</root>"),
	)

	decoder := xml.NewDecoder(wrapped)

	var ops []model.Operation
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		se, ok := tok.(xml.StartElement)
		if !ok || se.Name.Local != "operation" {
			continue
		}

		var op model.Operation
		if err := decoder.DecodeElement(&op, &se); err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	return ops, nil
}
