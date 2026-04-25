package parser

import (
	"fmt"
	"strings"
)

// Selector 目标选择器
type Selector struct {
	Type       string            // @e, @a, @s, @f, @p, ...
	EntityType string            // 实体类型：file, process
	Filters    map[string]string // 附加过滤条件
}

// ParseSelector 解析目标选择器
// 支持: @f[name="test.txt"], @file[name="test.txt"]
//
//	@p[pid=1234], @proc[pid=1234], @process[name="sleep"]
func ParseSelector(input string) (*Selector, error) {
	input = strings.TrimSpace(input)

	if !strings.HasPrefix(input, "@") {
		return nil, fmt.Errorf("无效的选择器，必须以 @ 开头")
	}

	sel := &Selector{
		Filters: make(map[string]string),
	}

	// 查找 [
	bracketStart := strings.Index(input, "[")
	selectorType := input
	filterStr := ""

	if bracketStart != -1 {
		// 确保以 ] 结尾
		if !strings.HasSuffix(input, "]") {
			return nil, fmt.Errorf("选择器未闭合: 缺少 ]")
		}
		selectorType = input[:bracketStart]
		filterStr = input[bracketStart+1 : len(input)-1]
	}

	sel.Type = selectorType

	// 根据选择器类型推断实体类型
	switch selectorType {
	case "@f", "@file":
		sel.EntityType = "file"
	case "@p", "@proc", "@process":
		sel.EntityType = "process"
	default:
		return nil, fmt.Errorf("不支持的选择器类型: %s (可用: @f/@file 用于文件, @p/@proc/@process 用于进程)", selectorType)
	}

	// 解析过滤器
	if filterStr != "" {
		if err := parseFilterString(filterStr, sel.Filters); err != nil {
			return nil, err
		}
	}

	return sel, nil
}

func parseFilterString(filterStr string, filters map[string]string) error {
	i := 0
	for i < len(filterStr) {
		// 跳过空格和逗号分隔符
		for i < len(filterStr) && (filterStr[i] == ' ' || filterStr[i] == ',') {
			i++
		}
		if i >= len(filterStr) {
			break
		}

		// 查找 key
		eq := strings.Index(filterStr[i:], "=")
		if eq == -1 {
			return fmt.Errorf("过滤器格式错误: 需要 key=value")
		}
		key := strings.TrimSpace(filterStr[i : i+eq])
		i += eq + 1 // 跳过 '='

		// 解析 value（可能带引号）
		var value string
		if i < len(filterStr) && filterStr[i] == '"' {
			i++ // 跳过前引号
			end := strings.Index(filterStr[i:], "\"")
			if end == -1 {
				return fmt.Errorf("字符串未闭合")
			}
			value = filterStr[i : i+end]
			i += end + 1 // 跳过后引号
		} else {
			end := i
			for end < len(filterStr) && filterStr[end] != ',' && filterStr[end] != ' ' {
				end++
			}
			value = filterStr[i:end]
			i = end
		}
		filters[key] = value
	}
	return nil
}

// String 返回选择器的字符串表示
func (s *Selector) String() string {
	if len(s.Filters) == 0 {
		return s.Type
	}
	parts := make([]string, 0, len(s.Filters))
	for k, v := range s.Filters {
		if strings.Contains(v, " ") {
			parts = append(parts, fmt.Sprintf("%s=\"%s\"", k, v))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return fmt.Sprintf("%s[%s]", s.Type, strings.Join(parts, ","))
}
