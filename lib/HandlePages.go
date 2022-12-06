package lib

import (
	"forum/html"
	"net/http"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	p := html.DashboardParams{
		Title:   "Forum Re-make",
		Message: "Welcome to the forum",
		User:    "Mr Test",
		LogedIn: true,
	}
	html.Dashboard(w, p)
}

func ProfileShow(w http.ResponseWriter, r *http.Request) {
	p := html.ProfileShowParams{
		Title:   "Categories",
		Message: "Hello from Category show",
	}
	html.ProfileShow(w, p)
}

func ProfileEdit(w http.ResponseWriter, r *http.Request) {
	p := html.ProfileEditParams{
		Title:   "Category Edit",
		Message: "Hello from Category edit",
	}
	html.ProfileEdit(w, p)
}
