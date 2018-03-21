package main

import (
        "encoding/json"
        "fmt"
//        "github.com/jinzhu/now"
        //"gopkg.in/mgo.v2/bson"
        //"time"
        "log"
        "net/http"
  //    "strings"
        //"gopkg.in/mgo.v2"
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
func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
        var user User
       _ = json.NewDecoder(req.Body).Decode(&user)

 token := jwt.New(jwt.SigningMethodHS256)
fmt.Print(user.Username)
fmt.Print(user.Password)

 tokenString, error := token.SignedString([]byte("secret"))

        if error != nil {
                fmt.Println(error)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)
}
func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
//        router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
  //      router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")
        log.Fatal(http.ListenAndServe(":12345", router))
}






