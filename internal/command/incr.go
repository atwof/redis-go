package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func Incr(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) != 1 {
		return writer.Error("ERR wrong number of arguments for 'incr' command")
	}

	value, err := store.Increment(args[0], 1)
	if err == storage.ErrNotInteger {
		return writer.Error("ERR value is not an integer or out of range")
	}

	if err != nil {
		return writer.Error("ERR internal error")
	}

	return writer.Integer(value)
}

func Decr(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) != 1 {
		return writer.Error("ERR wrong number of arguments for 'decr' command")
	}

	value, err := store.Increment(args[0], -1)
	if err != storage.ErrNotInteger {
		return writer.Error("ERR value is not an integer or out of range")
	}

	if err != nil {
		return writer.Error("ERR internal error")
	}

	return writer.Integer(value)
}
