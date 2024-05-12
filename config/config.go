package config

type Config struct {
	Server   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	Database DatabaseConfig `yaml:"database" env-prefix:"DB_"`
	Elastic  ElasticConfig  `yaml:"elastic" env-prefix:"ES_"`
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
	Host   string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port   uint   `yaml:"port" env:"PORT" env-default:"42069"`
	Prefix string `yaml:"prefix" env:"PREFIX" env-default:"services"`
}

type ElasticConfig struct {
	Port     int    `yaml:"port" env:"PORT" env-default:"9200"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	CertPath string `yaml:"cert_path" env:"CERT_PATH"`
	ApiKey   string `yaml:"api_key" env:"API_KEY"`
}
