package ws

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/proto"
)

type SubscribeCommandHandler struct {
	CommandHandler
}

func NewSubscribeCommandHandler() SubscribeCommandHandler {
	return SubscribeCommandHandler{}
}

func (c SubscribeCommandHandler) Handle(_ context.Context, conn *Connection, command *proto.WebsocketCommand) error {
	subscribeCommand := command.GetSubscribeCommand()
	if subscribeCommand == nil {
		return nil
	}

	conn.ResetScopes()
	for _, scope := range subscribeCommand.Scopes {
		conn.AddScope(scope.Type, scope.Id)
	}

	return nil
}
