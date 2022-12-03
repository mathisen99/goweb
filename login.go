package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type VerifyUserInput struct {
	UserEmail          string
	UserPassword       string
	UserSessionStorage string
}

type VerifyUserOutput struct {
	Result       string
	Content      string
	SessionLogin string
}

func checkLogin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println(params)
	// get the user input
	var input VerifyUserInput
	input.UserEmail = request.FormValue("email")
	input.UserPassword = request.FormValue("password")
	input.UserSessionStorage = request.FormValue("session_storage")

	// check the user input
	// if the user input is valid, create a session and return the session id
	// if the user input is invalid, return an error
	var output VerifyUserOutput
	output.Result = "success"
	output.Content = "You have successfully logged in"
	output.SessionLogin = "1234567890" // this is a fake session id

}
