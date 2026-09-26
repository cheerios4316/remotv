package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"remotv/cmd/server/command"
	tcpserver "remotv/cmd/server/internal/tcp-server"

	"github.com/go-chi/chi/v5"
)

type RequestBody interface {
	Validate() error
}

func parseBody[T RequestBody](content io.ReadCloser) (T, error) {
	var body T
	if err := json.NewDecoder(content).Decode(&body); err != nil {
		var zero T
		return zero, err
	}

	return body, nil
}

func jsonResponse(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func jsonError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func HandlePostCommand(tcpServer *tcpserver.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		clientName := chi.URLParam(r, "clientId")

		if clientName == "" {
			http.Error(w, "Empty client name", 400)
			return
		}

		var body, err = parseBody[command.Command](r.Body)
		if err != nil {
			http.Error(w, "Bad request body", http.StatusBadRequest)
			return
		}

		if err := body.Validate(); err != nil {
			jsonError(w, err)
			return
		}

		clients := tcpServer.ConnectedClients()
		client, ok := clients[clientName]
		if !ok {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}

		msg, err := body.ToMessage()
		if err != nil {
			http.Error(w, "Error sending message to client", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(client, msg)

		jsonResponse(w, "ok")
	}
}

func HandleGetClients(server *tcpserver.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients := server.ConnectedClients()

		jsonResponse(w, clients)
	}
}
