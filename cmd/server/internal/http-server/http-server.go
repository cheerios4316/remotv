package httpserver

import (
	"fmt"
	"net/http"
	"remotv/cmd/server/internal/http-server/handlers"
	tcpserver "remotv/cmd/server/internal/tcp-server"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Port      int
	TcpServer *tcpserver.Server
}

func (s *Server) Listen() error {
	router := chi.NewRouter()
	router.Post("/remote/{clientId}", handlers.HandlePostCommand(s.TcpServer))
	router.Get("/remote/clients", handlers.HandleGetClients(s.TcpServer))

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), router)
}
