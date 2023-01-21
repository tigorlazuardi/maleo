package maleozap_test

import (
	"go.uber.org/zap"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/maleozap"
)

func ExampleNew_production() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	mlog := maleozap.New(logger)
	mlog.Sync()
	maleo.Global().SetLogger(mlog)
}
