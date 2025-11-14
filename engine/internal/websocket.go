package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WebsocketHub struct {
	clients    map[*WebsocketClient]bool
	register   chan *WebsocketClient
	unregister chan *WebsocketClient
	broadcast  chan *proto.WebsocketMessage
	incoming   chan *proto.WebsocketMessage
}

type WebsocketClient struct {
	Scopes map[proto.WebsocketScopeType][]string
	Conn   *websocket.Conn
	send   chan *proto.WebsocketMessage
	hub    *WebsocketHub
}

func NewWebsocketClient(conn *websocket.Conn, hub *WebsocketHub) *WebsocketClient {
	return &WebsocketClient{
		Conn:   conn,
		send:   make(chan *proto.WebsocketMessage, 256),
		hub:    hub,
		Scopes: make(map[proto.WebsocketScopeType][]string),
	}
}

func NewWebsocketHub() *WebsocketHub {
	return &WebsocketHub{
		clients:    make(map[*WebsocketClient]bool),
		register:   make(chan *WebsocketClient, 256),
		unregister: make(chan *WebsocketClient, 256),
		broadcast:  make(chan *proto.WebsocketMessage, 256),
	}
}

func (h *WebsocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			h.HandleMessageBroadcast(message)
		case message := <-h.incoming:
			h.HandleIncomingMessage(message)
		}
	}
}

func (h *WebsocketHub) RegisterNewClient(client *WebsocketClient) {
	h.register <- client
}

func (h *WebsocketHub) RemoveClient(client *WebsocketClient) {
	h.unregister <- client
}

func (h *WebsocketHub) HandleMessageBroadcast(message *proto.WebsocketMessage) {
	log.Printf("Sending message to clients: %v", message)
}

func (h *WebsocketHub) HandleIncomingMessage(message *proto.WebsocketMessage) {
	log.Printf("Handling incoming message from client: %v", message)
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
		var msg proto.WebsocketMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		c.hub.broadcast <- &msg
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
				err := c.Conn.WriteJSON(message)
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
