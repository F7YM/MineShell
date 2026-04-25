package cmds

import (
	"fmt"
	"os"
)

type SummonCommand struct{}

func init() {
	Register(&SummonCommand{})
}

func (s *SummonCommand) Name() string {
	return "summon"
}

func (s *SummonCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: summon <文件名>")
	}

	fileName := args[0]
	// 创建空文件
	os.WriteFile(fileName, []byte{}, 0644)

	return nil
}

func (s *SummonCommand) Help() string {
	return "创建空文件"
}
