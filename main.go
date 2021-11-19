package main

import (
	"wxBot4g/internal/biz"
	"wxBot4g/internal/handler"
	"wxBot4g/internal/handler/middleware/event"
	"wxBot4g/wcbot"
)

var (
	bot *wcbot.WcBot
)

func newHandler(Bot *wcbot.WcBot, jobUseCase *biz.JobUseCase) *handler.WeChatBot {
	return handler.NewWeChatBot(
		Bot,
		handler.WithMiddleware(
			event.NewEventServer(jobUseCase),
		),
	)
}

func main() {
	bot = wcbot.New()
	bot.Debug = true
	bot.QrCodeInTerminal() //默认在 wxqr 目录生成二维码，调用此函数，在终端打印二维码
	h := initBot(bot)
	bot.AddHandler(h)

	bot.Run()
}
