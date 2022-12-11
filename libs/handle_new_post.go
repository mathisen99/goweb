package forum

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var Green = "\033[32m"
var Yellow = "\033[33m"
var Red = "\033[31m"
var Reset = "\033[0m"

// Calling HandleNewCategory function whenever there is a request to the URL
func HandleNewPost(w http.ResponseWriter, r *http.Request) {

	// Checking if a different URL than the registration URL is being requested
	if strings.Compare(r.URL.Path, "/new_post") != 0 {
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

	switch r.Method {
	case "POST":
		//geting information for categories, posts_numbers, last_posts, creating_users
		categories, posts_numbers, last_posts, creating_users, mess := GetIndexPageInformation()
		if strings.Compare(mess, "200 OK") != 0 {
			RenderError(w, "")
			return
		}
		// Parse the form data
		err = r.ParseForm()
		if err != nil {
			RenderError(w, "500 INTERNAL SERVER ERROR: PARSING FORM FAILED")
			return
		}
		//values and check if they are valid length
		title := r.FormValue("title")
		if len(title) > 100 {
			fmt.Println(Red, "Server WARNING -> Title To long!", Reset)
			renderCreatePost(w, "", "Title is too long", categories, posts_numbers, last_posts, creating_users)
			return
		}
		content := r.FormValue("content")
		if !TrollCheck(content) {
			fmt.Println(Red, "Server WARNING -> Troll detected!", Reset)
			renderCreatePost(w, mess, "<h1 style=\"color: brown;\">"+"Troll Detection!!! please write something normal"+"</h1>", categories, posts_numbers, last_posts, creating_users)
			return
		}
		//check title and content for badwords

		//create new post
		CreatePost(user_id, title, content)

		//loop list of categories checked by user
		list := r.Form["categories"]
		for _, v := range list {
			//get last post id
			last_id := GetLastPostID()
			//get category id
			category_id := GetCategoryID(v)
			//set category relation
			SetCategoryRelation(last_id, category_id)
		}

		//Redirect to index after post has been created
		http.Redirect(w, r, "/", http.StatusFound)
		return

	case "GET":
		cookie_string, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			fmt.Println("error2")
			return
		}

		for _, p := range sessions {
			if p.Cookie == cookie_string {
				prev := GetUserPrivilege(p.Username)
				if prev >= 1 {
					fmt.Println(Green, "Server -> User with privilege 1 or higher is trying to access the new post page", Reset)
				} else {
					fmt.Println(Red, "Server WARNING -> User with not enough privilege is trying to access the new post page", Reset)
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
					return
				}
			}
		}
		categories, posts_numbers, last_posts, creating_users, mess := GetIndexPageInformation()
		if strings.Compare(mess, "200 OK") != 0 {
			RenderError(w, "")
			return
		}

		renderCreatePost(w, x.Username, mess, categories, posts_numbers, last_posts, creating_users)
	default:
		// Handling errors in case another method than GET or POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}

// Rendering the index.html page with neccessary given information
func renderCreatePost(w http.ResponseWriter, username, error_message string, categories []Category, posts_numbers map[Category]int, last_posts map[Category]Post, creating_users map[Post]User) {
	if error_message == "200 OK" {
		error_message = ""
	}
	new_post_page := &IndexPage{Username: username, ErrorMessage: error_message, Categories: categories, PostsNumbers: posts_numbers, LastPosts: last_posts, CreatingUsers: creating_users}
	ExecutePage(w, "new_post.html", new_post_page)
}

func TrollCheck(s string) bool {
	badwords := ReadBadWords()
	if len(s) > 3000 {
		fmt.Println(Red, "Server WARNING -> Content To long!", Reset)
		return false
	}
	for _, v := range strings.Split(s, " ") {
		if len(v) > 100 {
			fmt.Println(Red, "Server WARNING -> Word To long!", Reset)
			return false
		}
	}

	for _, v := range strings.Split(s, " ") {
		for _, k := range badwords {
			if strings.Compare(v, k) == 0 {
				fmt.Println(Red, "Server WARNING -> Badword detected!", Reset)
				return false
			}
		}
	}

	return true
}

// Reading the badwords.txt file and returning a list of badwords
func ReadBadWords() []string {
	file, err := os.Open("fs/badwords.txt")
	if err != nil {
		fmt.Println(Red, "Server WARNING -> Badwords file not found!", Reset)
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var badwords []string
	for scanner.Scan() {
		badwords = append(badwords, scanner.Text())
	}
	return badwords
}
