module github.com/stasatdaglabs/kasboard/processing

go 1.16

require (
	github.com/containerd/containerd v1.4.12 // indirect
	github.com/go-pg/pg/v9 v9.1.3
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/jessevdk/go-flags v1.5.0
	github.com/kaspanet/kaspad v0.10.4
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
)

replace github.com/kaspanet/kaspad => ../../../kaspanet/kaspad
