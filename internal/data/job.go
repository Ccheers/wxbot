package data

import (
	"context"
	"encoding/json"
	"time"
	"wxBot4g/internal/biz"
	"wxBot4g/internal/pkg/itob"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

type JobRepoImpl struct {
	data *Data
}

func NewJobRepoImpl(data *Data) biz.JobRepo {
	return &JobRepoImpl{data: data}
}

func (j *JobRepoImpl) PutJob(ctx context.Context, job *biz.Job) (*biz.Job, error) {
	err := j.data.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(j.TableName()))

		if job.ID == 0 {
			id, err := b.NextSequence()
			if err != nil {
				return err
			}
			job.ID = id
		}

		bts, err := json.Marshal(job)
		if err != nil {
			return err
		}
		err = b.Put(itob.Itob(job.ID), bts)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *JobRepoImpl) DeleteJob(ctx context.Context, jobID uint64) error {
	return j.data.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(j.TableName()))
		return b.Delete(itob.Itob(jobID))
	})
}

func (j *JobRepoImpl) GetAllJobs(ctx context.Context, duration time.Duration) ([]*biz.Job, error) {
	jobs := make([]*biz.Job, 0)
	err := j.data.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(j.TableName())).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			job := new(biz.Job)
			err := json.Unmarshal(v, job)
			if err != nil {
				logrus.Errorf(err.Error())
				continue
			}
			jobs = append(jobs, job)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (j *JobRepoImpl) GetJobByID(ctx context.Context, jobID uint64) (job *biz.Job, err error) {
	err = j.data.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(j.TableName()))

		v := b.Get(itob.Itob(jobID))
		err := json.Unmarshal(v, job)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (j *JobRepoImpl) TableName() string {
	return "job"
}
