package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var tmpl *template.Template
var dbSessions = make(map[string]string)
var dbUsers = make(map[string]user)

type user struct {
	Name     string
	UserName string
	Password string
}

func init() {
	tmpl = template.Must(template.ParseGlob("Templates/*"))
	dbUsers["anan@gmail.com"] = user{"anandhu", "anan@gmail.com", "123"}
}

func main() {
	fmt.Printf("server running on port: 3000")
	http.HandleFunc("/", loginHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// loginHandler function

func loginHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")

	if err == nil {
		if _, ok := dbSessions[cookie.Value]; ok {
			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	}
	if req.Method == http.MethodPost {

		uname := req.FormValue("username")
		pass := req.FormValue("password")

		if _, ok := dbUsers[uname]; !ok {
			http.Error(w, "username does not match", http.StatusForbidden)
			return
		}
		if pass != dbUsers[uname].Password {
			http.Error(w, "password does not match", http.StatusForbidden)
			return
		}

		if pass == dbUsers[uname].Password {

			// create cookie

			uid := uuid.NewString()
			cookie = &http.Cookie{
				Name:  "session",
				Value: uid,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, req, "/home", http.StatusSeeOther)
			return
		}
	}

	tmpl.ExecuteTemplate(w, "login.html", nil)
}
