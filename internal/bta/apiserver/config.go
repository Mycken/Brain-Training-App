package apiserver

type Config struct {
	BingAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"logol_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BingAddr: ":8080",
		LogLevel: "debug",
	}
}
