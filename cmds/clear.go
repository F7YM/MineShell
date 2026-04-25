package cmds

import (
	"os"
	"os/exec"
	"runtime"
)

type ClearCommand struct{}

func init() {
	Register(&ClearCommand{})
}

func (c *ClearCommand) Name() string {
	return "clear"
}

func (c *ClearCommand) Execute(args []string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (c *ClearCommand) Help() string {
	return "clear  -  清空屏幕"
}
