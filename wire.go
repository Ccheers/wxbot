// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"wxBot4g/internal/biz"
	"wxBot4g/internal/cron"
	"wxBot4g/internal/data"
	"wxBot4g/internal/handler"
	"wxBot4g/wcbot"

	"github.com/google/wire"
)

func initBot(bot *wcbot.WcBot) *handler.WeChatBot {
	panic(wire.Build(cron.Provider, data.Provider, biz.Provider, newHandler))
}
