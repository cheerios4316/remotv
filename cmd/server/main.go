package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	httpserver "remotv/cmd/server/internal/http-server"
	"remotv/cmd/server/internal/input"
	tcpserver "remotv/cmd/server/internal/tcp-server"
	"syscall"
)

func main() {
	fmt.Println("=== Remotv Server ===")
	fmt.Println("Starting up...")
	flags := input.ParseFlags()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down...")
	}()

	tcpServer := tcpserver.Server{Port: flags.TcpPort}
	go tcpServer.Listen()
	fmt.Printf("Listening for TCP connections on port %d\n", flags.TcpPort)

	httpServer := httpserver.Server{Port: flags.HttpPort, TcpServer: &tcpServer}
	fmt.Printf("Listening for HTTP requests on port %d\n", flags.HttpPort)
	go httpServer.Listen()

	<-ctx.Done()
}
