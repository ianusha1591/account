package main

import (
        "encoding/json"
        "fmt"
//        "github.com/jinzhu/now"
        "gopkg.in/mgo.v2/bson"
        //"time"
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
        c := session.DB("test5").C("people")
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

c := session.DB("test5").C("people")
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


//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
  //       "username": user.Username,
    //           "password": user.Password,

      //  })
claims := &jwt.StandardClaims{
              ExpiresAt: 1500,
              Issuer:    "test",
      }

token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &User{
		"username": user.Username,
		"password":  user.Password,
	},claims)

//Claims := make(jwt.MapClaims)
//Claims := make(jwt.MapClaims)
//Claims["iss"] = "testClaim"
//token.Claims["exp"] = 100
//time.Now().Add(time.Duration(10) * time.Second)
//fmt.Print(Claims["exp"])
//claims := &jwt.StandardClaims{
//		ExpiresAt: 1500,
//		Issuer:    "test",
//	}
//mySigningKey := []byte("secret")
//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	ss, err := token.SignedString(mySigningKey)
//	fmt.Printf("%v %v", ss, err)
fmt.Print(user.Username)
fmt.Print(user.Password)


        tokenString, error := token.SignedString([]byte("secret"))
if error != nil {
                fmt.Println(error)
        }
        json.NewEncoder(w).Encode(User{Token: tokenString})
fmt.Print(tokenString)
err = c.Update(bson.M{"username": user.Username,"password":user.Password}, bson.M{"$set": bson.M{"token": tokenString}})
//var tokenstring = tokenString
//token, err = jwt.Parse(tokenString) {
//		return []byte("AllYourBase"), nil
//	}

//	if token.Valid {
//		fmt.Println("You look nice today")
//	} else if ve, ok := err.(*jwt.ValidationError); ok {
	//	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	//		fmt.Println("That's not even a token")
	//	} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	//		// Token is either expired or not active yet
	//		fmt.Println("Timing is everything")
	//	} else {
	//		fmt.Println("Couldn't handle this token:", err)
	//	}
	//} else {
	//	fmt.Println("Couldn't handle this token:", err)
	//}


		
 
}
}

func main() {
        router := mux.NewRouter()
        fmt.Println("Starting the application...")

        router.HandleFunc("/signup", signup).Methods("POST")
        router.HandleFunc("/login", login).Methods("POST")
log.Fatal(http.ListenAndServe(":12345", router))
}

