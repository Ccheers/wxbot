package event

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wxBot4g/internal/biz"
	"wxBot4g/models"
	"wxBot4g/wcbot"

	"github.com/sirupsen/logrus"
)

const splitChar = " "

const (
	segTypeOnce = 2 // 一次性提醒
	segTypeLoop = 3 // 循环提醒
)

func NewEventServer(useCase *biz.JobUseCase) func(msg *models.RealRecvMsg, Bot *wcbot.WcBot) (isBreak bool) {
	return func(msg *models.RealRecvMsg, Bot *wcbot.WcBot) (isBreak bool) {
		job, ok := parseEvent(msg.Content.Data)
		if ok {
			return false
		}
		_, err := useCase.PutJob(context.TODO(), job)
		if err != nil {
			logrus.Error(err.Error())
			return false
		}

		return true
	}
}

func parseEvent(content string) (*biz.Job, bool) {
	seg := strings.Split(content, splitChar)
	if len(seg) != 2 && len(seg) != 3 {
		return nil, false
	}
	switch len(seg) {
	case segTypeOnce:
		tExp, err := parseTime(seg[0])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		return &biz.Job{
			CronExpress: fmt.Sprintf("%s %s", tExp, string(biz.ADay)),
			Content:     seg[1],
		}, true
	case segTypeLoop:
		fExp, err := parseFrequency(seg[0])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		tExp, err := parseTime(seg[1])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		return &biz.Job{
			CronExpress: fmt.Sprintf("%s %s", tExp, fExp),
			Origin:      "",
			Content:     seg[2],
			JobType:     0,
		}, true
	}
	return nil, false
}

var (
	regFreType = regexp.MustCompile("(每天|每日|每周|每月|工作日)[^\\Day]*(\\Day+)")
)
var (
	errParseFrequency = errors.New("parse Frequency error")
)

func parseFrequency(content string) (string, error) {
	var tp biz.FreType
	var d uint

	res := regFreType.FindAllStringSubmatch(content, -1)
	if len(res) > 0 {
		day, err := strconv.ParseUint(res[0][2], 10, 64)
		switch res[0][1] {
		case "每天":
			fallthrough
		case "每日":
			tp = biz.ADay
		case "每周":
			tp = biz.AWeek
		case "每月":
			tp = biz.AMonth
		case "工作日":
			tp = biz.WorkDay
		default:
			return "", fmt.Errorf("%w: %s", errParseFrequency, "no match")
		}
		// 如果是每周或者每月点类型的，则需要检查是否是合法的
		if (tp == biz.AMonth || tp == biz.AWeek) && err != nil {
			return "", fmt.Errorf("%w: %s", errParseFrequency, err)
		}
		d = uint(day)
		return fmt.Sprintf(string(tp), d), nil
	}
	return "", fmt.Errorf("%w: %s", errParseFrequency, "no match")
}

var errParseTime = errors.New("parse time error")

func parseTime(content string) (string, error) {
	t, err := time.ParseInLocation("15:04", content, time.Local)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errParseTime, err)
	}
	return t.Format("4 15"), nil
}
