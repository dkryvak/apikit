package kube

import (
	"apikit/internal/config"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// EnsureAWSSSO makes sure there is a valid AWS session. client-go itself never
// triggers an interactive login, so this is the single external shell-out:
// if the session is missing/expired we run `aws sso login` interactively.
//
// Only `apikit kube` calls this — `apikit local` never needs AWS.
// Set APIKIT_SKIP_SSO=1 to skip the check entirely (e.g. when credentials are
// managed some other way than SSO).
func EnsureAWSSSO(ctx context.Context) error {
	if config.SkipSSO() {
		return nil
	}

	if _, err := exec.LookPath("aws"); err != nil {
		return fmt.Errorf(
			"AWS CLI not found on PATH.\n" +
				"'apikit kube' needs the AWS CLI plus a configured SSO profile and cluster access (kubeconfig).\n" +
				"  • Install the AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html\n" +
				"  • 'apikit local' works without AWS.\n" +
				"  • If you manage AWS credentials another way, set APIKIT_SKIP_SSO=1 to skip this check",
		)
	}

	if err := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity").Run(); err == nil {
		return nil
	}

	fmt.Fprintln(os.Stderr, "AWS session expired or not found. Logging in via SSO...")

	login := exec.CommandContext(ctx, "aws", "sso", "login")
	login.Stdin = os.Stdin
	login.Stdout = os.Stderr
	login.Stderr = os.Stderr
	if err := login.Run(); err != nil {
		return fmt.Errorf("aws sso login failed: %w", err)
	}
	return nil
}
