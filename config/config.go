package config

type Config struct {
	TTLCheckInterval int
	Port             string
	DumpFile         string
	DumpInterval     int
}

func New(ttlCheck int, port string, dumpFile string, dumpInterval int) *Config {
	cfg := new(Config)
	cfg.TTLCheckInterval = ttlCheck
	cfg.Port = port
	cfg.DumpFile = dumpFile
	cfg.DumpInterval = dumpInterval

	return cfg
}
