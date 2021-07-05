module github.com/stasatdaglabs/kasboard/processing

go 1.16

require (
	github.com/go-pg/pg/v9 v9.1.3
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/jessevdk/go-flags v1.5.0
	github.com/kaspanet/kaspad v0.10.4
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	google.golang.org/genproto v0.0.0-20210406143921-e86de6bf7a46 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/kaspanet/kaspad => ../../../kaspanet/kaspad
