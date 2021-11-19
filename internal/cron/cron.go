package cron

import (
	"context"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

//
//
//
//        ***************************     ***************************         *********      ************************
//      *****************************    ******************************      *********      *************************
//     *****************************     *******************************     *********     *************************
//    *********                         *********                *******    *********     *********
//    ********                          *********               ********    *********     ********
//   ********     ******************   *********  *********************    *********     *********
//   ********     *****************    *********  ********************     *********     ********
//  ********      ****************    *********     ****************      *********     *********
//  ********                          *********      ********             *********     ********
// *********                         *********         ******            *********     *********
// ******************************    *********          *******          *********     *************************
//  ****************************    *********            *******        *********      *************************
//    **************************    *********              ******       *********         *********************
//
//

var Provider = wire.NewSet(NewCron)

type Cron struct {
	cron *cron.Cron
}

func NewCron() *Cron {
	c := cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(logrus.StandardLogger())),
	)
	ins := &Cron{
		cron: c,
	}
	go ins.Run()
	return ins
}

func (c *Cron) Run() {
	logrus.Info("---------------------CRON START---------------------")
	c.cron.Run()
}

func (c *Cron) AddCron(exp string, job func()) (cron.EntryID, error) {
	id, err := c.cron.AddFunc(exp, job)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Cron) DelCron(id cron.EntryID) {
	c.cron.Remove(id)
}

func (c *Cron) Stop() context.Context {
	return c.cron.Stop()
}

func (c *Cron) Start() {
	c.cron.Start()
}
