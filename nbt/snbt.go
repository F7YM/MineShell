package nbt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// SNBT 节点类型
type NodeType int

const (
	TypeObject NodeType = iota
	TypeArray
	TypeString
	TypeInt
	TypeFloat
	TypeBool
)

// Node 表示一个 NBT 节点
type Node struct {
	Type     NodeType
	Value    any
	Children map[string]*Node // 用于 Object
	Items    []*Node          // 用于 Array
}

// ParseSNBT 解析 SNBT 字符串
func ParseSNBT(input string) (*Node, error) {
	input = strings.TrimSpace(input)
	p := &parser{
		input: []rune(input),
		pos:   0,
	}
	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("多余的字符: %s", string(p.input[p.pos:]))
	}
	return node, nil
}

// GetString 获取字符串值
func (n *Node) GetString(key string) (string, bool) {
	if n.Type != TypeObject {
		return "", false
	}
	child, ok := n.Children[key]
	if !ok || child.Type != TypeString {
		return "", false
	}
	return child.Value.(string), true
}

// GetInt 获取整数值
func (n *Node) GetInt(key string) (int64, bool) {
	if n.Type != TypeObject {
		return 0, false
	}
	child, ok := n.Children[key]
	if !ok {
		return 0, false
	}
	switch child.Type {
	case TypeInt:
		return child.Value.(int64), true
	case TypeFloat:
		return int64(child.Value.(float64)), true
	}
	return 0, false
}

// GetFloat 获取浮点数值
func (n *Node) GetFloat(key string) (float64, bool) {
	if n.Type != TypeObject {
		return 0, false
	}
	child, ok := n.Children[key]
	if !ok {
		return 0, false
	}
	switch child.Type {
	case TypeFloat:
		return child.Value.(float64), true
	case TypeInt:
		return float64(child.Value.(int64)), true
	}
	return 0, false
}

// GetBool 获取布尔值
func (n *Node) GetBool(key string) (bool, bool) {
	if n.Type != TypeObject {
		return false, false
	}
	child, ok := n.Children[key]
	if !ok || child.Type != TypeBool {
		return false, false
	}
	return child.Value.(bool), true
}

// GetObject 获取嵌套对象
func (n *Node) GetObject(key string) (*Node, bool) {
	if n.Type != TypeObject {
		return nil, false
	}
	child, ok := n.Children[key]
	if !ok || child.Type != TypeObject {
		return nil, false
	}
	return child, true
}

// GetArray 获取数组
func (n *Node) GetArray(key string) ([]*Node, bool) {
	if n.Type != TypeObject {
		return nil, false
	}
	child, ok := n.Children[key]
	if !ok || child.Type != TypeArray {
		return nil, false
	}
	return child.Items, true
}

// Has 检查 key 是否存在
func (n *Node) Has(key string) bool {
	if n.Type != TypeObject {
		return false
	}
	_, ok := n.Children[key]
	return ok
}

// String 递归打印节点树
func (n *Node) String() string {
	return n.stringIndent(0)
}

func (n *Node) stringIndent(indent int) string {
	prefix := strings.Repeat("  ", indent)
	switch n.Type {
	case TypeString:
		return fmt.Sprintf(`"%s"`, n.Value)
	case TypeInt:
		return fmt.Sprintf("%d", n.Value)
	case TypeFloat:
		return fmt.Sprintf("%f", n.Value)
	case TypeBool:
		return fmt.Sprintf("%t", n.Value)
	case TypeObject:
		if len(n.Children) == 0 {
			return "{}"
		}
		result := "{\n"
		keys := make([]string, 0, len(n.Children))
		for k := range n.Children {
			keys = append(keys, k)
		}
		for i, k := range keys {
			result += fmt.Sprintf("%s  %s: %s", prefix, k, n.Children[k].stringIndent(indent+1))
			if i < len(keys)-1 {
				result += ","
			}
			result += "\n"
		}
		result += prefix + "}"
		return result
	case TypeArray:
		if len(n.Items) == 0 {
			return "[]"
		}
		result := "[\n"
		for i, item := range n.Items {
			result += prefix + "  " + item.stringIndent(indent+1)
			if i < len(n.Items)-1 {
				result += ","
			}
			result += "\n"
		}
		result += prefix + "]"
		return result
	}
	return "?"
}

// parser 内部解析器
type parser struct {
	input []rune
	pos   int
}

func (p *parser) peek() rune {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *parser) next() rune {
	r := p.peek()
	p.pos++
	return r
}

func (p *parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(p.input[p.pos]) {
		p.pos++
	}
}

