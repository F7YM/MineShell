package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	wd, _ := os.Getwd()
	fmt.Printf("%s> ", wd)

	scanner.Scan()
	input := scanner.Text()
}
