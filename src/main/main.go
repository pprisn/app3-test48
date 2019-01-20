//Test проверка подлинности пользователя Facebook
package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "net/http"
  "net/url"
  "strings"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/facebook"
)

var AccessToken string
var VerifyToken string
var Port string
var ClientID string


var (
  oauthConf = &oauth2.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET",
    RedirectURL:  "https://app3-test48.herokuapp.com/oauth2callback",
    Scopes:       []string{"public_profile"},
    Endpoint:     facebook.Endpoint,
  }
  oauthStateString = "thisshouldberandom"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">facebook</a>
</body></html>
`
//Проинициализируем ключи Keyfb Keyyandex
func init() {
	Port = os.Getenv("PORT")
	oauthConf.ClientID = os.Getenv("ClientID")
        oauthConf.ClientSecret = os.Getenv("ClientSecret")	
	//AccessToken = os.Getenv("ACCESS_TOKEN")
	//VerifyToken = os.Getenv("VERIFY_TOKEN")
}


func handleMain(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(htmlIndex))
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
  Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
  if err != nil {
    log.Fatal("Parse: ", err)
  }
  parameters := url.Values{}
  parameters.Add("client_id", oauthConf.ClientID)
  parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
  parameters.Add("redirect_uri", oauthConf.RedirectURL)
  parameters.Add("response_type", "code")
  parameters.Add("state", oauthStateString)
  Url.RawQuery = parameters.Encode()
  url := Url.String()
  http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
  state := r.FormValue("state")
  if state != oauthStateString {
    fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }

  code := r.FormValue("code")

  token, err := oauthConf.Exchange(oauth2.NoContext, code)
  if err != nil {
    fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }

  resp, err := http.Get("https://graph.facebook.com/me?access_token=" + url.QueryEscape(token.AccessToken))
  if err != nil {
    fmt.Printf("Get: %s\n", err)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }
  defer resp.Body.Close()

  response, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("ReadAll: %s\n", err)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }

  log.Printf("parseResponseBody: %s\n", string(response))

  http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
  http.HandleFunc("/", handleMain)
  http.HandleFunc("/login", handleFacebookLogin)
  http.HandleFunc("/oauth2callback", handleFacebookCallback)
  log.Printf("Started running on Port = %s\n",Port)
  if err := http.ListenAndServe(":"+Port, nil); err != nil {
	log.Fatal(err)
  }

}