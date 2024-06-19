package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bikefrivolously/go-to-aws/internal/aws"
)

func getProfileByName(profiles []aws.Profile, name string) *aws.Profile {
	for _, p := range profiles {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

type Server struct {
	Address  string
	Port     int
	Profiles []aws.Profile
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	s.Profiles, err = aws.GetAwsProfiles()
	if err != nil {
		log.Println("Error getting list of profiles", err)
		return
	}
	t := template.Must(template.ParseFiles("templates/root.html"))
	t.Execute(w, s.Profiles)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	profile_name := r.FormValue("profile")
	region := r.FormValue("region")
	p := getProfileByName(s.Profiles, profile_name)
	if region != "" {
		err = p.GetLoginUrl(region)
	} else {
		err = p.GetDefaultLoginUrl()
	}
	if err != nil {
		fmt.Fprintln(w, "Error getting login URL:", err)
		return
	}
	t := template.Must(template.ParseFiles("templates/login.html"))
	t.Execute(w, p.LoginUrl)
}

func (s *Server) Run() {
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/login", s.loginHandler)
	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
