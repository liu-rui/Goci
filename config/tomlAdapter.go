package config

import "github.com/BurntSushi/toml"

type tomlAdapter struct {
}

func (config *tomlAdapter) Decode(v interface{}) error {
	_, err := toml.DecodeFile("config.toml", v)

	return err
}

func newTomlAdapter() *tomlAdapter {
	return &tomlAdapter{}
}
