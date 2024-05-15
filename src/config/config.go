package config


type Config struct {
	ListenAddr string
}

func Load(configPath string) *Config {
	return &Config{
		ListenAddr: ":8080",
	}
}