package tcpserver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestCommandResponsesAndTimeout(t *testing.T) {
	server := &Server{}
	client, connection := net.Pipe()
	defer client.Close()
	go server.handleConnection(connection)
	if _, err := fmt.Fprintln(client, "name:test-client"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for server.ConnectedClients()["test-client"] == nil {
		if time.Now().After(deadline) {
			t.Fatal("client was not registered")
		}
		time.Sleep(time.Millisecond)
	}
	go func() {
		reader := bufio.NewReader(client)
		for _, response := range []string{"ok=1|", "ok=0|execution failed"} {
			if _, err := reader.ReadString('\n'); err != nil {
				return
			}
			if _, err := fmt.Fprintln(client, response); err != nil {
				return
			}
		}
		// Receive the third command, but never acknowledge it.
		reader.ReadString('\n')
	}()
	for _, expected := range []string{"ok=1|", "ok=0|execution failed"} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		result, err := server.SendCommand(ctx, connection, "command")
		cancel()
		if err != nil || result != expected {
			t.Fatalf("response = %q, error = %v; want %q", result, err, expected)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := server.SendCommand(ctx, connection, "command")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v; want deadline exceeded", err)
	}
}

func TestTracksConnectedClientUntilDisconnect(t *testing.T) {
	server := &Server{}
	client, connection := net.Pipe()
	defer client.Close()

	done := make(chan struct{})
	go func() {
		server.handleConnection(connection)
		close(done)
	}()
	if _, err := fmt.Fprintln(client, "name:test-client"); err != nil {
		t.Fatal(err)
	}

	deadline := time.After(time.Second)
	for len(server.ConnectedClients()) != 1 {
		select {
		case <-deadline:
			t.Fatal("client was not registered")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	if server.ConnectedClients()["test-client"] != connection {
		t.Fatal("client was not registered by name")
	}

	client.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("connection handler did not stop after disconnect")
	}
	if len(server.ConnectedClients()) != 0 {
		t.Fatal("disconnected client was not removed")
	}
}

func TestDisconnectDoesNotRemoveReplacementClient(t *testing.T) {
	server := &Server{}
	connect := func() (net.Conn, <-chan struct{}) {
		t.Helper()
		client, connection := net.Pipe()
		t.Cleanup(func() { client.Close() })
		done := make(chan struct{})
		go func() {
			server.handleConnection(connection)
			close(done)
		}()
		if _, err := fmt.Fprintln(client, "name:test-client"); err != nil {
			t.Fatal(err)
		}
		deadline := time.After(time.Second)
		for server.ConnectedClients()["test-client"] != connection {
			select {
			case <-deadline:
				t.Fatal("client was not registered")
			default:
				time.Sleep(time.Millisecond)
			}
		}
		return client, done
	}
	oldClient, oldDone := connect()
	newClient, newDone := connect()
	oldClient.Close()
	select {
	case <-oldDone:
	case <-time.After(time.Second):
		t.Fatal("old connection handler did not stop")
	}
	if len(server.ConnectedClients()) != 1 {
		t.Fatal("old disconnect removed the replacement client")
	}
	newClient.Close()
	select {
	case <-newDone:
	case <-time.After(time.Second):
		t.Fatal("replacement connection handler did not stop")
	}
	if len(server.ConnectedClients()) != 0 {
		t.Fatal("disconnected replacement client was not removed")
	}
}
