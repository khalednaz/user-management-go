package zoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"go-user-app/database"
)

const (
	accountID    = "FDaPujOkQjCq9one6fxLrw"
	clientID     = "BGZctwfNQ7GoHl7mgtEeww"
	clientSecret = "EGCS3RxwQlyVTuUA72RP4AIkjIn0wfCS"
)

func getServerToServerToken() (string, error) {
	url := "https://zoom.us/oauth/token?grant_type=account_credentials&account_id=" + accountID
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token error: %s", string(bodyBytes))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return "", err
	}

	token, ok := data["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("could not parse access_token")
	}

	return token, nil
}

func ZoomFormPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/zoom_form.html")
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func CreateZoomMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/zoom_form", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Error(w, "Unauthorized: no user email", http.StatusUnauthorized)
		return
	}
	userEmail := cookie.Value

	var userID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	token, err := getServerToServerToken()
	if err != nil {
		http.Error(w, "Token Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	topic := r.FormValue("topic")
	startTimeStr := r.FormValue("start_time")
	duration := r.FormValue("duration")

	startTime, err := time.Parse("2006-01-02T15:04", startTimeStr)
	if err != nil {
		http.Error(w, "Invalid start time: "+err.Error(), http.StatusBadRequest)
		return
	}
	startTimeUTC := startTime.UTC().Format("2006-01-02T15:04:05Z")

	meetingPayload := map[string]interface{}{
		"topic":      topic,
		"type":       2,
		"start_time": startTimeUTC,
		"duration":   duration,
		"timezone":   "UTC",
		"settings": map[string]interface{}{
			"join_before_host": true,
		},
	}

	payloadBytes, _ := json.Marshal(meetingPayload)

	req, err := http.NewRequest("POST", "https://api.zoom.us/v2/users/me/meetings", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Request Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Zoom API Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	startURL, startOK := result["start_url"].(string)
	joinURL, joinOK := result["join_url"].(string)

	if !startOK || !joinOK {
		http.Error(w, "Failed to create Zoom meeting: "+string(body), http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO meetings (user_id, topic, start_time, duration, join_url, start_url) VALUES (?, ?, ?, ?, ?, ?)",
		userID, topic, startTimeUTC, duration, joinURL, startURL)
	if err != nil {
		http.Error(w, "Failed to save meeting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard?success=meeting_created", http.StatusSeeOther)

}
