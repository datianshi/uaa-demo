package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var viewHtml = `
<html>
<body>
Token:  {{.Token}} </br>
Value: {{.Refresh}}
</body>
</html>
`

var testHtml = "HellowWorld"
var store = sessions.NewCookieStore([]byte("something-very-secret"))

type TokenValue struct {
	Token   string
	Refresh string
}

func init() {
	gob.Register(&TokenValue{})
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(testHtml))
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		fmt.Println(err)
	}
	token := session.Values["token"]
	if token == nil {
		http.Redirect(w, r, "http://www.google.com", http.StatusFound)
	}
	tplate, err := template.New("name").Parse(viewHtml)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tplate.Execute(w, token)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	token := &TokenValue{
		Token:   "abcde",
		Refresh: "ghijk",
	}
	session, err := getSession(r)
	if err != nil {
		fmt.Println(err)
	}
	session.Values["token"] = token
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("Can not persiste to the session: %v", err.Error())
	}
	http.Redirect(w, r, "/view", http.StatusFound)

}

func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "token-session")
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/token/", tokenHandler)
	http.HandleFunc("/test/", testHandler)
	addr := fmt.Sprintf(":%v", os.Getenv("PORT"))
	fmt.Println(addr)
	http.ListenAndServe(addr, nil)
}
