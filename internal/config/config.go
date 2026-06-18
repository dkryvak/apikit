// Package config holds build-time configuration. All values are baked into the
// binary at build via:
//
//	-ldflags "-X apikit/internal/config.<name>=<value>"
//
// driven by .env (see Taskfile) or the release tag in CI. There is no runtime
// env lookup for these — to change a value, rebuild. The only exception is
// SkipSSO, which is a genuine per-run toggle.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Build-time values (injected via -ldflags; defaults below apply to a plain
// `go build` without flags).
var (
	version             = "dev"
	image               = "" // no default — see Image()
	podLifetimeSeconds  = "1800"
	jobTTLSeconds       = "120"
	podReadyTimeoutSecs = "90"
)

// Version is the build version for `apikit --version`.
func Version() string { return version }

// Image is the apikit pod image baked at build. Errors if it was not set
// (a binary without an image cannot start a pod).
func Image() (string, error) {
	if v := strings.TrimSpace(image); v != "" {
		return v, nil
	}
	return "", fmt.Errorf(`pod image not set: build with -ldflags "-X apikit/internal/config.image=..."`)
}

// PodLifetimeSeconds is the pod sleep duration AND the Job hard cap.
func PodLifetimeSeconds() int { return atoiOr(podLifetimeSeconds, 1800) }

// JobTTLSeconds is TTLSecondsAfterFinished (cleanup after the Job finishes).
func JobTTLSeconds() int { return atoiOr(jobTTLSeconds, 120) }

// PodReadyTimeoutSeconds is the wait-for-pod-ready timeout.
func PodReadyTimeoutSeconds() int { return atoiOr(podReadyTimeoutSecs, 90) }

// SkipSSO is a runtime toggle (NOT baked): set APIKIT_SKIP_SSO=1 in the shell to
// skip the AWS SSO preflight when running `apikit kube`.
func SkipSSO() bool { return strings.TrimSpace(os.Getenv("APIKIT_SKIP_SSO")) != "" }

func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && n > 0 {
		return n
	}
	return def
}
