package cmds

import (
	"fmt"
	"os/user"
	"strings"
)

type SayCommand struct{}

func init() {
	Register(&SayCommand{})
}

func (s *SayCommand) Name() string {
	return "say"
}

func (s *SayCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: say <消息>")
	}

	message := strings.Join(args, " ")
	userName, _ := user.Current()
	fmt.Printf("[%s] %s\n", userName.Username, message)
	return nil
}

func (s *SayCommand) Help() string {
	return "在屏幕上显示一条消息"
}
