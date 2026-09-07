package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
	"strconv"
	"strings"
	"time"
)

func Set(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) < 2 {
		return writer.Error("ERR wrong number of arguments for 'set' command")
	}

	key := args[0]
	value := args[1]

	if len(args) == 2 {
		store.Set(key, value)
		return writer.SimpleString("OK")
	}

	if len(args) == 4 && strings.ToUpper(args[2]) == "EX" {
		seconds, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil || seconds <= 0 {
			return writer.Error("ERR invalid expire time")
		}

		store.SetWithExpiration(key, value, time.Duration(seconds)*time.Second)

		return writer.SimpleString("OK")
	}

	return writer.Error("ERR syntax error")
}
