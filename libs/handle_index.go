package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// colors for debug prints
const green = "\033[32m"
const yellow = "\033[33m"
const red = "\033[31m"
const reset = "\033[0m"

// Initializing a global variable templ with all the necessary pages
var templ = template.Must(template.ParseFiles("templates/forums.html"))
var username = ""

// Calling HandleWelcome function whenever there is a request to the welcome URL
func HandleWelcome(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the welcome URL is being requested
	if strings.Compare(r.URL.Path, "/welcome") != 0 {
		renderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}

	// Opening the database and handling errors in case of failure
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		renderError(w, mess)
		return
	}
	defer db.Close()

	switch r.Method {
	case "GET":
		// Authenticating the user with the client cookie
		mess := AuthenticateUser(w, r)
		if strings.Compare(mess, username) != 0 {
			renderError(w, mess)
			return
		}

		// Getting all the necessary information for the index.html page
		categories, posts_numbers, last_posts, creating_users, mess := getIndexPageInformation(db)
		if strings.Compare(mess, "200 OK") != 0 {
			renderError(w, mess)
			return
		}

		// Rendering the index.html page for the logged-in user
		renderIndex(w, username, categories, posts_numbers, last_posts, creating_users)
	default:
		// Handling errors in case another method than GET or POST is being requested
		renderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Calling HandleIndex function whenever there is a request to the root URL
func HandleIndex(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the root URL is being requested
	if strings.Compare(r.URL.Path, "/") != 0 {
		renderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}

	// Opening the database and handling errors in case of failure
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		renderError(w, mess)
		return
	}
	defer db.Close()

	switch r.Method {
	case "POST":

		// Trying to parse the login form and handling errors in case of failure
		err := r.ParseForm()
		if err != nil {
			renderError(w, "500 INTERNAL SERVER ERROR: PARSING FORM FAILED")
			return
		}

		// Getting the credentials (given by the user) from the login form
		username = r.FormValue("uname")
		password := r.FormValue("pwd")

		// Getting the password of the given user from the database
		expected_pass, mess := GetPassword4User(db, username)
		if strings.Compare(mess, "200 OK") != 0 {
			renderError(w, mess)
			fmt.Println(red, "Server -> ERROR:", mess, reset)
			return
		}

		hashed_pass, err := HashPassword(password)
		if err != nil {
			renderError(w, "500 INTERNAL SERVER ERROR: HASHING PASSWORD FAILED")
			return
		}
		//fix hashed password to be in use here.. first add a register function
		fmt.Println(yellow, "Server -> Login atempt!", username, " ", hashed_pass, reset)
		if CheckPasswordHash(hashed_pass, expected_pass) {
			fmt.Println(green, "Server -> Password is correct", reset)
			fmt.Println(green, "Server -> User", username, "successfully logged in", reset)
		}

		// Checking if the given password is different than the expected one
		if strings.Compare(password, expected_pass) != 0 {
			renderError(w, "401 UNAUTHORIZED: PASSWORD MISMATCHED")
			fmt.Fprint(w, "401 UNAUTHORIZED: PASSWORD MISMATCHED")
			return
		}
		fmt.Println(green, "Server -> user "+username+" logged in", reset)
		// Setting the client cookie with a generated session token
		mess = SetClientCookieWithSessionToken(w, username, password)
		if strings.Compare(mess, "200 OK") != 0 {
			renderError(w, mess)
			return
		}

		// Redirecting the user to the welcome URL after the authentication
		http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)

	case "GET":
		// Getting all the necessary information for the index.html page
		categories, posts_numbers, last_posts, creating_users, mess := getIndexPageInformation(db)
		if strings.Compare(mess, "200 OK") != 0 {
			renderError(w, mess)
			return
		}

		// Rendering the index.html page for the unregistered user
		renderIndex(w, username, categories, posts_numbers, last_posts, creating_users)
	default:
		// Handling errors in case another method than GET or POST is being requested
		renderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Getting all the necessary information for the index.html page
func getIndexPageInformation(db *sql.DB) ([]Category, map[Category]int, map[Category]Post, map[Post]User, string) {

	// Getting all the categories from the database
	categories, mess := GetAllCategories(db)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the total number of posts per category
	posts_numbers, mess := GetPostsNumberPerCategory(db, categories)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the last post per category
	last_posts, mess := GetLastPostPerCategory(db, categories)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	// Getting the creating user per post
	creating_users, mess := GetCreatingUserPerPost(db, last_posts)
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, nil, nil, nil, mess
	}

	return categories, posts_numbers, last_posts, creating_users, "200 OK"
}

// Rendering the index.html page with neccessary given information
func renderIndex(w http.ResponseWriter, username string, categories []Category, posts_numbers map[Category]int, last_posts map[Category]Post, creating_users map[Post]User) {
	index_page := &IndexPage{Username: username, Categories: categories, PostsNumbers: posts_numbers, LastPosts: last_posts, CreatingUsers: creating_users}
	executeIndex(w, index_page)
}

// Rendering a given error with neccessary given information
func renderError(w http.ResponseWriter, message string) {
	index_page := &IndexPage{}
	switch message[:3] {
	case "404":
		w.WriteHeader(http.StatusNotFound)
	case "401":
		w.WriteHeader(http.StatusUnauthorized)
	case "400":
		w.WriteHeader(http.StatusBadRequest)
	case "500":
		w.WriteHeader(http.StatusInternalServerError)
	}
	executeIndex(w, index_page)
}

// Executing the index.html page and handling errors in case of failure
func executeIndex(w http.ResponseWriter, index_page *IndexPage) {
	err := templ.ExecuteTemplate(w, "forums.html", index_page)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "500 INTERNAL SERVER ERROR: PAGE CORRUPTED", http.StatusInternalServerError)
		return
	}
}
