package ut

type Config struct {
	RedisDSN string `yaml:"redis_dsn"`
	MysqlDSN string `yaml:"mysql_dsn"`
}
