package config

import "time"

type MysqlConfig struct {
	Username string `env:"USERNAME" envDefault:"root"`
	Password string `env:"PASSWORD" envDefault:"secret"`
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Address  string `env:"ADDRESS" envDefault:"127.0.0.1:3306"`
	Database string `env:"DATABASE" envDefault:"test"`
}

type ServerConfig struct {
	HttpPort              int           `env:"HTTP_PORT" envDefault:"8080"`
	HttpServerReadTimeout time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"5s"`
	GrpcPort              int           `env:"GRPC_PORT" envDefault:"9090"`
	MySql                 MysqlConfig   `envPrefix:"MYSQL_CONFIG_"`
}
