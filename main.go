package main

import (
	"log"
	"net/http"

	"go-user-app/auth"
	"go-user-app/database"
	"go-user-app/handlers"
	"go-user-app/zoom"
)

func main() {
	database.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.LoginPage)
	http.HandleFunc("/signup", handlers.SignUpPage)
	http.HandleFunc("/signup_submit", handlers.SignUpHandler)
	http.HandleFunc("/login_submit", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.Dashboard)
	http.HandleFunc("/create", handlers.CreateUser)
	http.HandleFunc("/delete", handlers.DeleteUser)
	http.HandleFunc("/update", handlers.UpdateUser)
	http.HandleFunc("/search", handlers.SearchUser)
	http.HandleFunc("/auth/google/login", auth.GoogleLogin)
	http.HandleFunc("/auth/google/callback", auth.GoogleCallback)
	http.HandleFunc("/zoom_form", zoom.ZoomFormPage)
	http.HandleFunc("/create_meeting", zoom.CreateZoomMeeting)
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