func (p *parser) parseValue() (*Node, error) {
	p.skipWhitespace()
	switch c := p.peek(); {
	case c == '{':
		return p.parseObject()
	case c == '[':
		return p.parseArray()
	case c == '"' || c == '\'':
		return p.parseString()
	case c == 't' || c == 'f':
		return p.parseBool()
	case c == '-' || (c >= '0' && c <= '9'):
		return p.parseNumber()
	case c == 0:
		return nil, fmt.Errorf("意外的 EOF")
	default:
		return nil, fmt.Errorf("不支持的字符 '%c' (位置 %d)", c, p.pos)
	}
}

func (p *parser) parseObject() (*Node, error) {
	node := &Node{
		Type:     TypeObject,
		Children: make(map[string]*Node),
	}
	p.next() // 跳过 '{'
	p.skipWhitespace()

	if p.peek() == '}' {
		p.next()
		return node, nil
	}

	for {
		p.skipWhitespace()
		key, err := p.parseKey()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if p.peek() != ':' {
			return nil, fmt.Errorf("期望 ':' 但得到 '%c' (位置 %d)", p.peek(), p.pos)
		}
		p.next() // 跳过 ':'
		p.skipWhitespace()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		node.Children[key] = value

		p.skipWhitespace()
		if p.peek() == '}' {
			p.next()
			return node, nil
		}
		if p.peek() != ',' {
			return nil, fmt.Errorf("期望 ',' 或 '}' 但得到 '%c' (位置 %d)", p.peek(), p.pos)
		}
		p.next() // 跳过 ','
	}
}

func (p *parser) parseArray() (*Node, error) {
	node := &Node{
		Type:  TypeArray,
		Items: make([]*Node, 0),
	}
	p.next() // 跳过 '['
	p.skipWhitespace()

	if p.peek() == ']' {
		p.next()
		return node, nil
	}

	for {
		p.skipWhitespace()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		node.Items = append(node.Items, value)

		p.skipWhitespace()
		if p.peek() == ']' {
			p.next()
			return node, nil
		}
		if p.peek() != ',' {
			return nil, fmt.Errorf("期望 ',' 或 ']' 但得到 '%c' (位置 %d)", p.peek(), p.pos)
		}
		p.next() // 跳过 ','
	}
}

func (p *parser) parseKey() (string, error) {
	if p.peek() == '"' || p.peek() == '\'' {
		node, err := p.parseString()
		if err != nil {
			return "", err
		}
		return node.Value.(string), nil
	}

	start := p.pos
	for p.pos < len(p.input) && (unicode.IsLetter(p.input[p.pos]) || unicode.IsDigit(p.input[p.pos]) || p.input[p.pos] == '_' || p.input[p.pos] == '-') {
		p.pos++
	}
	if start == p.pos {
		return "", fmt.Errorf("期望 key 但得到 '%c' (位置 %d)", p.peek(), p.pos)
	}
	return string(p.input[start:p.pos]), nil
}

func (p *parser) parseString() (*Node, error) {
	quote := p.next() // 获取引号
	var result strings.Builder

	for {
		c := p.next()
		if c == 0 {
			return nil, fmt.Errorf("未闭合的字符串")
		}
		if c == '\\' {
			escaped := p.next()
			switch escaped {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case '\\':
				result.WriteRune('\\')
			case '"':
				result.WriteRune('"')
			case '\'':
				result.WriteRune('\'')
			default:
				result.WriteRune(escaped)
			}
		} else if c == quote {
			return &Node{Type: TypeString, Value: result.String()}, nil
		} else {
			result.WriteRune(c)
		}
	}
}

func (p *parser) parseBool() (*Node, error) {
	if len(p.input)-p.pos >= 4 && string(p.input[p.pos:p.pos+4]) == "true" {
		p.pos += 4
		return &Node{Type: TypeBool, Value: true}, nil
	}
	if len(p.input)-p.pos >= 5 && string(p.input[p.pos:p.pos+5]) == "false" {
		p.pos += 5
		return &Node{Type: TypeBool, Value: false}, nil
	}
	return nil, fmt.Errorf("无效的布尔值 (位置 %d)", p.pos)
}

func (p *parser) parseNumber() (*Node, error) {
	start := p.pos
	isFloat := false
	for p.pos < len(p.input) && (p.input[p.pos] == '-' || p.input[p.pos] == '+' || p.input[p.pos] == '.' || (p.input[p.pos] >= '0' && p.input[p.pos] <= '9') || p.input[p.pos] == 'e' || p.input[p.pos] == 'E') {
		if p.input[p.pos] == '.' || p.input[p.pos] == 'e' || p.input[p.pos] == 'E' {
			isFloat = true
		}
		p.pos++
	}
	numStr := string(p.input[start:p.pos])

	if isFloat {
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("无效的浮点数 '%s': %v", numStr, err)
		}
		return &Node{Type: TypeFloat, Value: val}, nil
	}

	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的整数 '%s': %v", numStr, err)
	}
	return &Node{Type: TypeInt, Value: val}, nil
}
