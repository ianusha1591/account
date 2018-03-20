package main
import(
"fmt"
"time"
"github.com/dgrijalva/jwt-go"
)

func main(){
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
           "username":"anu",
               "password":"anusha",
        })
//Claims := make(jwt.MapClaims)
 tokenString, err := token.SignedString([]byte("secret"))
           if err !=nil{
fmt.Print(err)
}
fmt.Println(tokenString)




Claims := UserClaims{
    UserProfile{username: "user.Username", password: user.Password},
    jwt.StandardClaims{
        Issuer: "test-project",
ExpiresAt: 100
    },
}
token, err := jwtreq.ParseFromRequestWithClaims(req, jwtreq.AuthorizationHeaderExtractor, &claims, signingKeyFn)
    if err != nil {
        rw.WriteHeader(500)
        rw.Write([]byte("Failed to parse token"))
        log.Println("Failed to parse token")
        return
    }
if !token.Valid {
        log.Println("Invalid token")
     
    }

     } 
