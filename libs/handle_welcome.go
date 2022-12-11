package forum

import (
	"net/http"
	"strings"
)

// Calling HandleWelcome function whenever there is a request to the URL
func HandleWelcome(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the welcome URL is being requested
	if strings.Compare(r.URL.Path, "/welcome") != 0 {
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
	case "GET":
		// Authenticating the user with the client cookie
		mess := AuthenticateUser(w, r)
		if !UserLoggedIn(mess) && strings.Compare(mess[:4], "401 ") != 0 {
			RenderError(w, mess)
			return
		}
		if strings.Compare(mess[:4], "401 ") == 0 {
			w.WriteHeader(http.StatusUnauthorized)
		}

		// Rendering the index.html page for the logged-in user
		if UserLoggedIn(mess) {
			RenderIndex(w, mess, mess, categories, posts_numbers, last_posts, creating_users)
		} else {
			// Rendering the index.html page for the unregistered user with an unauthorized error message
			RenderIndex(w, "", mess, categories, posts_numbers, last_posts, creating_users)
		}
	default:
		// Handling errors in case another method than GET is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}
