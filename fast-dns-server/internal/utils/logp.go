package utils

import (
	"fast-dns-server/internal/logger"
	"fmt"
	"time"
)

// ShowStatueLog 打印带有颜色的日志信息
func (this *Utils) ShowStatueLog(tpy, title, content string) {
	// 获取当前时间
	currentTime := time.Now().Format("2006/01/02 15:04:05")

	// 定义颜色的 ANSI 转义序列
	var colorCode string

	// 根据 tpy 设置 title 的颜色
	switch tpy {
	case "success":
		colorCode = "\033[32m" // 绿色
	case "warning":
		colorCode = "\033[33m" // 橙色
	case "error":
		colorCode = "\033[31m" // 红色
	case "primary":
		colorCode = "\033[35m" // 紫色
	case "info":
		colorCode = "\033[90m" // 灰色
	case "cyan":
		colorCode = "\033[36m" // 青色
	default:
		colorCode = "\033[37m" // 默认颜色 (白色)
	}

	// 打印日志，时间部分为灰色，title 部分为不同颜色，content 部分保持默认颜色
	// title 部分根据 tpy 设置颜色，content 部分保持默认颜色
	fmt.Printf("%s %s[%s]\033[0m %s\n", currentTime, colorCode, title, content)

	logger.MyLogger.AddLog(fmt.Sprintf("%s [%s] %s\n", currentTime, title, content))

}
