package tcpserver

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Port int

	mu       sync.RWMutex
	clients  map[string]net.Conn
	sessions map[net.Conn]*session
}

type session struct {
	commands  chan struct{}
	responses chan string
	done      chan struct{}
}

// SendCommand serializes commands because replies do not carry request IDs.
// The connection handler is the only reader of the socket.
func (s *Server) SendCommand(ctx context.Context, conn net.Conn, message string) (string, error) {
	s.mu.RLock()
	client := s.sessions[conn]
	s.mu.RUnlock()
	if client == nil {
		return "", fmt.Errorf("client disconnected")
	}
	select {
	case client.commands <- struct{}{}:
	case <-ctx.Done():
		return "", ctx.Err()
	case <-client.done:
		return "", fmt.Errorf("client disconnected")
	}
	defer func() { <-client.commands }()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	// Closing an interrupted exchange prevents its late reply from being used
	// by a subsequent command, and interrupts a blocked write.
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	if _, err := fmt.Fprintln(conn, message); err != nil {
		conn.Close()
		return "", err
	}
	select {
	case response := <-client.responses:
		return response, nil
	case <-ctx.Done():
		conn.Close()
		return "", ctx.Err()
	case <-client.done:
		return "", fmt.Errorf("client disconnected")
	}
}

// ConnectedClients returns a snapshot keyed by each client's name.
func (s *Server) ConnectedClients() map[string]net.Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make(map[string]net.Conn, len(s.clients))
	for address, conn := range s.clients {
		clients[address] = conn
	}
	return clients
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		cn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(cn)
	}
}

func (s *Server) handleConnection(cn net.Conn) {
	defer cn.Close()

	address := cn.RemoteAddr().String()

	scanner := bufio.NewScanner(cn)
	if !scanner.Scan() {
		return
	}

	greet := scanner.Text()
	parts := strings.SplitN(greet, ":", 2)
	if len(parts) != 2 {
		return
	}

	key, val := parts[0], parts[1]
	var name string
	client := &session{commands: make(chan struct{}, 1), responses: make(chan string, 1), done: make(chan struct{})}

	if key == "name" && val != "" {
		name = val
		s.mu.Lock()
		if s.clients == nil {
			s.clients = make(map[string]net.Conn)
		}
		s.clients[name] = cn
		if s.sessions == nil {
			s.sessions = make(map[net.Conn]*session)
		}
		s.sessions[cn] = client
		fmt.Printf("[%v] Client connected: %v [%v]", time.Now().Format(time.RFC3339), name, address)
		s.mu.Unlock()
	} else {
		cn.Close()
		return
	}

	defer func() {
		close(client.done)
		s.mu.Lock()
		delete(s.sessions, cn)
		if s.clients[name] == cn {
			delete(s.clients, name)
		}
		s.mu.Unlock()
		fmt.Printf("[%v] client %v[%v] disconnected\n", time.Now().Format(time.RFC3339), name, address)
	}()

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("[%v]message from %v[%v]: %v\n", time.Now().Format(time.RFC3339), name, address, message)
		select {
		case client.responses <- message:
		default:
			// No command is waiting for this unsolicited message.
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

}
