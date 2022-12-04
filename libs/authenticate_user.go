package forum

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

// Creating a global variable sessions to hold all the currently active sessions
var sessions = make(map[string]Session)

// Setting the client cookie with a generated session token
func SetClientCookieWithSessionToken(w http.ResponseWriter, username, password string) string {
	fmt.Println("Setting the client cookie with a generated session token")
	// Creating a random session token and an expired time (1 minute from the current time)
	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "500 INTERNAL SERVER ERROR: GENERATING SESSION TOKEN FAILED"
	}
	session_token := u2.String()
	expired_time := time.Now().Add(60 * time.Second)

	// Creating a new session for the given user with the above-generated session token and expired time
	sessions[session_token] = Session{Username: username, ExpiredTime: expired_time}

	// Setting the cookie of the current client with the above-generated session token and expired time
	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: session_token, Expires: expired_time})
	fmt.Println("debug: session_token =", session_token)

	return "200 OK"
}

// Checking if the given session has already been expired
func sessionExpired(session Session) bool {
	return session.ExpiredTime.Before(time.Now())
}

// Authenticating the user with the client cookie
func AuthenticateUser(w http.ResponseWriter, r *http.Request) string {

	// Getting the session token from the requested cookie
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		fmt.Fprintln(os.Stderr, err)
		return "401 UNAUTHORIZED: CLIENT COOKIE NOT SET"
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "400 BAD REQUEST: REQUEST NOT ALLOWED"
	}
	session_token := cookie.Value

	// Getting the corresponding session from the given session token
	session, status := sessions[session_token]
	if !status {
		return "401 UNAUTHORIZED: INVALID SESSION TOKEN"
	}

	// Checking if the session has already expired and removing it if that is the case
	if sessionExpired(session) {
		delete(sessions, session_token)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Expires: time.Now()})
		return "401 UNAUTHORIZED: SESSION EXPIRED"
	}

	return session.Username
}
