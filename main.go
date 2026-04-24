package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/F7YM/MineShell/cmds"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("MineShell - Minecraft Style Shell")
	fmt.Print("输入 'help' 查看命令，'exit' 退出\n\n")

	for {
		// 获取当前工作目录作为提示符
		wd, err := os.Getwd()
		if err != nil {
			wd = "?"
		}
		fmt.Printf("%s> ", wd)

		if !scanner.Scan() {
			break // 遇到 EOF 或错误时退出
		}

		// 获取并处理输入
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		// 按空格分割命令和参数
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		args := parts[1:]

		// 从 cmds 包查找并执行命令
		if cmd, exists := cmds.GetCommand(cmdName); exists {
			if err := cmd.Execute(args); err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Printf("未知命令: %s，输入 'help' 查看可用命令\n", cmdName)
		}
	}
}
