package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketServer interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Upgrader websocket.Upgrader
	Registry *Registry
}

func NewServer() *Server {
	return &Server{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		Registry: NewRegistry(),
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	conn := NewConnection(r.RemoteAddr, wsConn, s.Registry)
	conn.Start()
}
