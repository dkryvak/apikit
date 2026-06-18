package console

import (
	"apikit/internal/kit/command"
	config2 "apikit/internal/kit/config"
	"apikit/internal/kit/console"
	"apikit/internal/vegangster/config"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PromptConfig(module command.ModuleName) (*config.VegangsterConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	var cfg config.VegangsterConfig

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

	operatorId, err := console.Prompt(reader, "Operator ID")
	if err != nil {
		return nil, err
	}
	cfg.OperatorId = operatorId

	brandId, err := console.Prompt(reader, "Brand ID")
	if err != nil {
		return nil, err
	}
	cfg.BrandId = brandId

	fmt.Println("\nPrivate Key (base64-encoded PKCS#8 format):")
	privateKey, err := console.PromptMultiline(reader, "Paste your private key (press Enter twice when done)")
	if err != nil {
		return nil, err
	}
	cfg.PrivateKeyStr = strings.TrimSpace(privateKey)

	fmt.Println("\nPublic Key (base64-encoded X.509 format):")
	publicKey, err := console.PromptMultiline(reader, "Paste your public key (press Enter twice when done)")
	if err != nil {
		return nil, err
	}
	cfg.PublicKeyStr = strings.TrimSpace(publicKey)

	return &cfg, nil
}

func PrintConfig(cfg *config.VegangsterConfig) {
	fmt.Printf("========================================\n")
	fmt.Printf("Config: %s\n", cfg.Alias())
	fmt.Printf("========================================\n")
	if cfg.Env != "" {
		fmt.Printf("Env: %s\n", cfg.Env)
	}
	fmt.Printf("Api Host: %s\n", cfg.ApiHost)
	fmt.Printf("Operator ID: %s\n", cfg.OperatorId)
	fmt.Printf("Brand ID: %s\n", cfg.BrandId)

	if cfg.PrivateKeyStr != "" {
		keyPreview := cfg.PrivateKeyStr
		if len(keyPreview) > 50 {
			keyPreview = keyPreview[:50] + "..."
		}
		fmt.Printf("Private Key: %s\n", keyPreview)

		if cfg.PrivateKey != nil {
			fmt.Printf("  Type: RSA\n")
			fmt.Printf("  Size: %d bits\n", cfg.PrivateKey.N.BitLen())
		}
	}

	if cfg.PublicKeyStr != "" {
		keyPreview := cfg.PublicKeyStr
		if len(keyPreview) > 50 {
			keyPreview = keyPreview[:50] + "..."
		}
		fmt.Printf("Public Key: %s\n", keyPreview)

		if cfg.PublicKey != nil {
			fmt.Printf("  Type: RSA\n")
			fmt.Printf("  Size: %d bits\n", cfg.PublicKey.N.BitLen())
		}
	}
	fmt.Printf("========================================\n")
}
