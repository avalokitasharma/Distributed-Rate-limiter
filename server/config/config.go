package config

type Config struct {
	RedisHost  string
	RedisPort  string
	ServerPort string
}

func NewConfig() *Config {
	return &Config{
		RedisHost:  "localhost",
		RedisPort:  "6379",
		ServerPort: "808",
	}
}
