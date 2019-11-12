package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoDBHosts = "172.17.84.205:27017"
	AuthDatabase = "test"
	AuthUserName = "test"
	AuthPassword = "123456"
	MaxCon       = 300
)

type Person struct {
	Name      string
	Phone     string
	City      string
	Age       int8
	IsMan     bool
	Interests []string
}

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

	session.SetPoolLimit(MaxCon)
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	type User struct {
		Name string
		City string
	}
	var usrs []User
	querySpecify(session, "test", "people", nil, bson.M{"Name": 1, "City": 1}, &usrs)

	/*
		for iter.Next(&usr) {
			log.Println(usr)
		}
	*/

	/*
		for _, item := range usrs {
			log.Println(item)
		}
	*/

}

func querySpecify(session *mgo.Session, dbname string, tablename string, query interface{}, selectFields interface{}, result interface{}) {
	copySession := session.Clone()
	defer copySession.Close()

	collection := copySession.DB(dbname).C(tablename)

	//Using iterator prevent from taking up too much memory
	type User struct {
		Name string
		City string
	}
	var usrs []User
	collection.Find(bson.M{}).All(&usrs)

	log.Println(usrs)
}
