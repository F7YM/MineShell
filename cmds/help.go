package cmds

import (
	"fmt"
	"sort"
)

type HelpCommand struct{}

func init() {
	Register(&HelpCommand{})
}

func (h *HelpCommand) Name() string {
	return "help"
}

func (h *HelpCommand) Execute(args []string) error {
	fmt.Println("可用的内建命令:")

	commands := GetAllCommands()
	// 按命令名排序
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	for _, cmd := range commands {
		fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Help())
	}
	return nil
}

func (h *HelpCommand) Help() string {
	return "显示帮助信息"
}
