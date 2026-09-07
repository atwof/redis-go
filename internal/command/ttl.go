package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func TTL(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) != 1 {
		return writer.Error("ERR wrong number of arguments for 'ttl' command")
	}

	return writer.Integer(store.TTL(args[0]))
}
