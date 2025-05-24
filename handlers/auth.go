package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"go-user-app/database"
	"go-user-app/scheduler"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	fmt.Println("Serving from directory:", dir)

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	_, err := database.DB.Exec("INSERT INTO users(username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	row := database.DB.QueryRow("SELECT id, email FROM users WHERE username=? AND password=?", username, password)

	var id int
	var email string
	err := row.Scan(&id, &email)
	if err != nil {
		http.Redirect(w, r, "/?error=1", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_email",
		Value: email,
		Path:  "/",
	})

	scheduler.StartEmailScheduler(email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "user_email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
