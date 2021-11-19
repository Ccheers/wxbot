package event

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ccheers/wxbot/internal/biz"
	"github.com/Ccheers/wxbot/internal/handler"
	"github.com/Ccheers/wxbot/models"
	"github.com/Ccheers/wxbot/wcbot"

	"github.com/sirupsen/logrus"
)

const splitChar = " "

const (
	segTypeOnce = 2 // 一次性提醒
	segTypeLoop = 3 // 循环提醒
)

var errJobExecFailed = errors.New("job exec failed")

const JobFuncEventFunc biz.JobFuncID = 1

func NewEventServer(useCase *biz.JobUseCase, bot *wcbot.WcBot) handler.MsgHandler {
	err := useCase.RegisterJobFunc(JobFuncEventFunc, callback(bot))
	if err != nil {
		panic(err)
	}

	return func(msg *models.RealRecvMsg) (isBreak bool) {
		job, ok := parseEvent(msg.Content.Data)
		if !ok {
			return false
		}
		job.FromUserID = msg.SendMsgUSer.Name
		job.JobFuncID = JobFuncEventFunc

		_, err := useCase.AddJob(context.TODO(), job)
		if err != nil {
			logrus.Error(err.Error())
			return false
		}

		logrus.Infof("event job: %+v", job)

		return true
	}
}

func parseEvent(content string) (*biz.Job, bool) {

	// 对内容进行修正，去除多余的空格
	content = fixContext(content)

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
			Origin:      content,
			Content:     seg[1],
			JobType:     biz.JobTypeOnce,
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
			Origin:      content,
			Content:     seg[2],
			JobType:     biz.JobTypeCron,
		}, true
	}
	return nil, false
}

var (
	regFreType = regexp.MustCompile("(每天|每日|每周|每月|工作日)(\\d*)")
)
var (
	errParseFrequency = errors.New("parse Frequency error")
)

var regexpSpace = regexp.MustCompile("\\s+")

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

func fixContext(content string) string {
	content = strings.TrimSpace(content)
	content = string(regexpSpace.ReplaceAll([]byte(content), []byte(splitChar)))
	return content
}

var errParseTime = errors.New("parse time error")

func parseTime(content string) (string, error) {
	t, err := time.ParseInLocation("15:04", content, time.Local)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errParseTime, err)
	}
	return t.Format("4 15"), nil
}

// 注册到 CRON 的 job，周期性任务到期会回调这个函数
func callback(bot *wcbot.WcBot) biz.CallbackFunc {
	return func(job *biz.Job) error {
		logrus.Infof("send msg to: %s content: %s", job.FromUserID, job.Content)
		ok := bot.SendMsg(job.FromUserID, fmt.Sprintf("[海军BOT]\n\r%s", job.Content), false)
		if !ok {
			return fmt.Errorf("%w: %s", errJobExecFailed, "send msg failed")
		}
		return nil
	}
}
