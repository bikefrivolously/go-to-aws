package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bikefrivolously/go-to-aws/internal/aws"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	profiles, err := aws.GetAwsProfiles()
	if err != nil {
		log.Println("Error getting list of profiles", err)
		return
	}
	t := template.Must(template.ParseFiles("templates/root.html"))
	t.Execute(w, profiles)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	profile_name := r.FormValue("profile")
	login_url, err := aws.GetLoginUrl(profile_name, "ca-central-1")
	if err != nil {
		fmt.Fprintln(w, "Error getting login URL:", err)
	}
	t := template.Must(template.ParseFiles("templates/login.html"))
	t.Execute(w, login_url)
}

func RunServer() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
