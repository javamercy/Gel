package domain

type Config struct {
	User UserConfig `toml:"user"`
}

type UserConfig struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}
