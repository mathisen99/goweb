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

	// Creating a file server for js directory and taking it into use
	js_fs := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", js_fs))

	// Calling HandleIndex/HandleWelcome function whenever there is a request to the root/welcome URL
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/welcome", HandleWelcome)

	// Starting listening to requests to port 8080 and logging errors in case of failure
	fmt.Fprintln(os.Stdout, "Server started successfully at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
