package cmds

import (
	"fmt"
	"os"
)

type TPCommand struct{}

func init() {
	Register(&TPCommand{})
}

func (t *TPCommand) Name() string {
	return "tp"
}

func (t *TPCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: tp <目录路径>")
	}

	path := args[0]
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = home
	}

	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("无法进入: %v", err)
	}

	return nil
}

func (t *TPCommand) Help() string {
	return "tp <路径> - 传送（切换目录）"
}
