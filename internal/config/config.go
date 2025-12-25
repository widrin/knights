package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Center   CenterConfig   `yaml:"center"`
	Login    LoginConfig    `yaml:"login"`
	Gateway  GatewayConfig  `yaml:"gateway"`
	Game     GameConfig     `yaml:"game"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Name        string `yaml:"name"`
	Address     string `yaml:"address"`
	Port        int    `yaml:"port"`
	ServiceType string `yaml:"service_type"`
}

// CenterConfig holds center service configuration
type CenterConfig struct {
	Address           string `yaml:"address"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
	Timeout           int    `yaml:"timeout"`
}

// LoginConfig holds login service configuration
type LoginConfig struct {
	Address     string `yaml:"address"`
	TokenExpire int    `yaml:"token_expire"`
}

// GatewayConfig holds gateway service configuration
type GatewayConfig struct {
	Address        string `yaml:"address"`
	Port           int    `yaml:"port"`
	MaxConnections int    `yaml:"max_connections"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	ServerID       string `yaml:"server_id"`
	MaxPlayers     int    `yaml:"max_players"`
	TickRate       int    `yaml:"tick_rate"`
	RoomMaxPlayers int    `yaml:"room_max_players"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &config, nil
}
