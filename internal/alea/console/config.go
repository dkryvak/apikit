package console

import (
	"apikit/internal/alea/config"
	"apikit/internal/kit/command"
	config2 "apikit/internal/kit/config"
	"apikit/internal/kit/console"
	"bufio"
	"fmt"
	"os"
)

func PromptConfig(module command.ModuleName) (*config.AleaConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	var cfg config.AleaConfig

	alias, err := console.Prompt(reader, "Config alias (e.g., 'myapi', 'prod', 'staging')")
	if err != nil {
		return nil, err
	}
	cfg.SetAlias(alias)

	if config2.Exists(alias, module) {
		overwrite, err := console.PromptYesNo(reader, fmt.Sprintf("Config '%s' already exists. Overwrite?", alias))
		if err != nil {
			return nil, err
		}
		if !overwrite {
			return nil, fmt.Errorf("config creation cancelled")
		}
	}

	envName, err := console.PromptOptional(reader, "Bind to env name (optional; required for 'apikit kube ...')")
	if err != nil {
		return nil, err
	}
	cfg.Env = envName

	customerApiHost, err := console.Prompt(reader, "Customer Api Host (e.g., 'https://customer-api.aleaplay.com')")
	if err != nil {
		return nil, err
	}
	cfg.CustomerApiHost = customerApiHost

	bonusApiHost, err := console.Prompt(reader, "Bonus Api Host (e.g., 'https://bonus-api.aleaplay.com')")
	if err != nil {
		return nil, err
	}
	cfg.BonusApiHost = bonusApiHost

	launchApiHost, err := console.PromptOptional(reader, "Launch Api Host (optional; e.g., 'https://play.aleaplay.com')")
	if err != nil {
		return nil, err
	}
	cfg.LaunchApiHost = launchApiHost

	environmentToken, err := console.Prompt(reader, "Environment Token (Bearer JWT)")
	if err != nil {
		return nil, err
	}
	cfg.EnvironmentToken = environmentToken

	secretKey, err := console.Prompt(reader, "Secret Key")
	if err != nil {
		return nil, err
	}
	cfg.SecretKey = secretKey

	casinoId, err := console.Prompt(reader, "Casino Id (Alea-CasinoId)")
	if err != nil {
		return nil, err
	}
	cfg.CasinoId = casinoId

	environmentId, err := console.PromptOptional(reader, "Environment Id (Alea-EnvironmentId, optional)")
	if err != nil {
		return nil, err
	}
	cfg.EnvironmentId = environmentId

	jurisdictionCode, err := console.Prompt(reader, "Jurisdiction Code (e.g., 'CAO')")
	if err != nil {
		return nil, err
	}
	cfg.JurisdictionCode = jurisdictionCode

	return &cfg, nil
}

func PrintConfig(cfg *config.AleaConfig) {
	fmt.Printf("========================================\n")
	fmt.Printf("Config: %s\n", cfg.Alias())
	fmt.Printf("========================================\n")
	if cfg.Env != "" {
		fmt.Printf("Env: %s\n", cfg.Env)
	}
	fmt.Printf("Customer Api Host: %s\n", cfg.CustomerApiHost)
	fmt.Printf("Bonus Api Host: %s\n", cfg.BonusApiHost)
	if cfg.LaunchApiHost != "" {
		fmt.Printf("Launch Api Host: %s\n", cfg.LaunchApiHost)
	}
	fmt.Printf("Environment Token: %s\n", cfg.EnvironmentToken)
	fmt.Printf("Secret Key: %s\n", cfg.SecretKey)
	fmt.Printf("Casino Id: %s\n", cfg.CasinoId)
	if cfg.EnvironmentId != "" {
		fmt.Printf("Environment Id: %s\n", cfg.EnvironmentId)
	}
	fmt.Printf("Jurisdiction Code: %s\n", cfg.JurisdictionCode)
	fmt.Printf("========================================\n")
}
