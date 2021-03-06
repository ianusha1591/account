package main

import (
        "encoding/json"
        "fmt"
 "time"
//"jwt.claims"
        "log"
        "net/http"
        "strings"
         "gopkg.in/mgo.v2"
        "github.com/dgrijalva/jwt-go"
        "github.com/gorilla/context"
        "github.com/gorilla/mux"
        "github.com/mitchellh/mapstructure"
)
type User struct {
        Username string `json:"username"`
        Password string `json:"password"`
// string `json:"expiration"`
 

        Token string `json:"token"`
}
type Exception struct {
        Message string `json:"message"`
}
func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
        var user User
       _ = json.NewDecoder(req.Body).Decode(&user)
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                "username": user.Username,
  //              "password": user.Password,
             //"expiration": 1521095590,
        })
claims := make(jwt.MapClaims)
claims["iss"] = "testClaim"
claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
fmt.Print(claims["exp"])
fmt.Print(user.Username)
fmt.Print(user.Password)


        tokenString, error := token.SignedString([]byte("secret"))

        if error != nil {
                fmt.Println(error)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)
session, err := mgo.Dial("127.0.0.1")
      if err != nil {
            panic(err)
  }

        defer session.Close()

         session.SetMode(mgo.Monotonic, true)



        // Collection People
        c := session.DB("test3").C("people")

         //Index
        index := mgo.Index{
               Key:        []string{"username", "password","token"},
                Unique:     true,
               DropDups:   true,
                Background: true,
 Sparse:     true,
       }

        err = c.EnsureIndex(index)
        if err != nil  && mgo.IsDup(err) {
panic(err)
        }

err = c.Insert(&User{Username:user.Username,Password:user.Password,Token:tokenString})

        if err != nil && mgo.IsDup(err) {

                panic(err)


}

}
func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
        params := req.URL.Query()
        token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("There was an error")
                }
                return []byte("secret"), nil
        })
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                var user User
                mapstructure.Decode(claims, &user)
                json.NewEncoder(w).Encode(user)
        } else {
                json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
        }
}
func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                authorizationHeader := req.Header.Get("authorization")
                if authorizationHeader != "" {
                        bearerToken := strings.Split(authorizationHeader, " ")
                        if len(bearerToken) == 2 {
                                token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                                        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                                                return nil, fmt.Errorf("There was an error")
                                        }
                                        return []byte("secret"), nil
                                })
                                if error != nil {
                                        json.NewEncoder(w).Encode(Exception{Message: error.Error()})
                                        return
                                }
                                if token.Valid {
                                        context.Set(req, "decoded", token.Claims)
                                        next(w, req)
                                } else {
                                        json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
                                }
                        }
                } else {
                        json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
                }
        })
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
        decoded := context.Get(req, "decoded")
        var user User
        mapstructure.Decode(decoded.(jwt.MapClaims), &user)
        json.NewEncoder(w).Encode(user)
}

func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
        router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
        router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")
        log.Fatal(http.ListenAndServe(":12345", router))
}



