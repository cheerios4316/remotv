package routines

import (
	"bufio"
	"context"
	"net"
	"remotv/cmd/client/internal/command"
)

func HandleReceiveMessage(ctx context.Context, ch chan error, conn net.Conn, runner command.Runner) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			ch <- err
			return
		}

		runner.Run(message)
	}
}
