package config

type Config struct {
	TTLCheckInterval int
	Port             string
}

func New(ttlCheck int, port string) *Config {
	cfg := new(Config)
	cfg.TTLCheckInterval = ttlCheck
	cfg.Port = port
	return cfg
}
