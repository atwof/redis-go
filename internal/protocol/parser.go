package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

func (p *Parser) ReadCommand() ([]string, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	if len(line) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	// RESP Array
	if line[0] == '*' {
		return p.readArray(line)
	}

	// Inline command
	return splitInlineCommand(line), nil
}

func (p *Parser) readArray(header string) ([]string, error) {
	count, err := strconv.Atoi(header[1:])
	if err != nil {
		return nil, fmt.Errorf(
			"invalid array length: %q",
			header,
		)
	}

	args := make([]string, 0, count)

	for i := 0; i < count; i++ {
		header, err := p.readLine()
		if err != nil {
			return nil, err
		}

		if len(header) == 0 || header[0] != '$' {
			return nil, fmt.Errorf(
				"expected bulk string, got %q",
				header,
			)
		}

		length, err := strconv.Atoi(header[1:])
		if err != nil {
			return nil, fmt.Errorf(
				"invalid bulk string length: %q",
				header,
			)
		}

		data := make([]byte, length)

		if _, err := io.ReadFull(p.reader, data); err != nil {
			return nil, err
		}

		// Consumir CRLF depois do conteúdo.
		crlf := make([]byte, 2)

		if _, err := io.ReadFull(p.reader, crlf); err != nil {
			return nil, err
		}

		if crlf[0] != '\r' || crlf[1] != '\n' {
			return nil, fmt.Errorf(
				"expected CRLF",
			)
		}

		args = append(args, string(data))
	}

	return args, nil
}

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(line) < 2 ||
		line[len(line)-2:] != "\r\n" {
		return "", fmt.Errorf(
			"invalid RESP line: %q",
			line,
		)
	}

	return line[:len(line)-2], nil
}

func splitInlineCommand(line string) []string {
	var result []string

	var current []byte
	inQuotes := false

	for _, char := range line {
		switch char {
		case '"':
			inQuotes = !inQuotes

		case ' ', '\t':
			if !inQuotes && len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}

		default:
			current = append(current, byte(char))
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}
