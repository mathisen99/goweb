package forum

import (
	"net/http"
	"strconv"
	"strings"
)

// Calling HandleNewUser function whenever there is a request to the URL
func HandleNewUser(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the registration URL is being requested
	if strings.Compare(r.URL.Path, "/new_user") != 0 {
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
	priv := x.Privilege

	switch r.Method {
	case "POST":
		// Trying to parse the registration form and handling errors in case of failure
		err := r.ParseForm()
		if err != nil {
			RenderError(w, "500 INTERNAL SERVER ERROR: PARSING FORM FAILED")
			return
		}

		// Getting the credentials (given by the admin) from the registration form
		privilege, _ := strconv.Atoi(r.FormValue("privilege"))
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		repassword := r.FormValue("repassword")
		privilege_error := ""
		user_error := ""
		email_error := ""
		password_error := ""

		if GetUserID(username) == -1 && ValidatePasswordUsername(username, false) && ValidatePasswordUsername(password, true) &&
			len(username) > 3 && len(username) < 11 && GetEmail(email) != email && password == repassword && len(password) > 7 &&
			len(password) < 21 && privilege != 0 {
			passwordH, _ := HashPassword(password)
			if CheckPasswordHash(passwordH, repassword) {
				if CreateUser(username, passwordH, email, privilege) == "200 OK" {
					http.Redirect(w, r, "/", http.StatusMovedPermanently)
					return
				}
			}
		} else if GetUserID(username) != -1 {
			user_error = "User already exist"
		}
		if privilege == 0 {
			privilege_error = "Choose privilege"
		}
		if len(username) < 4 || len(username) > 10 || !ValidatePasswordUsername(username, false) {
			user_error = "Username invalid"
		}
		if !ValidateMail(email) {
			email_error = "Incorrectly given email"
		}
		if GetEmail(email) == email {
			email_error = "Email already used"
		}
		if password != repassword {
			password_error = "Passwords doesn't match"
		}
		if !ValidatePasswordUsername(password, true) || len(password) < 8 || len(password) > 20 {
			password_error = "Password invalid"
		}
		renderNewUser(w, privilege_error, user_error, email_error, password_error)
		return

	case "GET":
		if priv == 5 {
			ExecutePage(w, "new_user.html", &FormErrorPage{PrivilegeErrorMessage: "", UserErrorMessage: "", EmailErrorMessage: "", PasswordErrorMessage: "", CategoryErrorMessage: "", DescriptionErrorMessage: ""})
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

func renderNewUser(w http.ResponseWriter, privilege_error, user_error, email_error, password_error string) {
	new_user_page := &FormErrorPage{PrivilegeErrorMessage: privilege_error, UserErrorMessage: user_error, EmailErrorMessage: email_error, PasswordErrorMessage: password_error}
	ExecutePage(w, "new_user.html", new_user_page)
}
