package easy

import (
	"gopkg.in/yaml.v2"
	"log"
)

// SysConfig 系统配置
type SysConfig struct {
	Server   *serverConfig
	Log      *logConfig
	Jwt      *jwtConfig
	Database *DatabaseConfig
	Redis    *RedisConfig
	Menus    []*MenuConfig
	Access   []*AccessConfig
}

func (this *SysConfig) Name() string {
	return "SysConfig"
}

func NewSysConfig() *SysConfig {
	return &SysConfig{}
}

type serverConfig struct {
	Port int32
}

type DatabaseConfig struct {
	Uri string
}

type RedisConfig struct {
	Uri      string
	Password string
}

type logConfig struct {
	Level   string
	Path    string
	Name    string
	Expired int
}

type jwtConfig struct {
	SecretKey string `yaml:"secretKey"`
}

type MenuConfig struct {
	Name string `json:"name" gorm:"column:name"`
	Code string `json:"code" gorm:"column:code"`
}

type AccessConfig struct {
	Name     string `json:"name"`
	Uri      string `json:"Uri"`
	Method   string `json:"method"`
	MenuCode string `yaml:"menuCode" json:"menuCode"`
	State    int    `json:"state"` // 1:通用api
}

func InitConfig() *SysConfig {
	config := NewSysConfig()
	if b := LoadConfigFile(); b != nil {
		err := yaml.Unmarshal(b, config)
		if err != nil {
			log.Fatal("解析配置文件失败", err)
		}
	}
	return config
}
