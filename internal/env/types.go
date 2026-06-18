package env

import "fmt"

// Env describes a Kubernetes target used by `apikit kube ...`.
// It is a shared, mode-neutral entity (see docs/remote-design.md):
// only `kube` execution consults it; `local` ignores it entirely.
//
// Note: the container image is intentionally NOT part of env — it is a fixed
// constant (kube.Image), because the image MUST contain the apikit binary.
type Env struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Context   string `json:"context"`
}

func New() *Env {
	return &Env{}
}

// ApplyDefaults fills derivable fields when they are empty.
//   - Namespace defaults to "<name>-casino"
//
// Context is left empty when unset (resolved against the local kubeconfig later).
func (e *Env) ApplyDefaults() {
	if e.Namespace == "" && e.Name != "" {
		e.Namespace = e.Name + "-casino"
	}
}

func (e *Env) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	if e.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}
