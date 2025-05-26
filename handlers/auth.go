package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"go-user-app/database"
	"go-user-app/scheduler"

	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{}
	if r.URL.Query().Get("error") == "invalid" {
		data["Error"] = "Invalid username or password."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{}
	errParam := r.URL.Query().Get("error")
	if errParam == "email_exists" {
		data["Error"] = "This email is already registered."
	} else if errParam == "insert_fail" {
		data["Error"] = "Failed to create account. Please try again."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// Check if email already exists
	var existingID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err != sql.ErrNoRows {
		http.Redirect(w, r, "/signup?error=email_exists", http.StatusSeeOther)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Redirect(w, r, "/signup?error=insert_fail", http.StatusSeeOther)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users(username, password, email) VALUES (?, ?, ?)", username, hashedPassword, email)
	if err != nil {
		http.Redirect(w, r, "/signup?error=insert_fail", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var id int
	var email, hashedPassword string

	row := database.DB.QueryRow("SELECT id, email, password FROM users WHERE username=?", username)
	err := row.Scan(&id, &email, &hashedPassword)
	if err != nil {
		// User not found
		tmpl, _ := template.ParseFiles("templates/login.html")
		data := map[string]interface{}{
			"Error": "User not found.",
		}
		tmpl.Execute(w, data)
		return
	}

	// Compare the entered password with the hashed one
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Wrong password
		tmpl, _ := template.ParseFiles("templates/login.html")
		data := map[string]interface{}{
			"Error": "Incorrect password.",
		}
		tmpl.Execute(w, data)
		return
	}

	// Successful login
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
