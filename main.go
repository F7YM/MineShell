package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/F7YM/MineShell/cmds"
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
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("MineShell - Minecraft Style Shell")
	fmt.Println("输入 'help' 查看命令，'exit' 退出")

	// 缓存用户信息（避免每次循环都获取）
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

	for {
		// 获取并简化工作目录
		wd, _ := os.Getwd()
		shortWd := simplifyPath(wd, homeDir)

		// 彩色提示符：用户@主机:路径
		fmt.Printf("%s%s%s@%s%s%s:%s%s%s> ",
			colorGreen, username, colorReset,
			colorCyan, hostname, colorReset,
			colorBlue, shortWd, colorReset,
		)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		args := parts[1:]

		if cmd, exists := cmds.GetCommand(cmdName); exists {
			if err := cmd.Execute(args); err != nil {
				fmt.Printf("%s错误: %v%s\n", colorRed, err, colorReset)
			}
		} else {
			fmt.Printf("%s未知命令: %s，输入 'help' 查看可用命令%s\n", colorRed, cmdName, colorReset)
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
