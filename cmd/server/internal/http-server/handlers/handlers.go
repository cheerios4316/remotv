package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"remotv/cmd/server/command"
	tcpserver "remotv/cmd/server/internal/tcp-server"
	"strings"
	"time"

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

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		result, err := tcpServer.SendCommand(ctx, client, msg)
		if err != nil {
			status := http.StatusBadGateway
			if ctx.Err() == context.DeadlineExceeded {
				status = http.StatusGatewayTimeout
			}
			http.Error(w, "Error receiving client response: "+err.Error(), status)
			return
		}

		// expected: ok=<0|1>|<msg>
		// eg. "ok=1|" or "ok=0|error doing stuff"
		parts := strings.SplitN(result, "|", 2)

		if len(parts) != 2 {
			jsonResponse(w, "received wrong response from called client")
			return
		}

		okStr, errMsg := parts[0], parts[1]
		okParts := strings.SplitN(okStr, "=", 2)
		if len(okParts) != 2 {
			jsonResponse(w, "received wrong response from called client")
			return
		}

		if okParts[1] != "1" {
			jsonResponse(w, fmt.Sprintf("called client failed with error %q", errMsg))
			return
		}

		jsonResponse(w, "request executed on the called client successfully")
	}
}

func HandleGetClients(server *tcpserver.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients := server.ConnectedClients()

		jsonResponse(w, clients)
	}
}
