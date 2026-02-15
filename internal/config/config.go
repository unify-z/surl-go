package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

var ConfigManagerInstance *ConfigManager

type RateLimitConfig struct {
	Limit int `toml:"limit"`
}

type WebServerConfig struct {
	Host string
	Port int
}

type SMTPConfig struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type SiteConfig struct {
	BaseURL               string `toml:"base_url"`
	SiteName              string `toml:"site_name"`
	AllowGuestToCreateURL bool   `toml:"allow_guest_to_create_url"`
	AllowRegistration     bool   `toml:"allow_registration"`
}

type LogConfig struct {
	Level    string `toml:"level"`
	FilePath string `toml:"file_path"`
}

type DatabaseConfig struct {
	Path string `toml:"path"`
}

type JwtConfig struct {
	SecretKey     string `toml:"secret_key"`
	TokenDuration int    `toml:"token_duration"`
}

type AppConfig struct {
	WebServer  WebServerConfig `toml:"web_server"`
	Database   DatabaseConfig  `toml:"database"`
	RateLimits RateLimitConfig `toml:"rate_limits"`
	Jwt        JwtConfig       `toml:"jwt"`
	Site       SiteConfig      `toml:"site"`
	SMTP       SMTPConfig      `toml:"smtp"`
	Log        LogConfig       `toml:"log"`
}

type ConfigManager struct {
	Config AppConfig
}

func NewConfigManager(config AppConfig) *ConfigManager {
	return &ConfigManager{Config: config}
}

func (cm *ConfigManager) GetWebServerConfig() WebServerConfig {
	return cm.Config.WebServer
}

func (cm *ConfigManager) GetDatabaseConfig() DatabaseConfig {
	return cm.Config.Database
}

func (cm *ConfigManager) GetRateLimitConfig() RateLimitConfig {
	return cm.Config.RateLimits
}

func (cm *ConfigManager) GetAppConfig() AppConfig {
	return cm.Config
}

func Load() (*ConfigManager, error) {
	bytes, err := os.ReadFile("./config.toml")
	if err != nil {
		panic(err)
	}
	var config AppConfig
	if err := toml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
	ConfigManagerInstance = NewConfigManager(config)
	return ConfigManagerInstance, nil
}
