package data

import (
	"fmt"
	"os"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"go.etcd.io/bbolt"
)

const boltdbPath = "/var/data/boltdb"

var Provider = wire.NewSet(NewData, NewJobRepoImpl)

type Data struct {
	db *bbolt.DB
}

type BucketName interface {
	TableName() string
}

func registerBucketNames() []BucketName {
	return []BucketName{
		new(JobRepoImpl),
	}
}

func NewData() *Data {
	err := os.MkdirAll(boltdbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	db, err := bbolt.Open(fmt.Sprintf("%s/%s", boltdbPath, "my.db"), 0600, nil)
	if err != nil {
		panic(err)
	}
	for _, name := range registerBucketNames() {
		err = db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name.TableName()))
			if err != nil {
				logrus.Errorf("create bucket error: %v", err)
			}
			return nil
		})
		if err != nil {
			logrus.Errorf("create bucket error: %v", err)
		}
	}
	return &Data{
		db: db,
	}
}
