package db

import (
	"core/pkg/config"
	"errors"
	"github.com/go-bongo/bongo"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var MONGO *MongoDB

type MongoDB struct {
	Connection *bongo.Connection
}

func (m *MongoDB) connect(config *bongo.Config) error {
	var err error
	m.Connection, err = bongo.Connect(config)
	if err != nil {
		return err
	}

	m.Connection.Session.SetSyncTimeout(time.Duration(10 * time.Second))
	m.Connection.Session.SetSocketTimeout(time.Duration(10 * time.Second))

	return nil
}

func (m *MongoDB) Get(collection string, caster interface{}, findBy interface{}, page, limit int) error {
	result := m.Connection.Collection(collection).Find(findBy)
	if result.Error != nil {
		return result.Error
	}

	if err := result.Query.Skip(page * limit).Limit(limit).All(caster); err != nil {
		return err
		m.Connection.Session.Refresh()
	}

	return nil
}

func (m *MongoDB) FindById(collection string, caster interface{}, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Not an objectID hex")
	}
	return m.Connection.Collection(collection).FindById(bson.ObjectIdHex(id), caster)
}

func (m *MongoDB) GetAll(collection string, caster interface{}, page, limit int) error {
	return m.Get(collection, caster, nil, page, limit)
}

func (m *MongoDB) Save(collection string, data bongo.Document) error {
	return m.Connection.Collection(collection).Save(data)
}

func NewMongo() error {

	MONGO = &MongoDB{}

	//mgo.SetDebug(true)
	//mgo.SetLogger(log.New(os.Stderr, "", log.LstdFlags))

	mongo_string := config.LoadedConfig.MongoConnection
	if mongo_string == "" {
		glog.Fatal("Missing MONGO_CONNECTION environment variable!")
	}

	mongo_config := &bongo.Config{mongo_string, config.LoadedConfig.MongoDatabase}

	if err := MONGO.connect(mongo_config); err != nil {
		glog.Fatalf("Error connecting to MongoDB")
		return err
	}

	glog.Info("Successfully connected to Mongo")

	return nil
}
