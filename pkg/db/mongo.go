package db

import (
	"github.com/go-bongo/bongo"
	"github.com/golang/glog"
	"os"
	"time"
)

var MONGO_CONNECTION *bongo.Connection = nil

func ConnectMongo() error {

	//mgo.SetDebug(true)
	//mgo.SetLogger(log.New(os.Stderr, "", log.LstdFlags))

	mongo_string := os.Getenv("MONGO_CONNECTION")
	if mongo_string == "" {
		glog.Fatal("Missing MONGO_CONNECTION environment variable!")
	}

	mongo_config := &bongo.Config{mongo_string, "caloriosa"}

	var err error

	MONGO_CONNECTION, err = bongo.Connect(mongo_config)
	if err != nil {
		glog.Fatal("Error connecting to Mongo! ", err)
	}

	glog.Info("Successfully connected to Mongo")

	MONGO_CONNECTION.Session.SetSyncTimeout(time.Duration(10 * time.Second))
	MONGO_CONNECTION.Session.SetSocketTimeout(time.Duration(10 * time.Second))

	return nil
}
