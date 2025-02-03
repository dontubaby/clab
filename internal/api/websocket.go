package server

import (
	//"encoding/json"
	//"fmt"
	"log"
	"net/http"

	//"cyber/internal/models"

	"github.com/lxzan/gws"
)

type WebSocketServer struct {
	server *gws.Server
}

func NewWebsocketServer() *WebSocketServer {
	return &WebSocketServer{
		server: gws.NewServer(NewWebSocketHandler(), nil),
	}
}

func (s *WebSocketServer) Start(addr string) {
	log.Printf("WebSocket server starting at: %s...\n", addr)
	if err := http.ListenAndServe(addr, s.Server); err != nil {
		log.Fatalf("WebSocket server failed: %v", err)
	}
}
