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
	MaxPoolSize  = 300
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

	session.SetPoolLimit(MaxPoolSize)
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")

	result := Person{}

	//update $set
	log.Println("\n---------update set------------------")
	c.Update(bson.M{"name": "Tony"}, bson.M{"$set": bson.M{"name": "Tony Tang", "age": 30}})
	iter := c.Find(bson.M{"name": "Tony Tang"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//update $inc
	log.Println("\n---------update inc------------------")
	//c.Update(bson.M{"name": "Tony Tang"}, bson.M{"$inc": bson.M{"age": 1}})
	iter = c.Find(bson.M{"name": "Tony Tang"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//update $push
	log.Println("\n---------update push------------------")
	c.Update(bson.M{"name": "Tony Tang"}, bson.M{"$push": bson.M{"interests": "dream"}})
	iter = c.Find(bson.M{"name": "Tony Tang"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//update $pull
	log.Println("\n---------update pull------------------")
	c.Update(bson.M{"name": "Tony Tang"}, bson.M{"$pull": bson.M{"interests": "dream"}})
	iter = c.Find(bson.M{"name": "Tony Tang"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}

	//Remove
	log.Println("\n---------remove ------------------")
	c.Remove(bson.M{"name": "Smith"})
	iter = c.Find(bson.M{"name": "Smith"}).Iter()
	for iter.Next(&result) {
		log.Println(result)
	}
}

func createData(session *mgo.Session, dbname string, tablename string) error {

	persons := []Person{
		Person{Name: "Tony", Phone: "123432", City: "Shanghai", Age: 33, IsMan: true, Interests: []string{"music", "tea", "collection"}},
		Person{Name: "Mary", Phone: "232562", City: "Beijing", Age: 43, IsMan: false, Interests: []string{"sport", "film"}},
		Person{Name: "Tom", Phone: "123432", City: "Suzhou", Age: 22, IsMan: true, Interests: []string{"music", "reading"}},
		Person{Name: "Bob", Phone: "123432", City: "Hangzhou", Age: 32, IsMan: true, Interests: []string{"shopping", "coffee"}},
		Person{Name: "Alex", Phone: "15772", City: "Shanghai", Age: 21, IsMan: true, Interests: []string{"music", "chocolate"}},
		Person{Name: "Alice", Phone: "43456", City: "Shanghai", Age: 42, IsMan: false, Interests: []string{"outing", "tea"}},
		Person{Name: "Ingrid", Phone: "123432", City: "Shanghai", Age: 22, IsMan: false, Interests: []string{"travel", "tea"}},
		Person{Name: "Adle", Phone: "123432", City: "Shanghai", Age: 20, IsMan: false, Interests: []string{"game", "coffee", "sport"}},
		Person{Name: "Smith", Phone: "54223", City: "Fuzhou", Age: 54, IsMan: true, Interests: []string{"music", "reading"}},
		Person{Name: "Bruce", Phone: "123432", City: "Shanghai", Age: 31, IsMan: true, Interests: []string{"film", "tea", "game", "shoping", "reading"}},
	}

	cloneSession := session.Clone()
	c := cloneSession.DB(dbname).C(tablename)

	for _, item := range persons {
		err := c.Insert(&item)
		if err != nil {
			panic(err)
		}
	}

	return nil
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

func QueryAll2(session *mgo.Session, dbname string, tablename string, query interface{}, results interface{}) error {
	collection := session.DB(dbname).C(tablename)
	err := collection.Find(query).All(results)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
