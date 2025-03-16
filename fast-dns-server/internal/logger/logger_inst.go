package logger

import (
	"fast-dns-server/internal/config"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

var MyLogger *Logger

type Logger struct {
	buffer      []string
	bufferSize  int
	filePath    string
	currentTime string
	maxFileSize int64
	cronInst    *cron.Cron
}

func NewLogger(filePath string) *Logger {
	// 确保日志文件存在，如果文件不存在则创建一个空文件
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建一个空文件
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error creating log file: %v", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Println("failure close log file.")
			}
		}()

		fmt.Println("Log file created:", filePath)
	} else if err != nil {
		log.Fatalf("Error checking log file: %v", err)
	}

	return &Logger{
		buffer:      []string{},
		bufferSize:  config.RootCfg.Details.Config.LogBuffer,
		filePath:    filePath,
		maxFileSize: config.RootCfg.Details.Config.LogSizeLimit,
		cronInst:    cron.New(),
	}
}
