package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func Exists(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) == 0 {
		return writer.Error("ERR wrong number of arguments for 'exists' command")
	}

	var count int64

	for _, key := range args {
		if store.Exists(key) {
			count++
		}
	}

	return writer.Integer(count)
}
