// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"wxbot/internal/biz"
	"wxbot/internal/cron"
	"wxbot/internal/data"
	"wxbot/internal/handler"
	"wxbot/wcbot"

	"github.com/google/wire"
)

func initBot(bot *wcbot.WcBot) *handler.WeChatBot {
	panic(wire.Build(cron.Provider, data.Provider, biz.Provider, newHandler))
}
