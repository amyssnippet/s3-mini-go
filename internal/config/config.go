package config

type Config struct {
	Port int
	APIPort string
	StorePath string
	KeyPath string
}

func Default() *Config {
	return &Config{
		
	}
}