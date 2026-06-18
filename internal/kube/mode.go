package kube

import (
	"os"
	"strings"
)

// Image is the fixed apikit container image used by the in-cluster Job.
// It MUST contain the apikit binary, so it is a constant — never configurable
// per-env (see docs/remote-design.md).
const Image = "dkryvak/apikit:latest"

// IsRemote reports whether the current invocation runs under the `kube`
// namespace (i.e. it must self-remote). Detection is based on the first
// command token in os.Args, which is framework-independent and recursion-safe:
// inside the pod we invoke `apikit local ...`, so IsRemote returns false there.
func IsRemote() bool {
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return a == "kube"
	}
	return false
}
