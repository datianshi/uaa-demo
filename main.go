package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const viewHtml = `
<html>
<body>
Token:  {{.Access}} </br>
Refresh: {{.Refresh}}
</body>
</html>
`

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type Config struct {
	Login        string
	Uaa          string
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

func init() {
	gob.Register(&uaa.Token{})
}

func viewHandler(w http.ResponseWriter, r *http.Request, uaaObject uaa.UAA) {
	session, err := getSession(r)
	pr(err)
	token := session.Values["token"]
	if token == nil {
		http.Redirect(w, r, uaaObject.LoginURL(), http.StatusFound)
	}
	tplate, err := template.New("name").Parse(viewHtml)
	pr(err)
	tplate.Execute(w, token)
}

func tokenHandler(w http.ResponseWriter, r *http.Request, uaaObject uaa.UAA) {
	code := r.URL.Query().Get("code")
	token, err := uaaObject.Exchange(code)
	pr(err)
	session, err := getSession(r)
	pr(err)
	session.Values["token"] = token
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("Can not persistent the session: %v", err.Error())
	}
	http.Redirect(w, r, "/view", http.StatusFound)

}

type UaaHandler func(http.ResponseWriter, *http.Request, uaa.UAA)

func makeHandler(fn UaaHandler, uaa uaa.UAA) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, uaa)
	}
}
func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "token-session")
}

func pr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func prAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	b, err := ioutil.ReadFile("config.json")
	prAndExit(err)
	var config Config
	err = json.Unmarshal(b, &config)
	prAndExit(err)
	var uaaObject = uaa.NewUAA(config.Login, config.Uaa, config.ClientId, config.ClientSecret, "")
	uaaObject.RedirectURL = config.RedirectUrl
	http.HandleFunc("/view/", makeHandler(viewHandler, uaaObject))
	http.HandleFunc("/token/", makeHandler(tokenHandler, uaaObject))
	addr := fmt.Sprintf(":%v", os.Getenv("PORT"))
	err = http.ListenAndServe(addr, nil)
	pr(err)
}
