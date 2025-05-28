package config

// Config 结构体定义了应用程序的配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server" desc:"server configuration"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig 结构体定义了服务器配置
type ServerConfig struct {
	Addr string `mapstructure:"addr" desc:"Server address with port"`
}

// LoggingConfig 结构体定义了日志配置
type LoggingConfig struct {
	Level           string `mapstructure:"level"`
	OutputToConsole bool   `mapstructure:"output_to_console"`
	OutputToFile    bool   `mapstructure:"output_to_file"`
	OutputFilePath  string `mapstructure:"output_file_path"`
	UseColorLevel   bool   `mapstructure:"use_color_level"`
	MaxSizeMB       int    `mapstructure:"max_size_mb"`
	MaxBackups      int    `mapstructure:"max_backups"`
	MaxAgeDays      int    `mapstructure:"max_age_days"`
	Compress        bool   `mapstructure:"compress"`
	CallerSkip      int    `mapstructure:"caller_skip"`
	AddStacktrace   bool   `mapstructure:"add_stacktrace"`
}

// DatabaseConfig 结构体定义了数据库配置
type DatabaseConfig struct {
	Type       string `mapstructure:"type"`
	Connection string `mapstructure:"connection"`
	Migrate    bool   `mapstructure:"migrate"`
}
