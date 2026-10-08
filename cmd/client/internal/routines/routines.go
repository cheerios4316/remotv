package routines

import (
	"bufio"
	"context"
	"fmt"
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

		result := "ok=1|"
		if err := runner.Run(message); err != nil {
			result = fmt.Sprintf("ok=0|%s", err.Error())
		}
		if _, err := fmt.Fprintln(conn, result); err != nil {
			select {
			case ch <- err:
			case <-ctx.Done():
			}
			return
		}
	}
}
