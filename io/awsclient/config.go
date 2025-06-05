package awsclient

import "time"

type ClientConfig struct {
	ReadTimeoutDuration time.Duration
	ListTimeoutDuration time.Duration
	MaxReadRetries      int
	MaxListRetries      int
	MaxPutRetries       int
}

var defaultClientCfg = ClientConfig{
	ReadTimeoutDuration: 500 * time.Millisecond,
	ListTimeoutDuration: 500 * time.Millisecond,
	MaxReadRetries:      3,
	MaxListRetries:      3,
	MaxPutRetries:       3,
}

var clientConfig = defaultClientCfg

func GetClientConfig() ClientConfig {
	return clientConfig
}

func SetClientConfig(config ClientConfig) {
	clientConfig = config
}
