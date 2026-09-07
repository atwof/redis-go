package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
)

func Ping(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) == 0 {
		return writer.SimpleString("PONG")
	}

	return writer.BulkString(args[0])
}
