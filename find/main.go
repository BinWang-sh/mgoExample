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

	result := Person{}
	err = QueryOne(session, "test", "people", &bson.M{"name": "Ale"}, &result)

	log.Println("Phone:", result.Phone)

	//Find with objectId
	id := "5dca1bb523c7a64c24e53ea7"
	objectId := bson.ObjectIdHex(id)
	err = QueryOne(session, "test", "people", &bson.M{"_id": objectId}, &result)
	log.Println(result)

	c := session.DB("test").C("people")

	//$ne city != Shanghai
	log.Println("\n---------city != Shanghai ------------------")
	iter := c.Find(bson.M{"city": bson.M{"$ne": "Shanghai"}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$gt age > 31
	log.Println("\n---------age > 31 ------------------")
	iter = c.Find(bson.M{"age": bson.M{"$gt": 31}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$lt age < 35
	log.Println("\n---------age < 35 ------------------")
	iter = c.Find(bson.M{"age": bson.M{"$lt": 35}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$gte age >= 33
	log.Println("\n---------age >= 33 ------------------")
	iter = c.Find(bson.M{"age": bson.M{"$gte": 33}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$lte age <= 33
	log.Println("\n---------age <= 33 ------------------")
	iter = c.Find(bson.M{"age": bson.M{"$lte": 33}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$in city in Shanghai Hangzhou
	log.Println("\n---------city in Shanghai Hangzhou ------------------")
	iter = c.Find(bson.M{"city": bson.M{"$in": []string{"Shanghai", "Hangzhou"}}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$nin city not in Shanghai Hangzhou
	log.Println("\n---------city not in Shanghai Hangzhou ------------------")
	iter = c.Find(bson.M{"city": bson.M{"$nin": []string{"Shanghai", "Hangzhou"}}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$exists city exist
	log.Println("\n---------city exist ------------------")
	iter = c.Find(bson.M{"city": bson.M{"$exists": true}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$exists $in city exists and value is null
	log.Println("\n---------city exist and value is null------------------")
	iter = c.Find(bson.M{"city": bson.M{"$in": []interface{}{nil}, "$exists": true}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$size interests size is 3
	log.Println("\n---------interests size is 3------------------")
	iter = c.Find(bson.M{"interests": bson.M{"$size": 3}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//$all interests include music and reading
	log.Println("\n---------interests include music and tea------------------")
	iter = c.Find(bson.M{"interests": bson.M{"$all": []string{"music", "reading"}}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//key.index first interest is music
	log.Println("\n---------first interest is music------------------")
	iter = c.Find(bson.M{"interests.0": "music"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//and
	log.Println("\n---------city == Shanghai and age >= 33 ------------------")
	iter = c.Find(bson.M{"city": "Shanghai", "age": bson.M{"$gte": 33}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//or
	log.Println("\n---------city == Shanghai and age >= 33 ------------------")
	iter = c.Find(bson.M{"$or": []bson.M{bson.M{"city": "Hangzhou"}, bson.M{"phone": "123432"}}}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	/*
		//Find all collections
		iter = QueryAll(session, "test", "people", nil)
		for iter.Next(&result) {
			log.Println(result)
		}
	*/

}

func QueryOne(session *mgo.Session, dbname string, tablename string, query interface{}, result interface{}) error {
	copySession := session.Clone()
	defer copySession.Close()

	collection := copySession.DB(dbname).C(tablename)
	err := collection.Find(query).One(result)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func QueryAll(session *mgo.Session, dbname string, tablename string, query interface{}) *mgo.Iter {
	copySession := session.Clone()
	defer copySession.Close()

	collection := copySession.DB(dbname).C(tablename)

	//Using iterator prevent from taking up too much memory
	iter := collection.Find(query).Iter()

	if iter != nil {
		return iter
	}

	return nil
}
