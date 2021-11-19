package biz

import "github.com/google/wire"

var Provider = wire.NewSet(NewJobUseCase)
