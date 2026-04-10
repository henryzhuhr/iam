// Package config provides configuration structures for the IAM application.
package config

import (
	"iam/infra/cache"
	"iam/infra/database"
	"iam/infra/queue"

	"github.com/zeromicro/go-zero/rest"
)

// Config is the top-level configuration for the IAM application.
type Config struct {
	rest.RestConf
	LocaleDir string                 `json:"LocaleDir,optional"`
	DB        database.MySQLConfig    `json:"DB"`
	Redis     cache.RedisConfig       `json:"Redis"`
	Kafka     queue.KafkaConfig       `json:"Kafka,optional"`
	Pprof     PprofConfig             `json:"Pprof,optional"`
}

// PprofConfig pprof performance analysis configuration.
type PprofConfig struct {
	Enabled bool `json:"Enabled,default=false"` // whether to enable pprof
	Port    int  `json:"Port,default=6060"`     // pprof service port
}
