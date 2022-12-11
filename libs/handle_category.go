package forum

import (
	"net/http"
	"strconv"
	"strings"
)

// Calling HandleCategory function whenever there is a request to the root URL
func HandleCategory(w http.ResponseWriter, r *http.Request) {
	pathString := strings.TrimPrefix(r.URL.Path, "/category/")
	params := strings.Split(pathString, "/")
	//	isCategoryExist := false
	category_id, err := strconv.Atoi(params[0])

	// Getting all the necessary information for the category.html page
	posts, comments_numbers, last_comment, creating_users, mess := GetCategoryPageInformation(category_id)
	if strings.Compare(mess, "200 OK") != 0 {
		RenderError(w, mess)
		return
	}

	// Getting the session token from the requested cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//Getting privilege from session
	x := sessions[cookie.Value]
	priv := x.Privilege

	switch r.Method {
	case "GET":
		if priv >= 0 {
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

			// Rendering the category.html page for the logged-in user
			if UserLoggedIn(mess) {
				RenderCategory(w, mess, mess, posts, comments_numbers, last_comment, creating_users)
			} else {
				// Rendering the category.html page for the unregistered user
				RenderCategory(w, "", error_message, posts, comments_numbers, last_comment, creating_users)
			}
		}
	default:
		// Handling errors in case another method than GET is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Getting all the necessary information for the category.html page
func GetCategoryPageInformation(category_id int) ([]Post, map[Post]int, map[Post]Comment, map[Comment]User, string) {

	// Getting all the posts from the database
	posts, mess := GetAllPosts(category_id)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the total number of comments per post
	comments_numbers, mess := GetCommentNumberPerPost(posts)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the last comment per post
	last_comments, mess := GetLastCommentPerPost(posts)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the creating user per comment
	creating_users, mess := GetCreatingUserPerComment(last_comments)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	return posts, comments_numbers, last_comments, creating_users, "200 OK"
}

// Rendering the index.html page with neccessary given information
func RenderCategory(w http.ResponseWriter, username, error_message string, posts []Post, comments_numbers map[Post]int, last_comment map[Post]Comment, creating_users map[Comment]User) {
	category_page := &CategoryPage{Username: username, ErrorMessage: error_message, Posts: posts, CommentNumbers: comments_numbers, LastComment: last_comment, CreatingUsers: creating_users}
	ExecutePage(w, "category.html", category_page)
}
