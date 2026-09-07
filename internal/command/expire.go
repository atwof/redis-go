package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
	"strconv"
	"time"
)

func Expire(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) != 2 {
		return writer.Error("ERR wrong number of arguments for 'expire' command")
	}

	seconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return writer.Error("ERR invalid expire time")
	}

	if store.Expire(args[0], time.Duration(seconds)*time.Second) {
		return writer.Integer(1)
	}

	return writer.Integer(0)
}
