package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var (
	Port                  = 3000
	ForcePort             = false
	LogDir                = ""
	Version               = "v0.0.0"
	IsMasterNode          = true
	DebugEnabled          = false
	MemoryCacheEnabled    = false
	SyncFrequency         = 60
	SessionSecret         = ""
	DisplayInCurrencyEnabled = false
	QuotaPerUnit          = 10000.0
	BatchUpdateEnabled    = false
	BatchUpdateInterval   = 5
	EnableMetric          = false
	Theme                 = "default"
	ServerAddress         = ""
	UIAddress             = ""
)

var GlobalConfig *Config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Minoru   MinoruConfig   `mapstructure:"minoru"`
	Upload   UploadConfig   `mapstructure:"upload"`
	OAuth    OAuthConfig    `mapstructure:"oauth"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MinoruConfig struct {
	Mode    string `mapstructure:"mode"`
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Timeout int    `mapstructure:"timeout"`
}

type UploadConfig struct {
	MaxSize     int64  `mapstructure:"max_size"`
	StoragePath string `mapstructure:"storage_path"`
}

type OAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	AuthURL      string `mapstructure:"auth_url"`
	TokenURL     string `mapstructure:"token_url"`
	UserInfoURL  string `mapstructure:"user_info_url"`
	RedirectURI  string `mapstructure:"redirect_uri"`
	APIKeyURL    string `mapstructure:"api_key_url"`
	LogoutURL    string `mapstructure:"logout_url"`
	WellKnownURL string `mapstructure:"well_known_url"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Name)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func InitConfig(cfgPath string) error {
	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func SetupConfig() {
	if port := os.Getenv("PORT"); port != "" {
		ForcePort = true
		Port, _ = strconv.Atoi(port)
	}
	if logDir := os.Getenv("LOG_DIR"); logDir != "" {
		LogDir = logDir
	}
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		ServerAddress = serverAddress
	}
	if uiAddress := os.Getenv("UI_ADDRESS"); uiAddress != "" {
		UIAddress = uiAddress
	}
	if debug := os.Getenv("DEBUG"); debug == "true" {
		DebugEnabled = true
	}
}
