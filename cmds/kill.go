package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/F7YM/MineShell/parser"
)

type KillCommand struct{}

func init() {
	Register(&KillCommand{})
}

func (k *KillCommand) Name() string {
	return "kill"
}

func (k *KillCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: kill <目标选择器>")
	}

	selector, err := parser.ParseSelector(args[0])
	if err != nil {
		return fmt.Errorf("选择器解析失败: %v", err)
	}

	switch selector.EntityType {
	case "file":
		return killFile(selector)
	case "process":
		return killProcess(selector)
	default:
		return fmt.Errorf("未知实体类型: %s (可用: file, process)", selector.EntityType)
	}
}

func (k *KillCommand) Help() string {
	return "kill @e[type=file,name=\"文件名\"]  -  删除文件\n" +
		"     kill @e[type=process,pid=1145]  -  终止进程\n" +
		"     kill @e[type=process,name=\"进程名\"]  -  按名称终止进程"
}

func killFile(sel *parser.Selector) error {
	name, ok := sel.Filters["name"]
	if !ok {
		return fmt.Errorf("缺少 name 过滤器")
	}

	if err := os.Remove(name); err != nil {
		return fmt.Errorf("无法删除文件 %s: %v", name, err)
	}

	fmt.Printf("杀死了%s\n", name)
	return nil
}

func killProcess(sel *parser.Selector) error {
	// 优先使用 pid
	if pidStr, ok := sel.Filters["pid"]; ok {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return fmt.Errorf("无效的 pid: %s", pidStr)
		}
		return killByPid(pid)
	}

	// 按名称查找进程
	if name, ok := sel.Filters["name"]; ok {
		return killByName(name)
	}

	return fmt.Errorf("需要 pid 或 name 过滤器")
}

func killByPid(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("找不到进程 %d: %v", pid, err)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		err = process.Kill()
		if err != nil {
			return fmt.Errorf("无法终止进程 %d: %v", pid, err)
		}
	}

	fmt.Printf("杀死了%d\n", pid)
	return nil
}

func killByName(name string) error {
	// 使用 pgrep 查找进程
	cmd := exec.Command("pgrep", name)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("找不到名为 %s 的进程", name)
	}

	pids := strings.TrimSpace(string(output))
	if pids == "" {
		return fmt.Errorf("找不到名为 %s 的进程", name)
	}

	// 杀死所有匹配的进程
	for _, pidStr := range strings.Split(pids, "\n") {
		pid, _ := strconv.Atoi(pidStr)
		if pid > 0 {
			if err := killByPid(pid); err != nil {
				fmt.Printf("%d无法被杀死: %v\n", pid, err)
			}
		}
	}

	return nil
}
