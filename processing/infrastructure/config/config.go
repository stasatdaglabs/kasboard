package config

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type Config struct {
	RPCServerAddress string `long:"rpcserver" description:"Kaspad RPC server to connect to"`
}

func Parse() (*Config, error) {
	config := &Config{}
	parser := flags.NewParser(config, flags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	if config.RPCServerAddress == "" {
		return nil, errors.Errorf("--rpcserver is required.")
	}

	return config, nil
}
