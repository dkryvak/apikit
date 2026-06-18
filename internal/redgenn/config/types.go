package config

import (
	"apikit/internal/kit/config"
	"fmt"
)

type RedgennConfig struct {
	alias string

	Env string `json:"env,omitempty"` // bound kube env (required for 'apikit kube', ignored for 'apikit local')

	ApiHost       string `json:"apiHost"`
	LaunchApiHost string `json:"launchApiHost"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	PartnerCode   string `json:"partnerCode"`
	WlCode        string `json:"wlCode"`
}

func New() *RedgennConfig {
	return &RedgennConfig{}
}

func (cfg *RedgennConfig) SetAlias(alias string) {
	cfg.alias = alias
}

func (cfg *RedgennConfig) Alias() string {
	return cfg.alias
}

func (cfg *RedgennConfig) BoundEnv() string {
	return cfg.Env
}

func (cfg *RedgennConfig) Validate() error {
	if cfg.ApiHost == "" {
		return fmt.Errorf("apiHost is required")
	}
	if cfg.LaunchApiHost == "" {
		return fmt.Errorf("launchApiHost is required")
	}
	if cfg.Login == "" {
		return fmt.Errorf("login is required")
	}
	if cfg.Password == "" {
		return fmt.Errorf("password is required")
	}
	if cfg.PartnerCode == "" {
		return fmt.Errorf("partnerCode is required")
	}
	if cfg.WlCode == "" {
		return fmt.Errorf("wlCode is required")
	}
	return nil
}

func (cfg *RedgennConfig) Hooks() []config.Hook {
	return []config.Hook{}
}
