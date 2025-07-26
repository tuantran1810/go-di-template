package config

import (
	"time"
)

type PubsubConfig struct {
	ProjectId      string        `env:"PROJECT_ID" envDefault:"project-id"`
	TopicId        string        `env:"TOPIC_ID" envDefault:"default"`
	SubscriptionId string        `env:"SUBSCRIPTION_ID" envDefault:"default"`
	Region         string        `env:"REGION" envDefault:"us-central1"`
	Endpoint       string        `env:"ENDPOINT" envDefault:"https://pubsub.googleapis.com"`
	KeyFile        string        `env:"KEY_FILE" envDefault:"/path/to/key.json"`
	DefaultTimeout time.Duration `env:"DEFAULT_TIMEOUT" envDefault:"5s"`
}

type ConsumerConfig struct {
	HttpPort int `env:"HTTP_PORT" envDefault:"8080"`
	GrpcPort int `env:"GRPC_PORT" envDefault:"9090"`

	Pubsub    PubsubConfig `envPrefix:"PUBSUB_"`
	LogSqlite SqliteConfig `envPrefix:"LOG_SQLITE_CONFIG_"`

	CdnHost                          string        `env:"CDN_HOST" envDefault:"https://cdn.com"`
	StorageBucket                    string        `env:"STORAGE_BUCKET" envDefault:"bucket"`
	StorageFolder                    string        `env:"STORAGE_LIVE_FOLDER" envDefault:"folder"`
	StorageUploadSignedUrlExpiration time.Duration `env:"STORAGE_UPLOAD_SIGNED_URL_EXPIRATION" envDefault:"10m"`
}
