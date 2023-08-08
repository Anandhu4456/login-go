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
type errors struct {
	UserNErr string
	PassErr  string
}

var errorval errors

func init() {
	tmpl = template.Must(template.ParseGlob("Templates/*"))
	dbUsers["anan@gmail.com"] = user{"anandhu", "anan@gmail.com", "123"}
}

func main() {
	fmt.Printf("server running on port: 3000")
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandler)
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
			errorval.UserNErr = "username error"
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		if pass != dbUsers[uname].Password {
			errorval.PassErr = "password error"
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		errorval.UserNErr = ""
		errorval.PassErr = ""

		if pass == dbUsers[uname].Password {

			// create cookie

			uid := uuid.NewString()
			cookie = &http.Cookie{
				Name:  "session",
				Value: uid,
			}
			http.SetCookie(w, cookie)
			dbSessions[cookie.Value] = uname
			http.Redirect(w, req, "/home", http.StatusSeeOther)
			return
		}
	}
	tmpl.ExecuteTemplate(w, "login.html", errorval)
}

// signupHandler function

func signupHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")

	if err == nil {
		if _, ok := dbSessions[cookie.Value]; ok {
			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	}
	if req.Method == http.MethodPost {

		name := req.FormValue("name")
		uname := req.FormValue("username")
		pass := req.FormValue("password")

		// check username already taken?
		if _, ok := dbUsers[uname]; ok {
			errorval.UserNErr = "username already taken"
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
		// store user in dbUsers
		dbUsers[uname] = user{name, uname, pass}
		uid := uuid.NewString()
		cookie = &http.Cookie{
			Name:  "session",
			Value: uid,
		}
		http.SetCookie(w, cookie)
		dbSessions[cookie.Value] = uname

		http.Redirect(w, req, "/home", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "signup.html", errorval)
}

// homeHandler function

func homeHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		errorval.UserNErr = ""
		errorval.PassErr = ""
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if _, ok := dbSessions[cookie.Value]; !ok {
		errorval.UserNErr = ""
		errorval.PassErr = ""
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var un string
	var usr user

	un = dbSessions[cookie.Value]
	usr = dbUsers[un]

	tmpl.ExecuteTemplate(w, "home.html", usr)
}

// logoutHandler function

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorval.UserNErr = ""
		errorval.PassErr = ""
		return
	}
	if _, ok := dbSessions[cookie.Value]; !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorval.UserNErr = ""
		errorval.PassErr = ""
		return
	}
	cookie.MaxAge = -1
	dbSessions[cookie.Value] = ""
	http.SetCookie(w, cookie)

	http.Redirect(w, req, "/", http.StatusSeeOther)
	errorval.UserNErr = ""
	errorval.PassErr = ""
	
}
