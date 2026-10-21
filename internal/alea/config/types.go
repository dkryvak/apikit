package config

import (
	"apikit/internal/kit/config"
	"fmt"
)

// AleaConfig holds credentials and settings for the Alea APIs.
// Alea exposes several hosts (customer GraphQL API, bonus REST API, launch API);
// these values populate the request headers and the shared GraphQL variables.
type AleaConfig struct {
	alias string

	Env string `json:"env,omitempty"` // bound kube env (required for 'apikit kube', ignored for 'apikit local')

	CustomerApiHost string `json:"customerApiHost"`         // GraphQL customer API (e.g. https://customer-api.aleaplay.com)
	BonusApiHost    string `json:"bonusApiHost"`            // bonus REST API (e.g. https://bonus-api.aleaplay.com)
	LaunchApiHost   string `json:"launchApiHost,omitempty"` // game launch API (e.g. https://play.aleaplay.com)

	EnvironmentToken string `json:"environmentToken"`        // Bearer JWT (Authorization header)
	SecretKey        string `json:"secretKey"`               // Alea secret key (Digest signature; per brand/environment)
	CasinoId         string `json:"casinoId"`                // Alea-CasinoId header
	EnvironmentId    string `json:"environmentId,omitempty"` // Alea-EnvironmentId header (optional)
	JurisdictionCode string `json:"jurisdictionCode"`        // GraphQL variable injected into queries
}

func New() *AleaConfig {
	return &AleaConfig{}
}

func (cfg *AleaConfig) SetAlias(alias string) {
	cfg.alias = alias
}

func (cfg *AleaConfig) Alias() string {
	return cfg.alias
}

func (cfg *AleaConfig) BoundEnv() string {
	return cfg.Env
}

func (cfg *AleaConfig) Validate() error {
	if cfg.CustomerApiHost == "" {
		return fmt.Errorf("customerApiHost is required")
	}
	if cfg.BonusApiHost == "" {
		return fmt.Errorf("bonusApiHost is required")
	}
	if cfg.EnvironmentToken == "" {
		return fmt.Errorf("environmentToken is required")
	}
	if cfg.SecretKey == "" {
		return fmt.Errorf("secretKey is required")
	}
	if cfg.CasinoId == "" {
		return fmt.Errorf("casinoId is required")
	}
	if cfg.JurisdictionCode == "" {
		return fmt.Errorf("jurisdictionCode is required")
	}
	return nil
}

func (cfg *AleaConfig) Hooks() []config.Hook {
	return []config.Hook{}
}
