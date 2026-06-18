package config

import (
	"apikit/internal/kit/config"
	"fmt"
)

type SoftgamingsConfig struct {
	alias string

	Env string `json:"env,omitempty"` // bound kube env (required for 'apikit kube', ignored for 'apikit local')

	ApiHost    string `json:"apiHost"`
	Key        string `json:"key"`
	Pwd        string `json:"pwd"`
	HmacSecret string `json:"hmacSecret"`
}

func New() *SoftgamingsConfig {
	return &SoftgamingsConfig{}
}

func (cfg *SoftgamingsConfig) SetAlias(alias string) {
	cfg.alias = alias
}

func (cfg *SoftgamingsConfig) Alias() string {
	return cfg.alias
}

func (cfg *SoftgamingsConfig) BoundEnv() string {
	return cfg.Env
}

func (cfg *SoftgamingsConfig) Validate() error {
	if cfg.ApiHost == "" {
		return fmt.Errorf("apiHost is required")
	}
	if cfg.Key == "" {
		return fmt.Errorf("key is required")
	}
	if cfg.Pwd == "" {
		return fmt.Errorf("pwd is required")
	}
	if cfg.HmacSecret == "" {
		return fmt.Errorf("hmacSecret is required")
	}
	return nil
}

func (cfg *SoftgamingsConfig) Hooks() []config.Hook {
	return []config.Hook{}
}
