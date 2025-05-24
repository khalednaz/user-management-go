package auth

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"go-user-app/database"
	"go-user-app/scheduler"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleapi "google.golang.org/api/oauth2/v2"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     "xxxxxxxx.apps.googleusercontent.com",
	ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

var oauthStateString = "randomstatestring" // Replace with secure random string in production

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Code exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	service, err := googleapi.New(client)
	if err != nil {
		http.Error(w, "Could not create userinfo client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, "Failed getting user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user already exists
	row := database.DB.QueryRow("SELECT id FROM users WHERE email=?", userInfo.Email)
	var userID int
	err = row.Scan(&userID)

	if err == sql.ErrNoRows {
		// New user: insert
		res, err := database.DB.Exec("INSERT INTO users(username, email, password) VALUES (?, ?, ?)", userInfo.Name, userInfo.Email, "google-oauth")
		if err != nil {
			http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		userID = int(id)
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(userID),
		Path:  "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "user_email",
		Value: userInfo.Email,
		Path:  "/",
	})

	// Start scheduler to send emails every minute
	scheduler.StartEmailScheduler(userInfo.Email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
