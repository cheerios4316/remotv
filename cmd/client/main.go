package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"remotv/cmd/client/internal/command"
	"remotv/cmd/client/internal/input"
	"remotv/cmd/client/internal/routines"
	"syscall"
)

func main() {
	flags := input.ParseFlags()
	cfg, err := input.LoadConfig(flags.ConfigPath)
	if err != nil {
		fmt.Printf("error in reading custom config file: using default config. %v\n", err.Error())
		cfg = input.DefaultConfig()
	}

	fmt.Printf("running with config %v\n", cfg)

	ipAddr := fmt.Sprintf("%s:%d", flags.URI, flags.Port)

	conn, err := net.Dial("tcp", ipAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		fmt.Fprintln(conn, "Disconnecting")
		fmt.Println("shutting down...")
		conn.Close()
	}()

	fmt.Fprintf(conn, "name:%v\n", flags.DeviceName)

	errCh := make(chan error)

	go routines.HandleReceiveMessage(ctx, errCh, conn, command.Runner{Config: cfg})

	go func() {
		err := <-errCh

		fmt.Fprintln(os.Stderr, "connection closed:", err)
		stop()
	}()

	<-ctx.Done()
	<-shutdownDone
}
