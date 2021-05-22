package config

// Config struct for application config.
type Config struct {
	Server    Server    `yaml:"server" validate:"required"`
	Logger    Logger    `yaml:"logger" validate:"required"`
	Database  Database  `yaml:"database" validate:"required"`
	Telemetry Telemetry `yaml:"telemetry" validate:"required"`
}

// Server holds data necessary for server configuration.
type Server struct {
	Name    string  `yaml:"name" validate:"required"`
	Host    string  `yaml:"host" validate:"required,gt=0"`
	Port    int     `yaml:"port" validate:"required,gt=0"`
	Timeout Timeout `yaml:"timeout" validate:"required"`
}

// Timeout holds server timeout settings.
type Timeout struct {
	Shutdown int `yaml:"shutdown" validate:"required"`
	Write    int `yaml:"write" validate:"required"`
	Read     int `yaml:"read" validate:"required"`
	Idle     int `yaml:"idle" validate:"required"`
}

// Logger holds logger specific configuration.
type Logger struct {
	Name       string  `yaml:"name" validate:"required"`
	Level      string  `yaml:"level" validate:"required"`
	JSONFormat bool    `yaml:"jsonFormat"`
	Writers    Writers `yaml:"writers" validate:"required"`
}

// Writers holds available logger writers.
type Writers struct {
	FileWriter FileWriter `yaml:"file" validate:"required"`
}

// FileWriter holds file writer configuration.
type FileWriter struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path" validate:"required"`
}

// Database holds databases configuration.
type Database struct {
	Mongo MongoDB `yaml:"mongo" validate:"required"`
}

// MongoDB holds MongoDB specific configuration.
type MongoDB struct {
	URI      string `yaml:"uri" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Pass     string `yaml:"pass" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

// Telemetry holds telemetry specific configuration.
type Telemetry struct {
	ServiceName string `yaml:"name" validate:"required"`
	Level       string `yaml:"level" validate:"required"`
	AccessToken string `yaml:"token" validate:"required"`
}
