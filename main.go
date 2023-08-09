package main

import (
	"fmt"
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type users struct {
	Name     string
	Username string
	Password string
}

var dbUsers = make(map[string]users)
var dbSession = make(map[string]string)

type errorBase struct {
	EmailError    string
	NameError     string
	PasswordError string
}

var errorV errorBase

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("Templates/*"))
	dbUsers["anan@gmail.com"] = users{"Anandhu", "anan@gmail.com", "123"}
	dbUsers["vin@gmail.com"] = users{"vinay", "vin@gmail.com", "1234"}
	dbUsers["az@gmail.com"] = users{"azad", "az@gmail.com", "12345"}

}

func main() {
	fmt.Printf("server running on port:8080")
	http.HandleFunc("/", login)
	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)

}

// LoginHandler function

func login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	cookies, err := req.Cookie("session")
	if err == nil {
		if _, ok := dbSession[cookies.Value]; ok {
			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	}

	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")

		if _, ok := dbUsers[username]; !ok {
			errorV.EmailError = "email error"
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		if password != dbUsers[username].Password {
			errorV.PasswordError = "password error"
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		errorV.EmailError = ""
		errorV.PasswordError = ""
		if password == dbUsers[username].Password {
			// creating cookie

			uid := uuid.NewV4()
			cookie := &http.Cookie{
				Name:  "session",
				Value: uid.String(),
			}
			http.SetCookie(w, cookie)
			dbSession[cookie.Value] = username

			http.Redirect(w, req, "/home", http.StatusSeeOther)
			return
		}

	}
	tpl.ExecuteTemplate(w, "login.html", errorV)

}

// signupHandler function

func signup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	cookies, err := req.Cookie("session")
	if err == nil {
		if _, ok := dbSession[cookies.Value]; ok {
			http.Redirect(w, req, "/home", http.StatusSeeOther)
		}
	}

	if req.Method == "POST" {
		name := req.FormValue("name")
		username := req.FormValue("username")
		password := req.FormValue("password")

		// check username already taken ?

		if _, ok := dbUsers[username]; ok {
			errorV.EmailError = "email already taken"
			http.Redirect(w, req, "/signup", http.StatusSeeOther)
			return
		}

		errorV.EmailError = ""
		// storing new user to dbUsers

		dbUsers[username] = users{name, username, password}
		uid := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: uid.String(),
		}

		http.SetCookie(w, cookie)
		dbSession[cookie.Value] = username

		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "signup.html", errorV)
}

func home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	cookie, err := req.Cookie("session")

	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}

	if _, ok := dbSession[cookie.Value]; !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}
	var un string
	var usr users
	un = dbSession[cookie.Value]
	usr = dbUsers[un]

	tpl.ExecuteTemplate(w, "home.html", usr)
}

func logout(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	cookie, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}

	if _, ok := dbSession[cookie.Value]; !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		errorV.EmailError = ""
		errorV.PasswordError = ""
		return
	}
	// cookie deleting
	cookie.MaxAge = -1
	dbSession[cookie.Value] = ""
	http.SetCookie(w, cookie)

	http.Redirect(w, req, "/", http.StatusSeeOther)
	errorV.EmailError = ""
	errorV.PasswordError = ""

}
