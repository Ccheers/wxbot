// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"wxBot4g/internal/biz"
	"wxBot4g/internal/cron"
	"wxBot4g/internal/data"
	"wxBot4g/internal/handler"
	"wxBot4g/wcbot"
)

// Injectors from wire.go:

func initBot(bot2 *wcbot.WcBot) *handler.WeChatBot {
	dataData := data.NewData()
	jobRepo := data.NewJobRepoImpl(dataData)
	cronCron := cron.NewCron()
	jobUseCase := biz.NewJobUseCase(jobRepo, cronCron)
	weChatBot := newHandler(bot2, jobUseCase, cronCron)
	return weChatBot
}
