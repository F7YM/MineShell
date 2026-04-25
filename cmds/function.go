package cmds

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FunctionCommand struct{}

func init() {
	Register(&FunctionCommand{})
}

func (f *FunctionCommand) Name() string {
	return "function"
}

func (f *FunctionCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: function <文件名.mcfunction>")
	}

	filename := args[0]

	// 检查文件后缀名
	if filepath.Ext(filename) != ".mcfunction" {
		return fmt.Errorf("文件后缀名必须是 .mcfunction，当前文件名: %s", filename)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filename)
	}

	// 读取并执行 mcfunction 文件
	return executeFunction(filename)
}

func (f *FunctionCommand) Help() string {
	return "执行 mcfunction 文件，文件后缀名必须是 .mcfunction"
}

// executeFunction 读取 mcfunction 文件并逐行执行
func executeFunction(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行（# 开头）
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析命令
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		cmdArgs := parts[1:]

		// 先查找内部命令
		if cmd, exists := GetCommand(cmdName); exists {
			if err := cmd.Execute(cmdArgs); err != nil {
				return fmt.Errorf("第 %d 行错误: %v\n  命令: %s", lineNum, err, line)
			}
		} else {
			// 执行外部命令
			cmdPath, err := exec.LookPath(cmdName)
			if err != nil {
				return fmt.Errorf("第 %d 行错误: 找不到命令 '%s'\n  完整命令: %s", lineNum, cmdName, line)
			}

			cmd := exec.Command(cmdPath, cmdArgs...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("第 %d 行命令执行失败: %v\n  命令: %s", lineNum, err, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %v", err)
	}

	return nil
}
