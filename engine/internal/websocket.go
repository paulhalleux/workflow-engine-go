package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/paulhalleux/workflow-engine-go/proto"
	gproto "google.golang.org/protobuf/proto"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WebsocketHub struct {
	clients    map[string]*WebsocketClient
	register   chan *WebsocketClient
	unregister chan string
	broadcast  chan *proto.WebsocketMessage
	incoming   chan *proto.WebsocketCommand
}

type WebsocketClient struct {
	ID     string
	Scopes map[proto.WebsocketScopeType][]string
	Conn   *websocket.Conn
	send   chan *proto.WebsocketMessage
	hub    *WebsocketHub
}

func NewWebsocketClient(conn *websocket.Conn, hub *WebsocketHub) *WebsocketClient {
	return &WebsocketClient{
		ID:     uuid.New().String(),
		Conn:   conn,
		send:   make(chan *proto.WebsocketMessage, 256),
		hub:    hub,
		Scopes: make(map[proto.WebsocketScopeType][]string),
	}
}

func NewWebsocketHub() *WebsocketHub {
	return &WebsocketHub{
		clients:    make(map[string]*WebsocketClient),
		register:   make(chan *WebsocketClient, 256),
		unregister: make(chan string, 256),
		broadcast:  make(chan *proto.WebsocketMessage, 256),
		incoming:   make(chan *proto.WebsocketCommand, 256),
	}
}

func (h *WebsocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			h.handleMessageBroadcast(&proto.WebsocketMessage{
				Type: proto.WebsocketMessageType_WEBSOCKET_MESSAGE_TYPE_REGISTERED,
				Payload: &proto.WebsocketMessage_RegisteredMessage{
					RegisteredMessage: &proto.RegisteredMessage{
						ClientId: client.ID,
					},
				},
			})
		case clientId := <-h.unregister:
			if client, ok := h.clients[clientId]; ok {
				delete(h.clients, clientId)
				close(client.send)
			}
		case message := <-h.broadcast:
			h.handleMessageBroadcast(message)
		case command := <-h.incoming:
			h.handleIncomingCommand(command)
		}
	}
}

func (h *WebsocketHub) RegisterNewClient(client *WebsocketClient) {
	h.register <- client
}

func (h *WebsocketHub) RemoveClient(client *WebsocketClient) {
	h.unregister <- client.ID
}

func (h *WebsocketHub) BroadcastMessage(message *proto.WebsocketMessage) {
	h.broadcast <- message
}

func (h *WebsocketHub) handleMessageBroadcast(message *proto.WebsocketMessage) {
	log.Printf("Sending message to clients: %v", message)
	for clientId, client := range h.clients {
		if !IsInScope(client.Scopes, message.Scope) {
			log.Printf("Skipping client %s for message %v due to scope mismatch", clientId, message)
			continue
		}

		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, clientId)
		}
	}
}

func (h *WebsocketHub) handleIncomingCommand(command *proto.WebsocketCommand) {
	log.Printf("Handling incoming command from client: %v", command)
	// Handle incoming messages from clients if needed
	if command.Type == proto.WebsocketCommandType_WEBSOCKET_COMMAND_TYPE_SUBSCRIBE {
		log.Printf("Subscribing to clients: %v", command)
		clientScopes := command.GetSubscribeCommand().GetScopes()
		log.Printf("Client scopes: %v", clientScopes)
		client, ok := h.clients[command.ClientId]
		if !ok {
			return
		}

		client.Scopes = make(map[proto.WebsocketScopeType][]string)
		for _, scope := range clientScopes {
			_, ok = client.Scopes[scope.Type]
			if !ok {
				client.Scopes[scope.Type] = make([]string, 0)
			}

			if scope.Id != nil {
				client.Scopes[scope.Type] = append(client.Scopes[scope.Type], *scope.Id)
			}
		}
	}
}

func (c *WebsocketClient) Read() {
	defer func() {
		c.hub.RemoveClient(c)
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading websocket message:", err)
			break
		}

		if mt != websocket.BinaryMessage {
			continue
		}

		var cmd proto.WebsocketCommand
		if err := gproto.Unmarshal(data, &cmd); err != nil {
			fmt.Println("Error unmarshaling protobuf command:", err)
			continue
		}

		c.hub.incoming <- &cmd
	}
}

func (c *WebsocketClient) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				data, err := gproto.Marshal(message)
				if err != nil {
					fmt.Println("Error: ", err)
					continue
				}
				err = c.Conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WebsocketClient) Close() {
	close(c.send)
}

func IsInScope(clientScopes map[proto.WebsocketScopeType][]string, messageScope *proto.WebsocketScope) bool {
	if messageScope == nil {
		return true
	}

	log.Printf("Checking scope: clientScopes=%v, messageScope=%v", clientScopes, messageScope)
	allowedValues, ok := clientScopes[messageScope.Type]
	if !ok {
		return false
	}

	if messageScope.Id == nil {
		return true
	}

	for _, v := range allowedValues {
		if v == *messageScope.Id {
			return true
		}
	}

	return false
}
