package handlers

import (
	"html/template"
	"log"
	"net/http"

	"go-user-app/database"
	"go-user-app/models"
)

type DashboardData struct {
	CurrentUser models.User
	Users       []models.User
	Meetings    []models.Meeting
}

func getCurrentUser(r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("user_email")
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = database.DB.QueryRow("SELECT id, username, email FROM users WHERE email = ?", cookie.Value).
		Scan(&user.ID, &user.Username, &user.Email)

	return user, err
}

func getAllMeetings() []models.Meeting {
	rows, err := database.DB.Query("SELECT topic, start_time, duration, join_url FROM meetings")
	if err != nil {
		log.Println("DB error:", err)
		return nil
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		err := rows.Scan(&m.Topic, &m.StartTime, &m.Duration, &m.JoinURL)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		meetings = append(meetings, m)
	}
	return meetings
}

func renderDashboard(w http.ResponseWriter, users []models.User, currentUser models.User, meetings []models.Meeting) {
	data := DashboardData{
		CurrentUser: currentUser,
		Users:       users,
		Meetings:    meetings,
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := database.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			http.Error(w, "Row scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	meetings := getAllMeetings()

	renderDashboard(w, users, currentUser, meetings)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := database.DB.Exec("INSERT INTO users(username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		http.Error(w, "Insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	_, err := database.DB.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	username := r.FormValue("username")
	email := r.FormValue("email")

	_, err := database.DB.Exec("UPDATE users SET username=?, email=? WHERE id=?", username, email, id)
	if err != nil {
		http.Error(w, "Update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	query := r.FormValue("query")
	rows, err := database.DB.Query("SELECT id, username, email FROM users WHERE username LIKE ?", "%"+query+"%")
	if err != nil {
		http.Error(w, "Search query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			http.Error(w, "Row scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	meetings := getAllMeetings()
	renderDashboard(w, users, currentUser, meetings)
}
