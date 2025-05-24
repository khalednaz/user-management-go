#  User Management Go

A full-stack meeting management system built with Go, featuring Google OAuth authentication, Zoom integration, scheduled email reminders, and a clean HTML-based dashboard.

---

## ✨ Features

### 🔐 Authentication

| Feature        | Description                                       |
|----------------|---------------------------------------------------|
| JWT Login      | Secure login system with token-based sessions     |
| Google OAuth   | One-click login via Google SSO                    |
| Route Security | Server-side session validation with middleware    |

### 📅 Zoom Integration

| Feature           | Description                                   |
|-------------------|-----------------------------------------------|
| Meeting Creation  | Create Zoom meetings directly from the UI     |
| Meeting Details   | View upcoming meeting information              |
| Join Sessions     | Users can join scheduled meetings easily       |

### 📧 Email Scheduler

| Feature               | Description                                          |
|-----------------------|------------------------------------------------------|
| Email Reminders       | Automatic email reminders for upcoming meetings      |
| Background Scheduler  | Periodic background job to check and send emails     |

---

## 🛠 Technology Stack

### Backend

- Go (Golang)
- SQLite (file-based DB)
- OAuth2 (Google)
- Zoom API
- Go Scheduler / Cron Jobs
- net/smtp (for sending emails)

### Frontend

- HTML / CSS
- Javascript
- Template rendering with Go's `html/template`

---

## ✅ Prerequisites

| Tool            | Version / Notes                  |
|-----------------|-----------------------------------|
| Go              | 1.21+                             |
| Zoom Account    | Developer account required        |
| Google Project  | OAuth2 Client setup required      |
| Gmail SMTP      | Used for sending email reminders  |

---

## 🚀 Setup Instructions

### 1. Clone Repository

```bash
git clone https://github.com/khalednaz/user-management-go.git
cd user-management-go


Google OAuth Setup
Go to Google Cloud Console

Create a new project

Enable Google+ API / OAuth 2.0

Create OAuth 2.0 Client ID

Add authorized redirect URI:


http://localhost:8080/auth/google/callback
5. Zoom API Setup
Go to Zoom Marketplace

Create an OAuth App

Set redirect URI:

http://localhost:8080/zoom/callback
Add these permissions:

Meeting:Write

Meeting:Read

Copy credentials into your .env

🗃️ Folder Structure
csharp
Copy
Edit
go-user-app/
│
├── auth/               # Google OAuth logic
├── database/           # DB connection (SQLite)
├── handlers/           # HTTP route handlers
├── models/             # Data models (User, Meeting)
├── scheduler/          # Email scheduler job
├── static/             # CSS, images
├── templates/          # HTML templates
├── utils/              # Helper methods (e.g. send mail)
├── zoom/               # Zoom API integration
└── main.go             # App entry point
