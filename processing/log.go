package main

import (
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kashboard/processing/logging"
)

var (
	log   = logging.Backend().Logger("KSBD")
	spawn = panics.GoroutineWrapperFunc(log)
)
