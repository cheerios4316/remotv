package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Port int

	mu      sync.RWMutex
	clients map[string]net.Conn
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

	if key == "name" && val != "" {
		name = val
		s.mu.Lock()
		if s.clients == nil {
			s.clients = make(map[string]net.Conn)
		}
		s.clients[name] = cn
		fmt.Printf("[%v] Client connected: %v [%v]", time.Now().Format(time.RFC3339), name, address)
		s.mu.Unlock()
	} else {
		cn.Close()
		return
	}

	defer func() {
		s.mu.Lock()
		if s.clients[name] == cn {
			delete(s.clients, name)
		}
		s.mu.Unlock()
		fmt.Printf("[%v] client %v[%v] disconnected\n", time.Now().Format(time.RFC3339), name, address)
	}()

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("[%v]message from %v[%v]: %v\n", time.Now().Format(time.RFC3339), name, address, message)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

}
