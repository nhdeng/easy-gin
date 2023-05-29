package easy

import (
	"log"
	"os"
)

// LoadConfigFile 读取配置文件
func LoadConfigFile() []byte {
	dir, _ := os.Getwd()
	file := dir + "/application.yaml"
	b, err := os.ReadFile(file)
	if err != nil {
		log.Panicln("读取配置文件失败", err)
	}
	return b
}
