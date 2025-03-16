package utils

import (
	"fmt"
	"runtime"
	"time"
)

func (Utils) clearConsole() {
	if runtime.GOOS == "windows" {
		// Windows平台
		fmt.Print("\033[H\033[2J")
	} else {
		// Linux 或 macOS
		fmt.Print("\033[H\033[2J")
	}
}

func (this *Utils) ShowLogEveryInterval() {
	for {
		// 显示日志
		fmt.Println("Logging some data...")

		// 清空控制台
		this.clearConsole()

		// 等待 5 秒钟
		time.Sleep(5 * time.Second)
	}
}
