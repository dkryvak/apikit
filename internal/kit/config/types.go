package config

type Hook func() error

type Config interface {
	Alias() string
	SetAlias(string)

	// BoundEnv returns the kube env this config is bound to (empty if none).
	// Used by 'apikit kube' to pick the cluster; ignored by 'apikit local'.
	BoundEnv() string

	Validate() error
	Hooks() []Hook
}
