package config

import (
	"github.com/jessevdk/go-flags"
	"github.com/kaspanet/kaspad/infrastructure/config"
	"github.com/pkg/errors"
)

type Config struct {
	RPCServerAddress         string `long:"rpc-server" description:"Kaspad RPC server to connect to. Should be of the form: <host>:<port>"`
	DatabaseConnectionString string `long:"connection-string" description:"Connection string for PostgrSQL database to connect to. Should be of the form: postgres://<username>:<password>@<host>:<port>/<database name>"`
	config.NetworkFlags
}

func Parse() (*Config, error) {
	config := &Config{}
	parser := flags.NewParser(config, flags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	if config.RPCServerAddress == "" {
		return nil, errors.Errorf("--rpc-server is required.")
	}
	if config.DatabaseConnectionString == "" {
		return nil, errors.Errorf("--connection-string is required.")
	}

	err = config.ResolveNetwork(parser)
	if err != nil {
		return nil, err
	}

	return config, nil
}
