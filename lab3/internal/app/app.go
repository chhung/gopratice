package app

import (
	"bufio"
	"fmt"
	"io"

	"lab3/internal/service"
)

func Run(input io.Reader, output io.Writer) error {
	if _, err := fmt.Fprint(output, "Input: "); err != nil {
		return err
	}

	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	if _, err := fmt.Fprintf(output, "Normalized Input: %s\n", service.NormalizeInput(line)); err != nil {
		return err
	}

	return nil
}
