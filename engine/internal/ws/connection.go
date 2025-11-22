package ws

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/paulhalleux/workflow-engine-go/proto"
	gproto "google.golang.org/protobuf/proto"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Connection struct {
	id         string
	ws         *websocket.Conn
	send       chan *proto.WebsocketMessage
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	registry   *Registry
	scopesMu   sync.RWMutex
	scopes     map[proto.WebsocketScopeType][]string
}

func NewConnection(id string, ws *websocket.Conn, reg *Registry) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	log.Printf("new websocket connection established: %s", id)
	return &Connection{
		id:         id,
		ws:         ws,
		send:       make(chan *proto.WebsocketMessage, 256),
		ctx:        ctx,
		cancelFunc: cancel,
		registry:   reg,
		scopes:     make(map[proto.WebsocketScopeType][]string),
	}
}

func (c *Connection) Start() {
	c.wg.Add(2)
	go c.readPump()
	go c.writePump()
}

func (c *Connection) Stop() {
	c.cancelFunc()
	c.wg.Wait()
	close(c.send)
	err := c.ws.Close()
	if err != nil {
		log.Printf("error closing websocket connection: %v", err)
	}
}

func (c *Connection) SendMessage(msg *proto.WebsocketMessage) error {
	select {
	case c.send <- msg:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *Connection) ResetScopes() {
	c.scopes = make(map[proto.WebsocketScopeType][]string)
}

func (c *Connection) AddScope(scopeType proto.WebsocketScopeType, scopeID *string) {
	c.scopesMu.Lock()
	defer c.scopesMu.Unlock()

	if _, exists := c.scopes[scopeType]; !exists {
		c.scopes[scopeType] = []string{}
	}

	if scopeID != nil {
		c.scopes[scopeType] = append(c.scopes[scopeType], *scopeID)
	}
}

func (c *Connection) IsSubscribedTo(scopeType proto.WebsocketScopeType, scopeID *string) bool {
	c.scopesMu.RLock()
	defer c.scopesMu.RUnlock()

	ids, exists := c.scopes[scopeType]
	if !exists {
		return false
	}

	if scopeID == nil {
		return true
	}

	for _, id := range ids {
		if id == *scopeID {
			return true
		}
	}

	return false
}

func (c *Connection) readPump() {
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
			mt, data, err := c.ws.ReadMessage()
			if err != nil {
				log.Printf("error reading websocket message: %v", err)
				return
			}

			if mt != websocket.BinaryMessage {
				continue
			}

			var command proto.WebsocketCommand
			err = gproto.Unmarshal(data, &command)
			if err != nil {
				log.Printf("error unmarshaling websocket command: %v", err)
				continue
			}

			// Handle incoming commands
			if command.Command != nil {
				handler, exists := c.registry.GetCommandHandler(command.Type)
				if exists {
					err := handler.Handle(c.ctx, c, &command)
					if err != nil {
						log.Printf("error handling command: %v", err)
					}
				} else {
					log.Printf("no handler registered for command type: %v", command.Type)
				}
			}

			// For demonstration, just log the received command
			log.Printf("received command: %v", &command)
		}
	}
}

func (c *Connection) writePump() {
	defer c.wg.Done()
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := gproto.Marshal(message)
			if err != nil {
				log.Printf("error marshaling websocket message: %v", err)
				continue
			}
			err = c.ws.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Printf("error writing websocket message: %v", err)
				return
			}

			// For demonstration, just log the sent message
			log.Printf("sent message: %v", message)
		case <-ticker.C:
			_ = c.ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
		case <-c.ctx.Done():
			return
		}
	}
}
