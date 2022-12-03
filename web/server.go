package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// colors for debug prints
const green = "\033[32m"
const red = "\033[31m"
const reset = "\033[0m"

func Server() {
	//folder for css files
	img := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", img))

	//folder for css files
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	//folder for js files
	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))
	//start server
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	fmt.Println("Server starting on port 8080")
	http.HandleFunc("/", HandlePage)
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	//log if error
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}

}

func HandlePage(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/":
		//home page
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatalln("There's an error with the template:", err)
		}
		t.Execute(w, nil)

	case "/categories":
		t, err := template.ParseFiles("templates/categories.html")
		if err != nil {
			log.Fatalln("There's an error with the template:", err)
		}
		//switch statement to handle GET and POST
		switch r.Method {
		case "GET":
			fmt.Println("GET")
			t.Execute(w, nil)

		case "POST":
			t.Execute(w, nil)
			//for post requests when we doing functions or buttons.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			//get the value of the button clicked
			category := r.FormValue("category")
			//Create category
			Create_category(category)
			//list all categories
			all_categories := Get_categories()
			for _, category := range all_categories {
				fmt.Fprintf(w, "<br>"+"Category: %v, ID: %v, Created: %s", category.Id, category.Category, category.CreatedAt)
			}
			fmt.Fprint(w, "<br><br>==========================================")

			Get_last_post := Get_last_post()
			fmt.Fprintf(w, "<br> Last post: %v", Get_last_post)

		}
	case "/posts":
		t, err := template.ParseFiles("templates/posts.html")
		if err != nil {
			log.Fatalln("There's an error with the template:", err)
		}
		//switch statement to handle GET and POST
		switch r.Method {
		case "GET":
			fmt.Println("GET")
			t.Execute(w, nil)

		case "POST":
			t.Execute(w, nil)
			//for post requests when we doing functions or buttons.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			//get the values of the form
			title := r.FormValue("title")
			content := r.FormValue("content")
			category_id := r.FormValue("category_id")
			post_id := r.FormValue("post_id")

			fmt.Println("Title: ", title)
			fmt.Println("Content: ", content)
			fmt.Println("Category ID: ", category_id)
			fmt.Println("Post ID: ", post_id)

			//convert category_id and post_user_id to int
			category_id_int, _ := strconv.Atoi(category_id)
			post_user_id_int, _ := strconv.Atoi(post_id)

			//Create post
			Create_post(post_user_id_int, category_id_int, title, content, 0, 0)

			//Create post
			//Create_post(topic)
			//list all topics
			all_topics := Get_all_posts()
			for _, topic := range all_topics {
				fmt.Fprintf(w, "<br>"+"Post id: %v, Post user id: %v, post category id: %v, title %s, content %s, liked %v disliked %v created at %s", topic.Id, topic.category_id, topic.user_id, topic.title, topic.content, topic.liked_no, topic.disliked_no, topic.createdAt)
			}
		}

	case "/users":
		t, err := template.ParseFiles("templates/users.html")
		if err != nil {
			log.Fatalln("There's an error with the template:", err)
		}
		t.Execute(w, nil)
		switch r.Method {
		case "GET":
			fmt.Println("GET")

		case "POST":
			//for post requests when we doing functions or buttons.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			//Get the values from the form
			username := r.FormValue("username")
			password := r.FormValue("psw")
			//password_repeat := r.FormValue("psw-repeat")
			email := r.FormValue("email")
			//hash the password
			hashedPassword, _ := HashPassword(password)
			fmt.Println(red)
			fmt.Println("Server -> Hashing of passowords: THIS IS NOT WORKING CORRECT SEE /check_login notes! this is what current password hash would look like :", hashedPassword)
			fmt.Println(reset)
			//create user
			Create_user(username, password, email)
			//List all users
			users := Get_users()
			for _, user := range users {
				fmt.Fprint(w, "<br>", user)
			}
		}

	case "/check_login":
		switch r.Method {

		case "GET":
			fmt.Println("GET")

		case "POST":

			//it seems the hash checking is not working as it sohuld it creates a new user hash on same string everytime so we cant compare it to the one in the database
			//NEED TO FIX THIS working without the hash checking for now
			Signin(w, r)
			fmt.Println(green)
			var data VerifyUserInput
			json.NewDecoder(r.Body).Decode(&data)
			fmt.Println("Server -> Verifying user started for " + data.UserEmail + ":" + data.UserPassword + ":" + data.UserSessionStorage)
			//check if user exists
			if !Check_user_exists(data.UserEmail) {
				fmt.Println(red)
				fmt.Println("Server -> User doesn't exist")
				fmt.Println("Server -> Sending user back to login page")
				t, err := template.ParseFiles("templates/index.html")
				if err != nil {
					log.Fatalln("There's an error with the template:", err)
				}
				t.Execute(w, nil)
				fmt.Println(green)
				return
			} else {
				fmt.Println(green)
				fmt.Println("Server -> User exists")
				fmt.Println("Server -> Checking password")
				hashpass, _ := HashPassword(data.UserPassword)
				fmt.Println("Server -> Hashed password: " + hashpass)
				//check if password is correct
				if !Check_user_password(data.UserEmail, data.UserPassword) {
					fmt.Println(red)
					fmt.Println("Server -> Password is incorrect")
					fmt.Println("Server -> Sending user back to login page")
					t, err := template.ParseFiles("templates/index.html")
					if err != nil {
						log.Fatalln("There's an error with the template:", err)
					}
					t.Execute(w, nil)
					fmt.Println(reset)
					return
				} else {
					fmt.Println(green)
					fmt.Println("Server -> Everything is correct!")
					fmt.Println("Server -> Sending user to home page")
					//find out how session storage works and implement it and send users to home page

				}
			}

		default:
			t, _ := template.ParseFiles("templates/error404.html")
			t.Execute(w, nil)
		}
	}

}
