package forum

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("File")
	if err != nil {
		fmt.Println(Red, "Server -> Error Retrieving the File", Reset)
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a buffer to store the header of the file in
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(Yellow, "Server -> Checking for malicious content...", Reset)

	// Detect content type of file
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		{
			fmt.Println(Red, "Server -> malicious content detected! just kidding :) they did not uppload a image thats all...  File: ", handler.Filename, "(", handler.Size, "bytes)", Reset)
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}
	}

	//Server Debug
	fmt.Print(Green)
	fmt.Printf(" Server -> malicious content check passed! File: %+v (%v bytes)", handler.Filename, handler.Size)
	fmt.Println(Reset)

	// Create file in our fs directory
	dst, err := os.Create("fs/" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ExecutePage(w, "upload.html", nil)
	case "POST":
		uploadFile(w, r)
	}
}
