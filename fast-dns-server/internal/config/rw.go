package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// FreshConfigFile2Disk 将当前的配置写入到文件

// FreshConfigFile2Disk 将当前的配置写入到文件
func (this *RootConfig) FreshConfigFile2Disk(filePath string) error {
	// 序列化结构体为 YAML
	yamlData, err := yaml.Marshal(this)
	if err != nil {
		return fmt.Errorf("unable to marshal yaml content: %v", err)
	}

	// 使用 os.WriteFile 替代 ioutil.WriteFile
	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("unable to write file: %v", err)
	}

	return nil
}

// ReadConfigFile 从文件中读取配置并填充结构体
func (this *RootConfig) ReadConfigFile(filePath string) error {
	// 读取文件内容，使用 os.ReadFile 替代 ioutil.ReadFile
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	// 反序列化 YAML 到结构体
	err = yaml.Unmarshal(fileContent, this)
	if err != nil {
		return fmt.Errorf("unable to unmarshal yaml content: %v", err)
	}

	return nil
}
