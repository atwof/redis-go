package protocol

import (
	"bufio"
	"fmt"
)

type Writer struct {
	writer *bufio.Writer
}

func NewWriter(writer *bufio.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) SimpleString(value string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", value)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *Writer) Error(value string) error {
	_, err := fmt.Fprintf(w.writer, "-%s\r\n", value)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *Writer) Integer(value int64) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", value)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *Writer) BulkString(value string) error {
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(value), value)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *Writer) Null() error {
	_, err := fmt.Fprintf(w.writer, "$-1\r\n")
	if err != nil {
		return err
	}

	return w.writer.Flush()
}
