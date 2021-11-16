package main

import (
	"wxBot4g/internal/handler"
	"wxBot4g/internal/handler/middleware/event"
	"wxBot4g/wcbot"
)

var (
	Bot *wcbot.WcBot
)

func main() {
	Bot = wcbot.New()
	Bot.Debug = true
	Bot.QrCodeInTerminal() //默认在 wxqr 目录生成二维码，调用此函数，在终端打印二维码

	h := handler.NewWeChatBot(Bot, handler.WithMiddleware(event.Server))
	Bot.AddHandler(h)

	Bot.Run()
}
