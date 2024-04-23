package main

import (
	"fmt"
	"github.com/atotto/clipboard"
)

func main() {
	content := "Hello, clipboard!"

	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println("Failed to write to clipboard:", err)
		return
	}

	// 读取剪贴板内容
	content, err = clipboard.ReadAll()
	if err != nil {
		fmt.Println("读取剪贴板失败:", err)
		return
	}

	fmt.Printf("剪贴板内容: %s\n", content)
}
