package forum

import (
	"net/http"
	"strings"
)

// Calling HandleLogout function whenever there is a request to the URL
func HandleLogout(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the logout URL is being requested
	if strings.Compare(r.URL.Path, "/logout") != 0 {
		RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}

	switch r.Method {
	case "POST":
		// Logging the currently logged-in user out
		mess := LogUserOut(w, r)
		if strings.Compare(mess, "200 OK") != 0 && strings.Compare(mess[:4], "401 ") != 0 {
			RenderError(w, mess)
			return
		}

		// Redirecting the user to the root URL if the logout succeeded
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	default:
		// Handling errors in case another method than POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}
