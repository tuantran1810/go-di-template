package config

import (
	"time"
)

type SqliteConfig struct {
	DatabasePath string `env:"DATABASE_PATH" envDefault:":memory:"`
}

type AwsStorageConfig struct {
	KeyId          string        `env:"KEY_ID"`
	AppKey         string        `env:"APP_KEY"`
	Endpoint       string        `env:"ENDPOINT"`
	Region         string        `env:"REGION"`
	ForcePathStyle bool          `env:"FORCE_PATH_STYLE" envDefault:"false"`
	DefaultTimeout time.Duration `env:"DEFAULT_TIMEOUT" envDefault:"5s"`
}

type ServerConfig struct {
	HttpPort int `env:"HTTP_PORT" envDefault:"8080"`
	GrpcPort int `env:"GRPC_PORT" envDefault:"9090"`

	Sqlite     SqliteConfig     `envPrefix:"SQLITE_CONFIG_"`
	AwsStorage AwsStorageConfig `envPrefix:"AWS_STORAGE_"`

	JwtSecret  string        `env:"JWT_SECRET"`
	JwtExpired time.Duration `env:"JWT_EXPIRED" envDefault:"720h"`

	CdnHost                          string        `env:"CDN_HOST" envDefault:"https://cdn.com"`
	StorageBucket                    string        `env:"STORAGE_BUCKET" envDefault:"bucket"`
	StorageFolder                    string        `env:"STORAGE_LIVE_FOLDER" envDefault:"folder"`
	StorageUploadSignedUrlExpiration time.Duration `env:"STORAGE_UPLOAD_SIGNED_URL_EXPIRATION" envDefault:"10m"`
}
