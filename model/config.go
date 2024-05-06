package model

type Config struct {
	Database DatabaseConfig `yaml:"database" env-prefix:"DB_"`
	Server   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	LogLevel string         `yaml:"log-level" env:"LOG_LEVEL" env-default:"INFO"`
}

type DatabaseConfig struct {
	Port     int    `yaml:"port" env:"PORT"`
	Host     string `yaml:"host" env:"HOST"`
	Name     string `yaml:"name" env:"NAME"`
	User     string `yaml:"user" env:"USER"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port uint   `yaml:"port" env:"PORT" env-default:"42069"`
}
