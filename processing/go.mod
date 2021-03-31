module github.com/stasatdaglabs/kasboard/processing

go 1.15

require (
	github.com/go-pg/pg/v9 v9.1.3
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/kaspanet/kaspad v0.10.0
	github.com/pkg/errors v0.9.1
)

replace github.com/kaspanet/kaspad => ../../../kaspanet/kaspad
