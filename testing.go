package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Email    string
	Password string
}

var db *sql.DB

func main() {
	// open database
	var err error
	db, err = sql.Open("sqlite3", "user.db")
	if err != nil {
		log.Fatal(err)
	}

	// create database table if no user exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			email VARCHAR NOT NULL,
			password VARCHAR NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	existingUser := struct{ Email string }{}
	err = db.QueryRow("SELECT email FROM user WHERE email=?", "lim113801@gmail.com").Scan(&existingUser.Email)
	if err == sql.ErrNoRows {
		result, err := db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", "lim113801@gmail.com", "zach")
		if err != nil {
			log.Fatal("Error !", err)
		}
		_ = result
	} else if err != nil {
		log.Fatal("Error ! No existing user found !:", err)
	}

	//login form
	tmpl, err := template.New("login").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Login</title>
			</head>
			<body>
				<h2>Login</h2>
				<form action="/login" method="post">
					<label>Email:</label>
					<input type="email" name="email" required><br>
					<label>Password:</label>
					<input type="password" name="password" required><br>
					<input type="submit" value="Login">
				</form>
			</body>
			</html>
		`)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		// login form handling

		loginDetails := User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		//compare
		var dbPass string
		tmplFinal := struct{ Success, Error bool }{false, false}

		err := db.QueryRow("SELECT password FROM user WHERE email=?", loginDetails.Email).Scan(&dbPass)
		if err != nil {
			tmplFinal.Error = true
		} else if loginDetails.Password == dbPass {
			tmplFinal.Success = true
		} else {
			tmplFinal.Error = true
		}

		tmpl.Execute(w, tmplFinal)

	})
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static", fs))

	log.Println(http.ListenAndServe(":8080", nil))

}
