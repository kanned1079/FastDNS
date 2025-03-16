package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 向日志缓冲区中添加一条日志
func (l *Logger) AddLog(logMessage string) {
	// 获取当前时间
	l.currentTime = time.Now().Format("2006/01/02 15:04:05")
	// 将日志添加到缓冲区
	l.buffer = append(l.buffer, fmt.Sprintf("%s %s", l.currentTime, logMessage))

	// 如果缓冲区的大小已达到指定数量，写入文件
	//log.Println("size: ", len(l.buffer))
	if len(l.buffer) >= l.bufferSize {
		l.writeLogsToFile()
	}
	l.buffer = []string{}
}

// 将日志写入文件
func (l *Logger) writeLogsToFile() {
	// 打开文件（如果文件不存在则创建）
	log.Println("starting insert logs...")
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer func() {
		log.Println("fresh log file success.")
		if err := file.Close(); err != nil {
			log.Println("failure close file")
		}
	}()

	// 将缓冲区中的日志写入文件
	for _, logEntry := range l.buffer {
		_, err := file.WriteString(logEntry)
		if err != nil {
			log.Printf("Error writing log to file: %v", err)
			return
		}
	}

	// 清空缓冲区

}
