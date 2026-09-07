package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func Del(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) == 0 {
		return writer.Error("ERR wrong number of arguments for 'del' command")
	}

	var deleted int64

	for _, key := range args {
		if store.Delete(key) {
			deleted++
		}
	}

	return writer.Integer(deleted)
}
