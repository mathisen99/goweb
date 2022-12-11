package main

import (
	"fmt"
	. "forum/libs"
	"log"
	"net/http"
	"os"
)

func main() {
	// Creating a file server for css directory and taking it into use
	css_fs := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", css_fs))

	// Creating a file server for fs directory and taking it into use
	fs := http.FileServer(http.Dir("./fs"))
	http.Handle("/fs/", http.StripPrefix("/fs/", fs))

	// Creating a file server for js directory and taking it into use
	js_fs := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", js_fs))

	// Calling specific Handle functions whenever there are requests to corresponding URLs
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/welcome", HandleWelcome)
	http.HandleFunc("/logout", HandleLogout)
	http.HandleFunc("/registration", HandleRegistration)
	http.HandleFunc("/new_user", HandleNewUser)
	http.HandleFunc("/new_category", HandleNewCategory)
	http.HandleFunc("/post/", HandlePost)
	http.HandleFunc("/new_post", HandleNewPost)
	http.HandleFunc("/add_like", AddPostLike)
	http.HandleFunc("/add_comment_like", AddCommentLike)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/user/", HandleUser)
	http.HandleFunc("/change_theme", HandleThemes)

	// Starting listening to requests to port 8080 and logging errors in case of failure
	fmt.Fprintln(os.Stdout, "Server started successfully at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
