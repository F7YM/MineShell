package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		wd, _ := os.Getwd() // 获取当前工作目录作为提示符
		fmt.Printf("%s> ", wd)

		scanner.Scan()
		input := scanner.Text()
		if !scanner.Scan() {
			break
		}
	}
}
