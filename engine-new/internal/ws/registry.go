package ws

import (
	"context"
	"sync"

	"github.com/paulhalleux/workflow-engine-go/proto"
)

type Registry struct {
	mu       sync.RWMutex
	commands map[proto.WebsocketCommandType]CommandHandler
}

type CommandHandler struct {
	Handle func(ctx context.Context, command *proto.WebsocketCommand) error
}

type MessageHandler struct {
	Handle func(ctx context.Context, message *proto.WebsocketMessage) error
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[proto.WebsocketCommandType]CommandHandler),
	}
}

func (r *Registry) RegisterCommand(cmdType proto.WebsocketCommandType, handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands[cmdType] = handler
}

func (r *Registry) GetCommandHandler(cmdType proto.WebsocketCommandType) (CommandHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.commands[cmdType]
	return handler, exists
}
