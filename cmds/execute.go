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

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 使用 sudo 重新调用自身执行命令
	cmd := exec.Command("sudo", "-u", player, execPath, "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// executeRun 直接执行命令
func executeRun(command string) error {
	// 先尝试解析为 MineShell 内建命令
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("没有指定命令")
	}

	cmdName := parts[0]
	args := parts[1:]

	// 检查是否是内建命令
	if cmd, exists := GetCommand(cmdName); exists {
		// 内建命令直接执行
		return cmd.Execute(args)
	}

	// 不是内建命令，作为系统命令执行
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("找不到命令: %s", cmdName)
	}

	sysCmd := exec.Command(cmdPath, args...)
	sysCmd.Stdin = os.Stdin
	sysCmd.Stdout = os.Stdout
	sysCmd.Stderr = os.Stderr

	return sysCmd.Run()
}
