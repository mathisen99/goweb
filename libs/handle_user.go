package forum

import (
	"net/http"
	"strings"
)

// Calling HandleUser function whenever there is a request to the user URL
func HandleUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Refreshing the current session (in case of a logged-in user)
		RefreshSession(w, r)

		// Authenticating the user with the client cookie
		mess := AuthenticateUser(w, r)
		if !UserLoggedIn(mess) && strings.Compare(mess[:4], "401 ") != 0 {
			RenderError(w, mess)
			return
		}

		// Checking if the same user has just logged-in from a different browser
		error_message := "200 OK"
		if strings.Compare(mess, "401 UNAUTHORIZED: INVALID SESSION TOKEN") == 0 {
			error_message = mess
			w.WriteHeader(http.StatusUnauthorized)
		}

		if UserLoggedIn(mess) {
			// Getting the given user info from the URL
			user_info := strings.TrimPrefix(r.URL.Path, "/user/")

			// Checking if an invalid username is given in the user URL
			if len(user_info) < len(mess) ||
				strings.Compare(user_info[:len(mess)], mess) != 0 ||
				(len(user_info) > len(mess) && rune(user_info[len(mess)]) != '?') {
				RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
				return
			}

			// Getting all the necessary information for the user.html page
			created_posts, liked_posts, creating_users, message := getUserPageInformation(mess)
			if strings.Compare(message, "200 OK") != 0 {
				RenderError(w, message)
				return
			}

			// Rendering the user.html page for the logged-in user
			renderUser(w, mess, mess, created_posts, liked_posts, creating_users)

		} else {
			// Getting all the necessary information for the index.html page
			categories, posts_numbers, last_posts, creating_users, message := GetIndexPageInformation()
			if strings.Compare(message, "200 OK") != 0 {
				RenderError(w, message)
				return
			}

			// Rendering the index.html page for the unregistered user
			RenderIndex(w, "", error_message, categories, posts_numbers, last_posts, creating_users)
		}
	default:
		// Handling errors in case another method than GET is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Getting all the necessary information for the user.html page
func getUserPageInformation(username string) ([]Post, []Post, map[Post]User, string) {

	// Getting all the posts created by the given user
	created_posts, mess := GetCreatedPostsOfGivenUser(username)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, mess
	}

	// Getting all the posts liked by the given user
	liked_posts, mess := GetLikedPostsOfGivenUser(username)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, mess
	}

	// Getting the creating user per post in the given list
	creating_users, mess := GetCreatingUserPerPostInList(liked_posts)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, mess
	}

	return created_posts, liked_posts, creating_users, "200 OK"
}

// Rendering the user.html page for the logged-in user
func renderUser(w http.ResponseWriter, username, error_message string,
	created_posts, liked_posts []Post, creating_users map[Post]User) {
	user_page := &UserPage{Username: username, ErrorMessage: error_message,
		CreatedPosts: created_posts, LikedPosts: liked_posts, CreatingUsers: creating_users}
	ExecutePage(w, "user.html", user_page)
}
