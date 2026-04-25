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
	if len(args) > 0 {
		// help <命令名> 显示具体命令的帮助
		for _, cmdName := range args {
			if cmd, exists := GetCommand(cmdName); exists {
				fmt.Printf("%s\n  用法: %s\n", cmd.Name(), cmd.Help())
			} else {
				fmt.Printf("未知命令: %s\n", cmdName)
			}
		}
		return nil
	}

	fmt.Println("可用的内建命令:")

	commands := GetAllCommands()
	// 按命令名排序
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	for _, cmd := range commands {
		fmt.Printf("  %s\n", cmd.Help())
	}
	return nil
}

func (h *HelpCommand) Help() string {
	return "help [命令名]  -  显示帮助信息"
}
