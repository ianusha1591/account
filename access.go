package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	config = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       nil,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://ianusha.auth0.co/authorize",
			TokenURL: "https://ianusha.auth0.com/oauth2/token",
		},
		RedirectURL: "https://ec2-34-212-34-176.us-west-2.compute.amazonaws.com:8080/handle",
	}

	clientId     = flag.String("cid", "", "Client ID")
	clientSecret = flag.String("csec", "", "Client Secret")
)

func main() {
	flag.Parse()

	config.ClientID = *clientId
	config.ClientSecret = *clientSecret

	http.HandleFunc("/", landing)
	http.HandleFunc("/handle", handler)
	http.ListenAndServe(":8080", nil)
}

// A landing page redirects to the OAuth provider to get the auth code.
func landing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL(""), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("code"))
	token, err := config.Exchange(context.Background(), r.FormValue("code"))
	fmt.Println("token: ", token, "\n err: ", err)
}
