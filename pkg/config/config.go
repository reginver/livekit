package config

import (
	"os"
	"time"

	"github.com/livekit/protocol/logger"
	"gopkg.in/yaml.v3"
)

// Config holds the full server configuration.
type Config struct {
	Port     uint32   `yaml:"port"`
	BindAddr string   `yaml:"bind_addresses"`
	RTC      RTCConfig `yaml:"rtc"`
	Redis    RedisConfig `yaml:"redis"`
	Audio    AudioConfig `yaml:"audio"`
	Room     RoomConfig  `yaml:"room"`
	Logging  LoggingConfig `yaml:"logging"`
	Keys     map[string]string `yaml:"keys"`
	NodeIP   string `yaml:"node_ip"`
	Region   string `yaml:"region"`
}

// RTCConfig holds WebRTC-related configuration.
type RTCConfig struct {
	UDPPort         uint32   `yaml:"udp_port"`
	TCPPort         uint32   `yaml:"tcp_port"`
	ICEPortRangeStart uint32 `yaml:"port_range_start"`
	ICEPortRangeEnd   uint32 `yaml:"port_range_end"`
	STUNServers       []string `yaml:"stun_servers"`
	TURNServers       []TURNServer `yaml:"turn_servers"`
	UseExternalIP     bool   `yaml:"use_external_ip"`
}

// TURNServer holds TURN server configuration.
type TURNServer struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Protocol   string `yaml:"protocol"`
	Username   string `yaml:"username"`
	Credential string `yaml:"credential"`
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Address        string        `yaml:"address"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	DB             int           `yaml:"db"`
	DialTimeout    time.Duration `yaml:"dial_timeout"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	SentinelAddrs  []string      `yaml:"sentinel_addresses"`
	SentinelMaster string        `yaml:"sentinel_master_name"`
}

// AudioConfig holds audio processing configuration.
type AudioConfig struct {
	ActiveLevel     uint8 `yaml:"active_level"`
	MinPercentile   uint8 `yaml:"min_percentile"`
	UpdateInterval  uint32 `yaml:"update_interval"`
	SmoothIntervals uint32 `yaml:"smooth_intervals"`
}

// RoomConfig holds default room settings.
type RoomConfig struct {
	AutoCreate         bool          `yaml:"auto_create"`
	EmptyTimeout       uint32        `yaml:"empty_timeout"`
	MaxParticipants    uint32        `yaml:"max_participants"`
	EnableRemoteUnmute bool          `yaml:"enable_remote_unmute"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	JSON      bool   `yaml:"json"`
	Level     string `yaml:"level"`
	Sample    bool   `yaml:"sample"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:     7880,
		BindAddr: "",
		RTC: RTCConfig{
			UDPPort:       7882,
			TCPPort:       7881,
			UseExternalIP: false,
		},
		Room: RoomConfig{
			AutoCreate:   true,
			EmptyTimeout: 300,
		},
		Audio: AudioConfig{
			ActiveLevel:    35,
			MinPercentile:  10,
			UpdateInterval: 500,
			SmoothIntervals: 2,
		},
		Logging: LoggingConfig{
			Level: "info",
		},
	}
}

// NewConfig loads configuration from a YAML file path.
// If path is empty, it returns the default configuration.
func NewConfig(configFile string) (*Config, error) {
	conf := DefaultConfig()
	if configFile == "" {
		return conf, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, conf); err != nil {
		return nil, err
	}

	logger.Infow("loaded configuration", "file", configFile)
	return conf, nil
}
