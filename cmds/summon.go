// cmds/summon.go
package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/F7YM/MineShell/nbt"
)

type SummonCommand struct{}

func init() {
	Register(&SummonCommand{})
}

func (s *SummonCommand) Name() string {
	return "summon"
}

func (s *SummonCommand) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: summon <实体类型> <NBT>\n用法:\n  summon file {name:\"filename\", content:\"text\", mode:644}\n  summon process {command:\"ls\", args:[\"-la\", \"-h\"]}")
	}

	entityType := args[0]
	nbtStr := strings.Join(args[1:], " ")

	// 解析 SNBT
	node, err := nbt.ParseSNBT(nbtStr)
	if err != nil {
		return fmt.Errorf("SNBT 解析失败: %v", err)
	}

	switch entityType {
	case "file":
		return summonFile(node)
	case "process":
		return summonProcess(node)
	default:
		return fmt.Errorf("未知实体类型: %s (可用: file, process)", entityType)
	}
}

func (s *SummonCommand) Help() string {
	return "创建实体 (文件/进程)"
}

func summonFile(node *nbt.Node) error {
	name, ok := node.GetString("name")
	if !ok {
		return fmt.Errorf("缺少 'name' 字段 (字符串)")
	}

	content, _ := node.GetString("content")

	mode := os.FileMode(0644)
	if m, ok := node.GetInt("mode"); ok {
		// 八进制模式解析
		modeStr := fmt.Sprintf("%o", m)
		parsed, err := strconv.ParseUint(modeStr, 8, 32)
		if err == nil {
			mode = os.FileMode(parsed)
		}
	}

	return os.WriteFile(name, []byte(content), mode)
}

func summonProcess(node *nbt.Node) error {
	command, ok := node.GetString("command")
	if !ok {
		return fmt.Errorf("缺少 'command' 字段 (字符串)")
	}

	var cmdArgs []string
	if argsNode, ok := node.GetArray("args"); ok {
		for _, arg := range argsNode {
			if arg.Type == nbt.TypeString {
				cmdArgs = append(cmdArgs, arg.Value.(string))
			} else {
				cmdArgs = append(cmdArgs, fmt.Sprint(arg.Value))
			}
		}
	} else if argsStr, ok := node.GetString("args"); ok {
		cmdArgs = strings.Fields(argsStr)
	}

	cmd := exec.Command(command, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("召唤了新的", command)
	return cmd.Start()
}
