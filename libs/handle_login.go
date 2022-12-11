package forum

import (
	"net/http"
	"strings"
)

// Calling HandleLogin function whenever there is a request to the URL
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the login URL is being requested
	if strings.Compare(r.URL.Path, "/login") != 0 {
		RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}

	// Getting all the necessary information for the index.html page
	categories, posts_numbers, last_posts, creating_users, mess := GetIndexPageInformation()
	if strings.Compare(mess, "200 OK") != 0 {
		RenderError(w, mess)
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
		username := r.FormValue("uname")
		password := r.FormValue("pwd")
		error_message := ""

		// Getting the password of the given user from the database
		expected_pass, mess := GetPassword4User(username)
		if strings.Compare(mess, "200 OK") != 0 && strings.Compare(mess[:4], "401 ") != 0 {
			RenderError(w, mess)
			return
		}
		if strings.Compare(mess[:4], "401 ") == 0 {
			w.WriteHeader(http.StatusUnauthorized)
		}

		// Checking if the given password is different than the expected one
		error_message = mess
		if strings.Compare(error_message, "200 OK") == 0 && !CheckPasswordHash(expected_pass, password) {
			error_message = "401 UNAUTHORIZED: PASSWORD MISMATCHED"
			w.WriteHeader(http.StatusUnauthorized)
		}

		// Setting the client cookie with a generated session token
		mess = SetClientCookieWithSessionToken(w, username)
		if strings.Compare(mess, "200 OK") != 0 {
			RenderError(w, mess)
			return
		}

		// Redirecting the user to the welcome URL if the authentication succeeded
		if strings.Compare(error_message, "200 OK") == 0 {
			http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)
		} else {
			// Rendering the index.html page for the unregistered user with an unauthorized error message
			RenderIndex(w, "", error_message, categories, posts_numbers, last_posts, creating_users)
		}
	default:
		// Handling errors in case another method than POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}
