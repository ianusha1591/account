package main

import (
        "encoding/json"
        "fmt"
//        "github.com/jinzhu/now"
        "gopkg.in/mgo.v2/bson"
        "time"
        "log"
        "net/http"
  //    "strings"
        "gopkg.in/mgo.v2"
        "github.com/dgrijalva/jwt-go"
 //     "github.com/gorilla/context"
        "github.com/gorilla/mux"
 //     "github.com/mitchellh/mapstructure"
)
type User struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Token string `json:"token"`
}
type Exception struct {
        Message string `json:"message"`
}
func signup(w http.ResponseWriter, req *http.Request) {
var user User
 _ = json.NewDecoder(req.Body).Decode(&user)
 fmt.Print(user.Username)
fmt.Print(user.Password)
session, err := mgo.Dial("127.0.0.1")
      if err != nil {
            panic(err)
  }

        defer session.Close()

         session.SetMode(mgo.Monotonic, true)



        // Collection People
        c := session.DB("test4").C("people")
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

err = c.Insert(&User{Username:user.Username,Password:user.Password})

        if err != nil && mgo.IsDup(err) {
 panic(err)


}

}
func login(w http.ResponseWriter, req *http.Request) {
var user User
 _ = json.NewDecoder(req.Body).Decode(&user)
session, err := mgo.Dial("127.0.0.1")
      if err != nil {
            panic(err)
  }

        defer session.Close()

         session.SetMode(mgo.Monotonic, true)
        // Collection People
   
c := session.DB("test4").C("people")
fmt.Print("connection established")
err = c.Find(bson.M{"username": user.Username,"password":user.Password}).One(&user)
if err != nil {
//    if err.Error() == "not found" {
        log.Println("No such document")
     

  }  else
{
log.Print("ok")
//Claims := make(jwt.MapClaims)
//Claims["iss"] = "testClaim"
//Claims["exp"] = 100
//Claims["exp"] = time.Now().Add(time.Duration(10) * time.Second)


token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
           "username": user.Username,
               "password": user.Password,
             
        })
Claims := make(jwt.MapClaims)
Claims["iss"] = "testClaim"
Claims["exp"] = 100
time.Now().Add(time.Duration(10) * time.Second)


fmt.Print(Claims["exp"])

fmt.Print(user.Username)
fmt.Print(user.Password)


        tokenString, error := token.SignedString([]byte("secret"))

        if error != nil {
                fmt.Println(error)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)
err = c.Update(bson.M{"username": user.Username,"password":user.Password}, bson.M{"$set": bson.M{"token": tokenString}})
//now := TimeFunc().Unix()
//	if (now > Claims["exp"] ) { fmt.Print("expired") }
}
}
//}

func retrieve(w http.ResponseWriter, req *http.Request) {
var user User
 _ = json.NewDecoder(req.Body).Decode(&user)
session, err := mgo.Dial("127.0.0.1")
      if err != nil {
            panic(err)
  }

        defer session.Close()

         session.SetMode(mgo.Monotonic, true)
        // Collection People

c := session.DB("test4").C("people")
fmt.Print("connection established")
err = c.Find(bson.M{"username": user.Username,"password":user.Password}).One(&user)
if err != nil {
//    if err.Error() == "not found" {
        log.Println("No such document")


  }  else
{
log.Print("ok")
fmt.Print(user)
}
}

func retrieveAll(w http.ResponseWriter, req *http.Request) {
var user User
 _ = json.NewDecoder(req.Body).Decode(&user)
session, err := mgo.Dial("127.0.0.1")
      if err != nil {
            panic(err)
  }

        defer session.Close()

         session.SetMode(mgo.Monotonic, true)
        // Collection People

c := session.DB("test4").C("people")
fmt.Print("connection established")

  err = c.Find(nil).All(&user)
    if err != nil {
        // TODO: Do something about the error
    } else {
s := make([]byte,15,15) 
s=&user
        fmt.Println("Results All: ", s) 
    }


}


func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/signup", signup).Methods("POST")
        router.HandleFunc("/login", login).Methods("POST")
 router.HandleFunc("/retrieve", retrieve).Methods("POST")
router.HandleFunc("/retrieveAll", retrieveAll).Methods("POST")
log.Fatal(http.ListenAndServe(":12345", router))
}





