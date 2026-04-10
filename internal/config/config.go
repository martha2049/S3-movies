package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/env"
)

type Config struct {
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Secure    bool


	SyncTimeout time.Duration
}

func Load() (*Config, error) {
	k := koanf.New(".")

	// читаем ENV с префиксом S3_
	err := k.Load(env.Provider("S3_", ".", func(s string) string {
		return strings.TrimPrefix(s, "S3_")
	}), nil)
	if err != nil {
		return nil, err
	}

	return &Config{
		S3Endpoint:  k.String("ENDPOINT"),
		S3AccessKey: k.String("ACCESS_KEY"),
		S3SecretKey: k.String("SECRET_KEY"),
		S3Secure:    k.Bool("SECURE"),
	}, nil
}