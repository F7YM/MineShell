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
		return fmt.Errorf("用法: execute as <用户> run <命令> 或 execute at <目录> run <命令> 或 execute as <用户> at <目录> run <命令>")
	}

	var user string
	var dir string
	var command string

	// 解析 as, at, run
	i := 0
	n := len(args)
	for i < n {
		switch args[i] {
		case "as":
			if i+1 >= n {
				return fmt.Errorf("缺少用户名")
			}
			user = args[i+1]
			i += 2
		case "at":
			if i+1 >= n {
				return fmt.Errorf("缺少目录")
			}
			dir = args[i+1]
			i += 2
		case "run":
			if i+1 >= n {
				return fmt.Errorf("缺少命令")
			}
			command = strings.Join(args[i+1:], " ")
			i = n
		default:
			return fmt.Errorf("未知子命令: %s，可用: as, at, run", args[i])
		}
	}

	if command == "" {
		return fmt.Errorf("缺少 run 子命令")
	}

	return executeWithOptions(user, dir, command)
}

func (e *ExecuteCommand) Help() string {
	return "以指定用户身份执行命令，或直接执行系统命令"
}

func executeWithOptions(user, dir, command string) error {
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}

	// 确定是否切换用户
	switchUser := user != "" && user != currentUser && user != "@s"

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 如果需要切换用户
	if switchUser {
		// 构建命令：先 cd 到指定目录（如果有），再执行 mineshell -c
		var cmdStr string
		if dir != "" {
			cmdStr = fmt.Sprintf("cd %s && %s -c %s", shellQuote(dir), execPath, shellQuote(command))
		} else {
			cmdStr = fmt.Sprintf("%s -c %s", execPath, shellQuote(command))
		}
		cmd := exec.Command("sudo", "-u", user, "sh", "-c", cmdStr)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// 不切换用户，当前进程执行
	// 处理目录切换
	var oldDir string
	if dir != "" {
		oldDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("获取当前目录失败: %v", err)
		}
		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("切换目录到 %s 失败: %v", dir, err)
		}
		defer os.Chdir(oldDir)
	}

	// 执行命令（内建或外部）
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("没有指定命令")
	}
	cmdName := parts[0]
	cmdArgs := parts[1:]

	if cmd, exists := GetCommand(cmdName); exists {
		return cmd.Execute(cmdArgs)
	}

	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("找不到命令: %s", cmdName)
	}
	sysCmd := exec.Command(cmdPath, cmdArgs...)
	sysCmd.Stdin = os.Stdin
	sysCmd.Stdout = os.Stdout
	sysCmd.Stderr = os.Stderr
	return sysCmd.Run()
}

// shellQuote 简单对字符串加引号，避免空格问题
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
