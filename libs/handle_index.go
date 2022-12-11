package forum

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// Checking if the given message has a given error code
func HasErrorCode(message, error_code string) bool {
	return strings.HasPrefix(message, error_code)
}

// Calling HandleIndex function whenever there is a request to the root URL
func HandleIndex(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the root URL is being requested
	if strings.Compare(r.URL.Path, "/") != 0 {
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

		// Checking if the same user has just logged-in from a different browser
		error_message := "200 OK"
		if strings.Compare(mess, "401 UNAUTHORIZED: INVALID SESSION TOKEN") == 0 {
			error_message = mess
			w.WriteHeader(http.StatusUnauthorized)
		}

		// Rendering the index.html page for the logged-in user
		if UserLoggedIn(mess) {
			RenderIndex(w, mess, mess, categories, posts_numbers, last_posts, creating_users)
		} else {
			// Rendering the index.html page for the unregistered user
			RenderIndex(w, "", error_message, categories, posts_numbers, last_posts, creating_users)
		}
	default:
		// Handling errors in case another method than GET is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Getting all the necessary information for the index.html page
func GetIndexPageInformation() ([]Category, map[Category]int, map[Category]Post, map[Post]User, string) {

	// Getting all the categories from the database
	categories, mess := GetAllCategories()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the total number of posts per category
	posts_numbers, mess := GetPostsNumberPerCategory(categories)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the last post per category
	last_posts, mess := GetLastPostPerCategory(categories)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the creating user per post in the given map
	creating_users, mess := GetCreatingUserPerPostInMap(last_posts)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	return categories, posts_numbers, last_posts, creating_users, "200 OK"
}

// Rendering the index.html page with neccessary given information
func RenderIndex(w http.ResponseWriter, username, error_message string, categories []Category,
	posts_numbers map[Category]int, last_posts map[Category]Post, creating_users map[Post]User) {
	index_page := &IndexPage{Username: username, ErrorMessage: error_message, Categories: categories,
		PostsNumbers: posts_numbers, LastPosts: last_posts, CreatingUsers: creating_users}
	ExecutePage(w, "index.html", index_page)
}

// Rendering a given error with neccessary given information
func RenderError(w http.ResponseWriter, error_message string) {
	index_page := &IndexPage{Username: "", ErrorMessage: error_message, Categories: nil,
		PostsNumbers: nil, LastPosts: nil, CreatingUsers: nil}
	switch error_message[:4] {
	case "404 ":
		w.WriteHeader(http.StatusNotFound)
	case "401 ":
		w.WriteHeader(http.StatusUnauthorized)
	case "400 ":
		w.WriteHeader(http.StatusBadRequest)
	case "500 ":
		w.WriteHeader(http.StatusInternalServerError)
	}
	ExecutePage(w, "error.html", index_page)
}

// Executing the page given the page name and the page structure
func ExecutePage(w http.ResponseWriter, page_name string, page_structure interface{}) {

	// Initializing a variable templ with all the necessary pages
	templ := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"HasErrorCode": HasErrorCode, "GetUserPrivilege": GetUserPrivilege}).ParseFiles(
		"templates/index.html", "templates/registration.html", "templates/post.html",
		"templates/new_category.html", "templates/category.html", "templates/new_post.html",
		"templates/new_user.html", "templates/user.html", "templates/header.html",
		"templates/footer.html", "templates/error.html", "templates/login.html"))

	// Executing the given page and handling errors in case of failure
	err := templ.ExecuteTemplate(w, page_name, page_structure)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "500 INTERNAL SERVER ERROR: PAGE CORRUPTED", http.StatusInternalServerError)
		return
	}
}
