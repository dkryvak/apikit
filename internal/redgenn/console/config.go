package console

import (
	"apikit/internal/kit/command"
	config2 "apikit/internal/kit/config"
	"apikit/internal/kit/console"
	"apikit/internal/redgenn/config"
	"bufio"
	"fmt"
	"os"
)

func PromptConfig(module command.ModuleName) (*config.RedgennConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	var cfg config.RedgennConfig

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

	apiHost, err := console.Prompt(reader, "Api Host (e.g., 'https://api.example.com')")
	if err != nil {
		return nil, err
	}
	cfg.ApiHost = apiHost

	launchApiHost, err := console.Prompt(reader, "Launch Api Host (e.g., 'https://api.example.com')")
	if err != nil {
		return nil, err
	}
	cfg.LaunchApiHost = launchApiHost

	login, err := console.Prompt(reader, "Login")
	if err != nil {
		return nil, err
	}
	cfg.Login = login

	password, err := console.Prompt(reader, "Password")
	if err != nil {
		return nil, err
	}
	cfg.Password = password

	partnerCode, err := console.Prompt(reader, "Partner Code")
	if err != nil {
		return nil, err
	}
	cfg.PartnerCode = partnerCode

	wlCode, err := console.Prompt(reader, "WL Code")
	if err != nil {
		return nil, err
	}
	cfg.WlCode = wlCode

	return &cfg, nil
}

func PrintConfig(cfg *config.RedgennConfig) {
	fmt.Printf("========================================\n")
	fmt.Printf("Config: %s\n", cfg.Alias())
	fmt.Printf("========================================\n")
	if cfg.Env != "" {
		fmt.Printf("Env: %s\n", cfg.Env)
	}
	fmt.Printf("Api Host: %s\n", cfg.ApiHost)
	fmt.Printf("Launch Api Host: %s\n", cfg.LaunchApiHost)
	fmt.Printf("Login: %s\n", cfg.Login)
	fmt.Printf("Password: %s\n", cfg.Password)
	fmt.Printf("Partner Code: %s\n", cfg.PartnerCode)
	fmt.Printf("WL Code: %s\n", cfg.WlCode)
	fmt.Printf("========================================\n")
}
