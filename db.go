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
        Username      string
        Password     string
        Timestamp time.Time
}
var (
        IsDrop = true
)
var user string
var num string
func main() {
read := bufio.NewReader(os.Stdin)
   fmt.Println("enter password?")
    password , _ := read.ReadString('\n')
 fmt.Println("Your password is ", password)

reader := bufio.NewReader(os.Stdin)
   fmt.Println("enter username?")
    username , _ := reader.ReadString('\n')
 fmt.Println("Your name is ", username)

        session, err := mgo.Dial("127.0.0.1")
        if err != nil {
                panic(err)
        }

        defer session.Close()

        session.SetMode(mgo.Monotonic, true)



        // Collection People
        c := session.DB("test2").C("people")
// Index
        index := mgo.Index{
                Key:        []string{"username", "password"},
                Unique:     true,
                DropDups:   true,
                Background: true,
 Sparse:     true,
        }

        err = c.EnsureIndex(index)
        if err != nil  && mgo.IsDup(err) {
                panic(err)
        }

        // Insert Datas
        err = c.Insert(&Person{Username: username, Password: password, Timestamp: time.Now()})

        if err != nil && mgo.IsDup(err) {

                panic(err)


}
}


