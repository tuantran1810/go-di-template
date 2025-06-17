package config

type SqliteConfig struct {
	DatabasePath string `env:"DATABASE_PATH" envDefault:":memory:"`
}

type FakeConsumerConfig struct {
	PerMs uint `env:"PER_MS" envDefault:"1000"`
}

type FakeClientConfig struct {
	LatencyMs uint `env:"LATENCY_MS" envDefault:"1000"`
}

type ServerConfig struct {
	HttpPort int `env:"HTTP_PORT" envDefault:"8080"`
	GrpcPort int `env:"GRPC_PORT" envDefault:"9090"`

	Sqlite       SqliteConfig       `envPrefix:"SQLITE_CONFIG_"`
	FakeConsumer FakeConsumerConfig `envPrefix:"FAKE_CONSUMER_CONFIG_"`
	FakeClient   FakeClientConfig   `envPrefix:"FAKE_CLIENT_CONFIG_"`
}
