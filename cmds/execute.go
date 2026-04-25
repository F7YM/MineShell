package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExecuteCommand struct{}

func init() {
	Register(&ExecuteCommand{})
}

func (e *ExecuteCommand) Name() string {
	return "execute"
}

func (e *ExecuteCommand) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: execute as <用户> run <命令> 或 execute run <命令>")
	}

	// 解析 execute 子命令
	switch args[0] {
	case "as":
		if len(args) < 4 || args[2] != "run" {
			return fmt.Errorf("用法: execute as <用户> run <命令>")
		}
		player := args[1]
		command := strings.Join(args[3:], " ")
		return executeAs(player, command)

	case "run":
		if len(args) < 2 {
			return fmt.Errorf("用法: execute run <命令>")
		}
		command := strings.Join(args[1:], " ")
		return executeRun(command)

	default:
		return fmt.Errorf("未知子命令: %s，可用: as, run", args[0])
	}
}

func (e *ExecuteCommand) Help() string {
	return "以指定用户身份执行命令，或直接执行系统命令"
}

// executeAs 以指定用户执行命令
func executeAs(player string, command string) error {
	// 获取当前用户
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}

	// 是当前用户，直接执行
	if player == currentUser || player == "@s" {
		return executeRun(command)
	}

	// 尝试使用 sudo 切换用户
	if player != currentUser {
		cmd := exec.Command("sudo", "-u", player, "sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return executeRun(command)
}

// executeRun 直接执行系统命令
func executeRun(command string) error {
	// 分割命令和参数
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("没有指定命令")
	}

	cmdName := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	// 查找命令路径
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("找不到命令: %s", cmdName)
	}

	// 执行命令
	cmd := exec.Command(cmdPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
