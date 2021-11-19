// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/Ccheers/wxbot/internal/biz"
	"github.com/Ccheers/wxbot/internal/cron"
	"github.com/Ccheers/wxbot/internal/data"
	"github.com/Ccheers/wxbot/internal/handler"
	"github.com/Ccheers/wxbot/wcbot"

	"github.com/google/wire"
)

func initBot(bot *wcbot.WcBot) *handler.WeChatBot {
	panic(wire.Build(cron.Provider, data.Provider, biz.Provider, newHandler))
}
