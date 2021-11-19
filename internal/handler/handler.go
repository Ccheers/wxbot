package handler

import (
	"wxBot4g/models"
	"wxBot4g/pkg/define"
	"wxBot4g/wcbot"

	"github.com/sirupsen/logrus"
)

type MsgHandler func(msg *models.RealRecvMsg, Bot *wcbot.WcBot) (isBreak bool)

type BotOption func(bot *WeChatBot)

type WeChatBot struct {
	Bot         *wcbot.WcBot
	middlewares []MsgHandler
}

func NewWeChatBot(bot *wcbot.WcBot, opts ...BotOption) *WeChatBot {
	ins := &WeChatBot{Bot: bot}
	for _, opt := range opts {
		opt(ins)
	}
	return ins
}

func (w *WeChatBot) HandleMessage(msg *models.RealRecvMsg) {
	//过滤不支持消息99
	if msg.MsgType == 99 || msg.MsgTypeId == 99 {
		return
	}

	//获取unknown的username
	contentUser := msg.Content.User.Name

	logrus.Debug(
		"消息类型:", define.MsgIdString(msg.MsgTypeId), " ",
		"数据类型:", define.MsgTypeIdString(msg.Content.Type), " ",
		"发送者:", msg.FromUserName, " ",
		"发送人:", msg.SendMsgUSer.Name, " ",
		"发送内容人:", contentUser, " ",
		"内容:", msg.Content.Data,
	)

	for _, h := range w.middlewares {
		isBreak := h(msg, w.Bot)
		if isBreak {
			break
		}
	}
}

func WithMiddleware(middlewares ...MsgHandler) func(bot *WeChatBot) {
	return func(bot *WeChatBot) {
		bot.middlewares = middlewares
	}
}
