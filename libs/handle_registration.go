package forum

import (
	"net/http"
	"strings"
)

// Calling HandleRegistration function whenever there is a request to the URL
func HandleRegistration(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the registration URL is being requested
	if strings.Compare(r.URL.Path, "/registration") != 0 {
		RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}

	switch r.Method {
	case "POST":
		// Trying to parse the login form and handling errors in case of failure
		err := r.ParseForm()
		if err != nil {
			RenderError(w, "500 INTERNAL SERVER ERROR: PARSING FORM FAILED")
			return
		}

		// Getting the credentials (given by the user) from the login form
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		repassword := r.FormValue("repassword")
		user_error := ""
		email_error := ""
		password_error := ""

		if GetUserID(username) == -1 && ValidatePasswordUsername(username, false) && ValidatePasswordUsername(password, true) &&
			len(username) > 3 && len(username) < 11 && GetEmail(email) != email && password == repassword && len(password) > 7 &&
			len(password) < 21 {
			passwordH, _ := HashPassword(password)
			if CheckPasswordHash(passwordH, repassword) {
				if CreateUser(username, passwordH, email, 1) == "200 OK" {
					http.Redirect(w, r, "/", http.StatusMovedPermanently)
					return
				}
			}
		} else if GetUserID(username) != -1 {
			user_error = "User already exist"
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
		renderRegistration(w, user_error, email_error, password_error)
		return

	case "GET":
		ExecutePage(w, "registration.html", &FormErrorPage{PrivilegeErrorMessage: "", UserErrorMessage: "", EmailErrorMessage: "", PasswordErrorMessage: "", CategoryErrorMessage: "", DescriptionErrorMessage: ""})
	default:
		// Handling errors in case another method than GET or POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

func renderRegistration(w http.ResponseWriter, user_error, email_error, password_error string) {
	registration_page := &FormErrorPage{UserErrorMessage: user_error, EmailErrorMessage: email_error, PasswordErrorMessage: password_error}
	ExecutePage(w, "registration.html", registration_page)
}

// email \w{1,64}[@]\w{1,64}[.]\w{1,64} regexp

// duck Donald4All
// Kitty H8everyone
// German IchHabeDurst1
// Police Hello999
// police Thieves911
// clown Circus77
// batman BatLover5
// robin Slave4Batman
// snake Sssssnake0
// pucko ZeroIQ00
