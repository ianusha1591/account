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
 token := jwt.New(jwt.SigningMethodHS256)
fmt.Print(user.Username)
fmt.Print(user.Password)
claims := make(jwt.MapClaims)
claims["exp"] =100
//    claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
    claims["iat"] = time.Now().Unix()
    //token.Claims = claims
fmt.Print(claims["exp"])
 tokenString, error := token.SignedString([]byte("secret"))

        if error != nil {
                fmt.Println(error)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)
err = c.Update(bson.M{"username": user.Username,"password":user.Password}, bson.M{"$set": bson.M{"token": tokenString}})
    now := claims["iat"].(int64)

 exp := claims["exp"].(int)
  if now > int64(exp) {
log.Print("your token is expired")
} else {
log.Print("your token is still valid, have a nice life")
}
}
}
func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/signup", signup).Methods("POST")
        router.HandleFunc("/login", login).Methods("POST")
log.Fatal(http.ListenAndServe(":12345", router))
}

