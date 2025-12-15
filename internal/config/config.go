package config

type Config struct {
	Port int
	APIPort string
	StorePath string
	KeyPath string
}

func Default() *Config {
	return &Config{
		Port: 9000,
		APIPort: ":6125",
		StorePath: "./my_files",
		KeyPath: "./keys",
	}
}