package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/F7YM/MineShell/cmds"
	"github.com/chzyer/readline"
)

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {
	// 缓存用户信息
	currentUser, _ := user.Current()
	username := currentUser.Username
	if username == "" {
		username = "user"
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "host"
	}

	homeDir, _ := os.UserHomeDir()

	// 配置 readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "",                         // 提示符我们手动彩色输出，所以这里留空
		HistoryFile:            "~/.mineshell_history.tmp", // 历史记录文件
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		fmt.Printf("初始化 readline 失败: %v\n", err)
		return
	}
	defer rl.Close()

	fmt.Println("MineShell - Minecraft Style Shell")
	fmt.Println("输入 'help' 查看命令，'exit' 退出")

	for {
		// 获取并简化工作目录
		wd, _ := os.Getwd()
		shortWd := simplifyPath(wd, homeDir)

		// 手动构建彩色提示符
		prompt := fmt.Sprintf("%s%s%s@%s%s%s:%s%s%s> ",
			colorGreen, username, colorReset,
			colorCyan, hostname, colorReset,
			colorBlue, shortWd, colorReset,
		)

		// 设置 readline 的提示符（每次循环可以动态改变）
		rl.SetPrompt(prompt)

		// 读取一行
		input, err := rl.Readline()
		if err != nil {
			// Ctrl+D 或 exit 时退出
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 将输入加入历史
		rl.SaveHistory(input)

		if input == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		args := parts[1:]

		// 先查找内部命令
		if cmd, exists := cmds.GetCommand(cmdName); exists {
			if err := cmd.Execute(args); err != nil {
				fmt.Printf("%s错误: %v%s\n", colorRed, err, colorReset)
			}
		} else {
			// 尝试作为外部命令执行
			if err := executeExternalCommand(cmdName, args); err != nil {
				fmt.Printf("%s未知命令: %s，输入 'help' 查看可用命令%s\n", colorRed, cmdName, colorReset)
			}
		}
	}
}

// simplifyPath 将主目录替换为 ~
func simplifyPath(path, homeDir string) string {
	if strings.HasPrefix(path, homeDir) {
		return "~" + strings.TrimPrefix(path, homeDir)
	}
	return path
}

// executeExternalCommand 执行外部系统命令
func executeExternalCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
