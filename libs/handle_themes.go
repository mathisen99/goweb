package forum

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Calling HandleNewCategory function whenever there is a request to the URL
func HandleThemes(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		fmt.Println("POST")
		// Parse the form data
		theme := r.FormValue("theme")

		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		//runing bash script to change theme
		if theme == "dark" {
			cmd := exec.Command("sh", path+"/change-dark.sh", theme)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		} else {
			cmd := exec.Command("sh", path+"/change-light.sh", theme)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		}

		//Redirect to index after post has been created
		http.Redirect(w, r, "/", http.StatusFound)
		return

	case "GET":
		//Redirect to index
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		// Handling errors in case another method than GET or POST is being requested
		RenderError(w, "400 BAD REQUEST: REQUESTED METHOD NOT SUPPORTED")
	}
}
