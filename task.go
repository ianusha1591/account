package main
import (
//"errors"
"encoding/json"
        "fmt"
//"errors"
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
//type TokenClaims struct {
//StandardClaims:jwt.StandardClaims {
  //         ExpiresAt: expireToken,
    //                    Issuer:    "bandzest-auth",
    //            },

//}
type User struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Token string `json:"token"`
        jwt.StandardClaims

}
type Exception struct {
        Message string `json:"message"`
}
func signup(w http.ResponseWriter, req *http.Request) {
var user User
 _= json.NewDecoder(req.Body).Decode(&user)
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
        log.Println("No such document")

 }  else
{
log.Print("ok")
//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
   //        "username": user.Username,
  //             "password": user.Password,

    //    })
expireToken := time.Now().Add(time.Minute * 1).Unix()


claims := &jwt.StandardClaims{
                        ExpiresAt: expireToken,
                        Issuer:    "bandzest-auth",
}
// expireToken := time.Now().Add(time.Minute * 1).Unix()          
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

fmt.Print(claims)

fmt.Print(user.Username)
fmt.Print(user.Password)
 tokenString, errorr := token.SignedString([]byte("secret"))
  if errorr != nil {
fmt.Println(errorr)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)

err = c.Update(bson.M{"username": user.Username,"password":user.Password}, bson.M{"$set": bson.M{"token": tokenString}})
fmt.Print("your token is valid till a minute")

token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method")
		}
		return []byte("secret"), nil
	})

	log.Println(err)

	log.Println("Valid -- ", token.Valid)

 time.Sleep(time.Minute * 1)
fmt.Print("token is expired, you should login once again to get another token")


 message := make(chan string) // no buffer
     count := 3

     go func() {
         for i := 1; i <= count; i++ {
               fmt.Println("token will be deleted after the expiration time")
               message <- fmt.Sprintf("deleting  the token")
          }
     }()

   time.Sleep(time.Minute * 2)
fmt.Println(<-message)
err = c.Update(bson.M{"username": user.Username,"password":user.Password}, bson.M{"$set": bson.M{"token": " "}})

fmt.Println("deleted the token as it is expired")
}
}

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

func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/signup", signup).Methods("POST")
router.HandleFunc("/login", login).Methods("POST")
 router.HandleFunc("/retrieve", retrieve).Methods("POST")
log.Fatal(http.ListenAndServe(":12345", router))
}

