package main

import (
	"context"
	"wxBot4g/internal/biz"
	"wxBot4g/internal/cron"
	"wxBot4g/internal/handler"
	"wxBot4g/internal/handler/middleware/event"
	"wxBot4g/wcbot"

	"github.com/sirupsen/logrus"
)

var (
	bot *wcbot.WcBot
)

func newHandler(bot *wcbot.WcBot, jobUseCase *biz.JobUseCase, c *cron.Cron) *handler.WeChatBot {
	h := handler.NewWeChatBot(
		bot,
		handler.WithMiddleware(
			event.NewEventServer(jobUseCase, bot),
		),
	)

	jobs, err := jobUseCase.GetAllJobs(context.Background())
	if err != nil {
		panic(err)
	}

	for _, job := range jobs {
		cid, err := c.AddCron(job.CronExpress, jobUseCase.WithCronFunc(job.ID))
		if err != nil {
			logrus.Error(err)
			continue
		}
		job.CronID = cid
		_, err = jobUseCase.UpdateJob(context.Background(), job)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Infof("add cron %+v", job)
	}

	return h
}

func main() {
	bot = wcbot.New()
	bot.Debug = true
	bot.QrCodeInTerminal() //默认在 wxqr 目录生成二维码，调用此函数，在终端打印二维码
	h := initBot(bot)
	bot.AddHandler(h)

	bot.Run()
}
