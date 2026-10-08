package command

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"remotv/cmd/client/internal/input"
)

type CmdName string

const (
	CmdBrowser  CmdName = "browser"
	CmdShutdown CmdName = "shutdown"
)

type Command struct {
	Action CmdName `json:"action"`
	Args   string  `json:"args"`
}

func (c *Command) GetCommandBin(r *Runner) string {
	return map[CmdName]string{
		CmdBrowser:  r.Config.Browser,
		CmdShutdown: r.Config.Shutdown,
	}[c.Action]
}

func (c *Command) HasArgs() bool {
	return map[CmdName]bool{
		CmdBrowser:  true,
		CmdShutdown: false,
	}[c.Action]
}

func (c *Command) Validate() error {
	switch c.Action {
	case CmdBrowser, CmdShutdown:
		return nil
	default:
		return fmt.Errorf("unknown command %q provided", c.Action)
	}
}

type Runner struct {
	Config input.Config
}

func (r *Runner) ParseMessage(message string) (Command, error) {
	var cmd Command

	if err := json.Unmarshal([]byte(message), &cmd); err != nil {
		return Command{}, err
	}

	if !cmd.HasArgs() {
		cmd.Args = ""
	}

	return cmd, nil
}

func (r *Runner) Run(message string) error {
	cmd, err := r.ParseMessage(message)
	if err != nil {
		fmt.Printf("invalid JSON message %q\n", message)
		return err
	}

	if err := cmd.Validate(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Received request to execute command %q with args %q\n", cmd.Action, cmd.Args)

	bin := cmd.GetCommandBin(r)

	osCmd := exec.Command(bin, cmd.Args)
	err = osCmd.Run()
	if err != nil {
		fmt.Printf("error during execution of received command: %s\n", err.Error())
		return err
	}

	return nil
}
