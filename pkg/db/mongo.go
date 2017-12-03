package db

import (
	"github.com/go-bongo/bongo"
	"os"
	"github.com/golang/glog"
)

var MONGO_CONNECTION *bongo.Connection = nil

func ConnectMongo() error {
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

	return nil
}