package config

import "time"

type MysqlConfig struct {
	Username string `env:"USERNAME" envDefault:"root"`
	Password string `env:"PASSWORD" envDefault:"secret"`
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Address  string `env:"ADDRESS" envDefault:"127.0.0.1:3306"`
	Database string `env:"DATABASE" envDefault:"test"`
}

type LoggingWorkerConfig struct {
	BufferCapacity int           `env:"BUFFER_CAPACITY" envDefault:"10"`
	FlushInterval  time.Duration `env:"FLUSH_INTERVAL" envDefault:"1s"`
}

type ServerConfig struct {
	HttpPort              int                 `env:"HTTP_PORT" envDefault:"8080"`
	HttpServerReadTimeout time.Duration       `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"5s"`
	GrpcPort              int                 `env:"GRPC_PORT" envDefault:"9090"`
	MySql                 MysqlConfig         `envPrefix:"MYSQL_CONFIG_"`
	LoggingWorker         LoggingWorkerConfig `envPrefix:"LOGGING_WORKER_CONFIG_"`
}
