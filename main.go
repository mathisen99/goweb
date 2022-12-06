package main

import (
	"forum/lib"
	"net/http"
)

func main() {
	//fileserver for css
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	//pages
	http.HandleFunc("/dashboard", lib.Dashboard)
	http.HandleFunc("/profile/show", lib.ProfileShow)
	http.HandleFunc("/profile/edit", lib.ProfileEdit)
	//lISTEN
	http.ListenAndServe(":8080", nil)
}
