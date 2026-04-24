package cmds

// Command 定义所有命令必须实现的接口
type Command interface {
	// 执行命令
	Execute(args []string) error
	// 命令帮助信息
	Help() string
	// 命令名称（如 "give", "tp"）
	Name() string
}

// 全局命令注册表
var commandRegistry = make(map[string]Command)

// Register 注册命令（在 init() 中调用）
func Register(cmd Command) {
	commandRegistry[cmd.Name()] = cmd
}

// GetCommand 获取命令
func GetCommand(name string) (Command, bool) {
	cmd, exists := commandRegistry[name]
	return cmd, exists
}

// GetAllCommands 获取所有命令
func GetAllCommands() []Command {
	cmds := make([]Command, 0, len(commandRegistry))
	for _, cmd := range commandRegistry {
		cmds = append(cmds, cmd)
	}
	return cmds
}
