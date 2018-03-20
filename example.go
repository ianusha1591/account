package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
"bufio"
"os"
//"io"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

var (
	IsDrop = true
)
var user string
func main() {
reader := bufio.NewReader(os.Stdin)
  //  var user string
    fmt.Println("enter username?")
    user , _ := reader.ReadString('\n')

    fmt.Println("Your name is ", user)

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	if IsDrop {
		err = session.DB("test").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	// Collection People
	c := session.DB("test").C("people")

	// Index
	index := mgo.Index{
		Key:        []string{"name", "phone"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	// Insert Datas
	err = c.Insert(&Person{Name: user, Phone: "+55 53 1234 4321", Timestamp: time.Now()},
		&Person{Name: user, Phone: "+66 33 1234 5678", Timestamp: time.Now()})

	if err != nil && mgo.IsDup(err) {

		panic(err)
	
}

	// Query One
	result := Person{}
	err = c.Find(bson.M{"name": &user}).Select(bson.M{"phone": 0}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phone", result)

	// Query All
	var results []Person
	err = c.Find(bson.M{"name": &user}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	// Update
	colQuerier := bson.M{"name": &user}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
	err = c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	// Query All
	err = c.Find(bson.M{"name": &user}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

}
