package command

import (
	"redis-go/internal/protocol"
	"redis-go/internal/storage"
	"strings"
)

type Handler func(args []string, store *storage.Store, writer *protocol.Writer) error

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: map[string]Handler{
			"PING":   Ping,
			"GET":    Get,
			"SET":    Set,
			"DEL":    Del,
			"EXISTS": Exists,
			"INCR":   Incr,
			"DECR":   Decr,
			"EXPIRE": Expire,
			"TTL":    TTL,
		},
	}
}

func (r *Registry) Execute(args []string, store *storage.Store, writer *protocol.Writer) error {
	if len(args) == 0 {
		return writer.Error("ERR empty command")
	}

	name := strings.ToUpper(args[0])

	handler, exists := r.handlers[name]

	if !exists {
		return writer.Error("ERR unknown command '" + strings.ToLower(name) + "'")
	}

	return handler(args[1:], store, writer)
}
