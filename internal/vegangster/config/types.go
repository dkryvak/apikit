package config

import (
	"apikit/internal/kit/config"
	"apikit/internal/vegangster/crypto"
	"crypto/rsa"
	"fmt"
)

type VegangsterConfig struct {
	alias string

	Env string `json:"env,omitempty"` // bound kube env (required for 'apikit kube', ignored for 'apikit local')

	ApiHost       string `json:"apiHost"`
	OperatorId    string `json:"operatorId"`
	BrandId       string `json:"brandId"`
	PrivateKeyStr string `json:"privateKey"`
	PublicKeyStr  string `json:"publicKey"`

	PrivateKey *rsa.PrivateKey `json:"-"`
	PublicKey  *rsa.PublicKey  `json:"-"`
}

func New() *VegangsterConfig {
	return &VegangsterConfig{}
}

func (cfg *VegangsterConfig) SetAlias(alias string) {
	cfg.alias = alias
}

func (cfg *VegangsterConfig) Alias() string {
	return cfg.alias
}

func (cfg *VegangsterConfig) BoundEnv() string {
	return cfg.Env
}

func (cfg *VegangsterConfig) Validate() error {
	if cfg.ApiHost == "" {
		return fmt.Errorf("apiHost is required")
	}
	if cfg.OperatorId == "" {
		return fmt.Errorf("operatorId is required")
	}
	if cfg.BrandId == "" {
		return fmt.Errorf("brandId is required")
	}
	if cfg.PrivateKeyStr == "" {
		return fmt.Errorf("privateKey is required")
	}
	if cfg.PublicKeyStr == "" {
		return fmt.Errorf("publicKey is required")
	}
	return nil
}

func (cfg *VegangsterConfig) Hooks() []config.Hook {
	return []config.Hook{
		func() error { return cfg.parseKeys() },
	}
}

func (cfg *VegangsterConfig) parseKeys() error {
	if cfg.PrivateKeyStr != "" {
		privateKey, err := crypto.ParsePrivateKey(cfg.PrivateKeyStr)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		cfg.PrivateKey = privateKey
	}

	if cfg.PublicKeyStr != "" {
		publicKey, err := crypto.ParsePublicKey(cfg.PublicKeyStr)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}
		cfg.PublicKey = publicKey
	}

	return nil
}
