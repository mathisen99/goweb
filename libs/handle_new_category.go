package forum

import (
	"net/http"
	"strings"
)

// Calling HandleNewCategory function whenever there is a request to the URL
func HandleNewCategory(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the registration URL is being requested
	if strings.Compare(r.URL.Path, "/new_category") != 0 {
		RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}
	// Getting the session token from the requested cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Getting username and check id from session
	x := sessions[cookie.Value]
	user_id := GetUserID(x.Username)
	priv := x.Privilege

	switch r.Method {
	case "POST":
		// Trying to parse the new category form and handling errors in case of failure
		err := r.ParseForm()
		if err != nil {
			RenderError(w, "500 INTERNAL SERVER ERROR: PARSING FORM FAILED")
			return
		}

		// Getting the credentials (given by the admin user) from the form
		category_name := r.FormValue("category-name")
		description := r.FormValue("category-description")
		category_error := ""
		description_error := ""

		if len(category_name) > 0 && len(category_name) <= 30 && len(description) <= 100 {
			if CreateCategory(category_name, description) == "200 OK" {
				if CreateDummyPost(category_name, user_id) == "200 OK" { //change user id
					if SetCategoryRelation(GetLastPostID(), GetCategoryID(category_name)) == "200 OK" {
						http.Redirect(w, r, "/", http.StatusMovedPermanently)
						return
					}
				}
			}
		} else if len(category_name) < 1 || len(category_name) > 30 {
			category_error = "Category length between 1-30"
		}
		if len(description) > 100 {
			description_error = "Description length maximum 100 characters"
		}
		renderCategory(w, category_error, description_error)
		return

	case "GET":
		if priv == 5 || priv == 3 {
			ExecutePage(w, "new_category.html", &FormErrorPage{PrivilegeErrorMessage: "", UserErrorMessage: "", EmailErrorMessage: "", PasswordErrorMessage: "", CategoryErrorMessage: "", DescriptionErrorMessage: ""})
			return
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

	default:
		// Handling errors in case another method than GET or POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

func renderCategory(w http.ResponseWriter, category_error, description_error string) {
	category_page := &FormErrorPage{CategoryErrorMessage: category_error, DescriptionErrorMessage: description_error}
	ExecutePage(w, "new_category.html", category_page)
}
