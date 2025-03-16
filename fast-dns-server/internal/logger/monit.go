package logger

import (
	"fmt"
	"log"
	"os"
)

func (l *Logger) StartCheckSizeInterval() {
	_, err := l.cronInst.AddFunc("@daily", func() {
		// 每天检查日志文件大小并清空文件内容
		go l.checkAndClearLogFile()
	})
	if err != nil {
		log.Print("failure set interval")
	}
}

// 检查日志文件大小并清空文件内容
func (l *Logger) checkAndClearLogFile() {
	fileInfo, err := os.Stat(l.filePath)
	if err != nil {
		// 如果文件不存在，则创建文件
		if os.IsNotExist(err) {
			_, err := os.Create(l.filePath)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
		} else {
			fmt.Println("Error checking log file:", err)
			return
		}
	} else {
		// 如果文件大小超过指定大小，则清空文件
		if fileInfo.Size() > l.maxFileSize {
			err := os.Truncate(l.filePath, 0) // 清空文件内容
			if err != nil {
				fmt.Println("Error clearing log file:", err)
			} else {
				fmt.Println("Log file cleared due to exceeding max size.")
			}
		}
	}
}
