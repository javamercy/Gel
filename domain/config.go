package domain

import "Gel/core/validation"

type UserConfig struct {
	Name  string `toml:"name,omitempty" validate:"required,min=1,max=256"`
	Email string `toml:"email,omitempty" validate:"required,email"`
}

func NewUserConfig(name, email string) (UserConfig, error) {
	userConfig := UserConfig{
		Name:  name,
		Email: email,
	}

	validator := validation.GetValidator()
	if err := validator.Struct(userConfig); err != nil {
		return UserConfig{}, err
	}
	return userConfig, nil
}

type Config struct {
	User UserConfig `toml:"user" validate:"required"`
}

func NewConfig(user UserConfig) (*Config, error) {
	config := &Config{
		User: user,
	}

	validator := validation.GetValidator()
	if err := validator.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
