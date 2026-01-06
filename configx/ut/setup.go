package ut

import (
	"testing"

	"github.com/GizmoVault/gotools/configx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type Flag int

const (
	FlagMySQL Flag = 1 << iota
	FlagRedis
)

const (
	defaultUTConfig = "ut.yaml"
)

func SetupUTConfig() *Config {
	return SetupUTConfigEx(defaultUTConfig, nil)
}

func SetupUTConfigEx(fileName string, configPaths []string) *Config {
	cfg := &Config{}

	var err error

	if len(configPaths) == 0 {
		_, err = configx.Load(fileName, cfg)
	} else {
		_, err = configx.LoadOnConfigPath(fileName, configPaths, cfg)
	}

	if err != nil {
		return nil
	}

	return cfg
}

func SetupAndCheckFlags(t *testing.T, flags Flag) *Config {
	return SetupAndCheckFlagsEx(t, defaultUTConfig, nil, flags)
}

func SetupAndCheckFlagsEx(t *testing.T, fileName string, configPaths []string, flags Flag) *Config {
	cfg := SetupUTConfigEx(fileName, configPaths)
	if cfg == nil {
		t.SkipNow()

		return nil
	}

	if flags&FlagRedis == FlagRedis {
		if cfg.RedisDSN == "" {
			t.SkipNow()

			return nil
		}

		options, err := redis.ParseURL(cfg.RedisDSN)
		require.NoError(t, err)

		cfg.RedisOpt = options
	}

	if flags&FlagMySQL == FlagMySQL {
		if cfg.MysqlDSN == "" {
			t.SkipNow()

			return nil
		}
	}

	return cfg
}
