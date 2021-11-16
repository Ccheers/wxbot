package event

import (
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

type freType uint

const (
	none    freType = 0
	aDay    freType = 1
	aWeek   freType = 1 << 2
	aMonth  freType = 1 << 3
	workDay freType = 1 << 4
)

type event struct {
	frequency *frequency
	when      time.Duration
	content   string
}

type frequency struct {
	tp freType
	d  uint
}

func NewEventServer(useCase *biz.JobUseCase) func(msg *models.RealRecvMsg, Bot *wcbot.WcBot) (isBreak bool) {
	return Server
}

func Server(msg *models.RealRecvMsg, Bot *wcbot.WcBot) (isBreak bool) {
	_, ok := parseEvent(msg.Content.Data)
	if ok {
		return false
	}

	return true
}

func parseEvent(content string) (*event, bool) {
	seg := strings.Split(content, splitChar)
	if len(seg) != 2 && len(seg) != 3 {
		return nil, false
	}
	switch len(seg) {
	case segTypeOnce:
		t, err := parseTime(seg[0])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		return &event{
			frequency: &frequency{
				tp: none,
			},
			when:    t,
			content: seg[1],
		}, true
	case segTypeLoop:
		f, err := parseFrequency(seg[0])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		t, err := parseTime(seg[1])
		if err != nil {
			logrus.Errorf("%s", err)
			return nil, false
		}

		return &event{
			frequency: f,
			when:      t,
			content:   seg[2],
		}, true
	}
	return nil, false
}

var (
	regFreType = regexp.MustCompile("(每天|每日|每周|每月|工作日)[^\\d]*(\\d+)")
)
var (
	errParseFrequency = errors.New("parse frequency error")
)

func parseFrequency(content string) (*frequency, error) {
	var tp freType
	var d uint

	res := regFreType.FindAllStringSubmatch(content, -1)
	if len(res) > 0 {
		day, err := strconv.ParseUint(res[0][2], 10, 64)
		switch res[0][1] {
		case "每天":
			fallthrough
		case "每日":
			tp = aDay
		case "每周":
			tp = aWeek
		case "每月":
			tp = aMonth
		case "工作日":
			tp = workDay
		default:
			tp = aMonth
		}
		// 如果是每周或者每月点类型的，则需要检查是否是合法的
		if tp&(aMonth|aWeek) > 0 && err != nil {
			return nil, fmt.Errorf("%w: %s", errParseFrequency, err)
		}
		d = uint(day)
		return &frequency{
			tp: tp,
			d:  d,
		}, nil
	}
	return nil, fmt.Errorf("%w: %s", errParseFrequency, "no match")
}

var errParseTime = errors.New("parse time error")

func parseTime(content string) (time.Duration, error) {
	t, err := time.ParseInLocation("15:04", content, time.Local)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errParseTime, err)
	}
	logrus.Debug(t.Format("2006-01-02 15:04:05"))
	zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	logrus.Debug(zero.Format("2006-01-02 15:04:05"))

	return t.Sub(zero), nil
}
