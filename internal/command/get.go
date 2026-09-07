package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func Get(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) != 1 {
		return writer.Error("ERR wrong number of arguments for 'get' command")
	}

	value, exists := store.Get(args[0])
	if !exists {
		return writer.Null()
	}

	return writer.BulkString(value)
}
