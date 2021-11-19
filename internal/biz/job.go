package biz

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"wxBot4g/internal/cron"

	cron2 "github.com/robfig/cron/v3"
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

type FreType string

type JobFuncID uint

const (
	None    FreType = "* * *"
	ADay    FreType = "* * *"
	AWeek   FreType = "* * %d"
	AMonth  FreType = "%d * *"
	WorkDay FreType = "* * 1,2,3,4,5"
)

type JobFunc func(job *Job) error

var (
	errRegisterJobFunc = errors.New("register job func error")
	errGetJobFunc      = errors.New("get job func error")
)

// JobUseCase job 是一次性任务
type JobUseCase struct {
	repo       JobRepo
	jobFuncMap sync.Map
	cron       *cron.Cron
}

func (j *JobUseCase) PutJob(ctx context.Context, job *Job) (*Job, error) {
	job, err := j.repo.PutJob(ctx, job)
	if err != nil {
		return nil, err
	}
	entryID, err := j.cron.AddCron(job.CronExpress, j.withJobFunc(job.ID))
	if err != nil {
		return nil, err
	}

	job.CronID = entryID
	if err != nil {
		return nil, err
	}

	return j.repo.PutJob(ctx, job)
}

func (j *JobUseCase) DeleteJob(ctx context.Context, jobID uint64) error {
	return j.repo.DeleteJob(ctx, jobID)
}

func (j *JobUseCase) GetAllJobs(ctx context.Context, duration time.Duration) ([]*Job, error) {
	return j.repo.GetAllJobs(ctx, duration)
}

func (j *JobUseCase) RegisterJobFunc(typeID JobFuncID, f JobFunc) error {
	_, ok := j.jobFuncMap.Load(typeID)
	if ok {
		return fmt.Errorf("%w: job type %d already registered", errRegisterJobFunc, typeID)
	}
	j.jobFuncMap.Store(typeID, f)
	return nil
}

func (j *JobUseCase) GetJobFunc(typeID JobFuncID) (JobFunc, error) {
	jFunc, ok := j.jobFuncMap.Load(typeID)
	if !ok {
		return nil, fmt.Errorf("%w: job type %d not registered", errGetJobFunc, typeID)
	}
	return jFunc.(JobFunc), nil
}

func (j *JobUseCase) withJobFunc(jobID uint64) func() {
	return func() {
		job, err := j.repo.GetJobByID(context.TODO(), jobID)
		if err != nil {
			logrus.Errorf("get job by id error: %v", err)
			return
		}
		f, err := j.GetJobFunc(job.JobFuncID)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		err = f(job)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
	}
}

func NewJobUseCase(repo JobRepo, c *cron.Cron) *JobUseCase {
	return &JobUseCase{repo: repo, cron: c, jobFuncMap: sync.Map{}}
}

type JobRepo interface {
	PutJob(ctx context.Context, job *Job) (*Job, error)
	DeleteJob(ctx context.Context, jobID uint64) error
	GetAllJobs(ctx context.Context, duration time.Duration) ([]*Job, error)
	GetJobByID(ctx context.Context, jobID uint64) (*Job, error)
}

type Job struct {
	ID          uint64        `json:"id"`
	CronID      cron2.EntryID `json:"cron_id"`
	CronExpress string        `json:"cron_express"`
	Origin      string        `json:"origin"`
	Content     string        `json:"content"`
	JobType     uint          `json:"job_type"`
	JobFuncID   JobFuncID     `json:"job_func_id"`
}