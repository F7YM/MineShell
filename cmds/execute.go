// cmds/execute.go
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
	if len(args) < 4 {
		return fmt.Errorf("用法: execute as <用户> run <命令>")
	}

	var user string
	var dir string
	var command string
	hasAs := false
	hasAt := false
	hasRun := false

	i := 0
	n := len(args)
	for i < n {
		switch args[i] {
		case "as":
			if hasAs {
				return fmt.Errorf("重复的子命令: as")
			}
			if i+1 >= n {
				return fmt.Errorf("缺少用户名")
			}
			user = args[i+1]
			hasAs = true
			i += 2
		case "at":
			if hasAt {
				return fmt.Errorf("重复的子命令: at")
			}
			if i+1 >= n {
				return fmt.Errorf("缺少目录")
			}
			dir = args[i+1]
			hasAt = true
			i += 2
		case "run":
			if hasRun {
				return fmt.Errorf("重复的子命令: run")
			}
			if i+1 >= n {
				return fmt.Errorf("缺少命令")
			}
			command = strings.Join(args[i+1:], " ")
			hasRun = true
			i = n
		default:
			return fmt.Errorf("未知子命令: %s，可用: as, at, run", args[i])
		}
	}

	if !hasRun {
		return fmt.Errorf("缺少 run 子命令")
	}

	if command == "" {
		return fmt.Errorf("run 后面缺少命令")
	}

	return executeWithOptions(user, dir, command)
}

func (e *ExecuteCommand) Help() string {
	return "execute as <用户> [at <目录>] run <命令>  -  以指定用户身份在指定目录执行命令"
}

func executeWithOptions(user, dir, command string) error {
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}

	switchUser := user != "" && user != currentUser && user != "@s"

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	if switchUser {
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

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
