package data

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

const boltdbPath = "/var/data/boltdb"

type Data struct {
	db *bbolt.DB
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
	return &Data{
		db: db,
	}
}
