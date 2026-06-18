package console

import (
	"apikit/internal/kit/command"
	config2 "apikit/internal/kit/config"
	"apikit/internal/kit/console"
	"apikit/internal/softgamings/config"
	"bufio"
	"fmt"
	"os"
)

func PromptConfig(module command.ModuleName) (*config.SoftgamingsConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	var cfg config.SoftgamingsConfig

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

	key, err := console.Prompt(reader, "Key")
	if err != nil {
		return nil, err
	}
	cfg.Key = key

	pwd, err := console.Prompt(reader, "Pwd")
	if err != nil {
		return nil, err
	}
	cfg.Pwd = pwd

	hmacSecret, err := console.Prompt(reader, "HmacSecret")
	if err != nil {
		return nil, err
	}
	cfg.HmacSecret = hmacSecret

	return &cfg, nil
}

func PrintConfig(cfg *config.SoftgamingsConfig) {
	fmt.Printf("========================================\n")
	fmt.Printf("Config: %s\n", cfg.Alias())
	fmt.Printf("========================================\n")
	if cfg.Env != "" {
		fmt.Printf("Env: %s\n", cfg.Env)
	}
	fmt.Printf("Api Host: %s\n", cfg.ApiHost)
	fmt.Printf("Key: %s\n", cfg.Key)
	fmt.Printf("Pwd: %s\n", cfg.Pwd)
	fmt.Printf("HmacSecret: %s\n", cfg.HmacSecret)
	fmt.Printf("========================================\n")
}
