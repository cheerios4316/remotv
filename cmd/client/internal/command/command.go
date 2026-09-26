package command

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"remotv/cmd/client/internal/input"
)

type Command struct {
	Action string `json:"action"`
	Args   string `json:"args"`
}

type Runner struct {
	Config input.Config
}

func (r *Runner) Run(message string) {
	fmt.Printf("received message %q", message)

	var cmd Command
	err := json.Unmarshal([]byte(message), &cmd)
	if err != nil {
		fmt.Printf("invalid JSON message %q\n", message)
		return
	}

	bin, ok := r.getMap()[cmd.Action]
	if !ok {
		fmt.Printf("received invalid command %q\n", cmd.Action)
		return
	}

	osCmd := exec.Command(bin, cmd.Args)
	err = osCmd.Run()
	if err != nil {
		fmt.Printf("error during execution of received command: %s", err.Error())
	}
}

func (r *Runner) getMap() map[string]string {
	return map[string]string{
		"browser": r.Config.Browser,
	}
}
