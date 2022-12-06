package lib

import (
	"fmt"
	"forum/html"
	"net/http"
)

// Anci Colors
var Green = "\033[32m"
var Red = "\033[31m"
var Yellow = "\033[33m"
var Reset = "\033[0m"

// Main dashboard page for the forum
func Dashboard(w http.ResponseWriter, r *http.Request) {
	CheckURL(w, r)
	p := html.DashboardParams{
		Title:   "Forum Re-make",
		Message: "Welcome to the forum",
		User:    "Mr Test",
		LogedIn: true,
	}

	err := html.Dashboard(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Functions that i use to test and leaarn templates
func ProfileShow(w http.ResponseWriter, r *http.Request) {
	p := html.ProfileShowParams{
		Title:   "Categories",
		Message: "Hello from Category show",
	}
	err := html.ProfileShow(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ProfileEdit(w http.ResponseWriter, r *http.Request) {
	p := html.ProfileEditParams{
		Title:   "Category Edit",
		Message: "Hello from Category edit",
	}

	err := html.ProfileEdit(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/profile" {
		fmt.Println(Red, "Server -> 404! Moving user back to index", Reset)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}
}
