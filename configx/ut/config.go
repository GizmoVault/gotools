package ut

import "github.com/redis/go-redis/v9"

type Config struct {
	RedisDSN string `yaml:"redis_dsn"`
	MysqlDSN string `yaml:"mysql_dsn"`

	//
	// redis
	//
	RedisOpt *redis.Options `yaml:"-"`
}
