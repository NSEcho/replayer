package config

type ReplayerConfig struct {
	Count         int
	Timeout       int
	Proxy         string
	PrintOnStdout bool
	PrintHeaders  bool
}
