package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	MongoDBHosts = "172.17.84.205:27017"
	AuthDatabase = "test"
	AuthUserName = "test"
	AuthPassword = "123456"
	MaxPoolSize  = 300
)

func main() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession failed:%\n", err)
	}

	session.SetPoolLimit(MaxPoolSize)
	defer session.Close()
}
