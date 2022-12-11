package forum

import (
	"net/http"
	"strconv"
	"strings"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	var postPage PostPage
	var err error

	pathString := strings.TrimPrefix(r.URL.Path, "/post/")
	params := strings.Split(pathString, "/")
	isPostExist := false

	if len(params) > 0 {
		id, errId := strconv.Atoi(params[0])
		if errId == nil {

			postPage.PostInfo, err = GetPost(id)
			if err == nil {
				isPostExist = true
			}

			//todo post clean
			if r.Method == "POST" {
				comment := r.FormValue("comment")
				userID := GetLoggedInUserID(w, r)

				//comment troll check
				if !TrollCheck(comment) {
					RenderError(w, "TROLL DETECTED write something normal")
					return
				}

				if userID != -1 {
					// todo check error + message to user if he is unauthorized
					CreateComment(userID, id, comment)
				}
			}

			if err == nil && isPostExist {
				//set comments data
				postPage.Comments, err = GetComments(id)
				if err != nil {
					RenderError(w, "500 INTERNAL SERVER ERROR: GETTING COMMENTS FAILED")
					return
				}

				// Authenticating the user with the client cookie
				mess := AuthenticateUser(w, r)
				if UserLoggedIn(mess) {
					userID := GetLoggedInUserID(w, r)
					postPage.Username = mess
					likeID, _ := IsPostLiked(userID, postPage.PostInfo.ID, 1)
					dislikeID, _ := IsPostLiked(userID, postPage.PostInfo.ID, 0)

					postPage.IsPostLiked = likeID > 0
					postPage.IsPostDisliked = dislikeID > 0

					postPage.CommentLike, postPage.CommentDislike, _ = GetUserLikedComments(postPage.Comments, userID)
				}

				ExecutePage(w, "post.html", postPage)
			}
		}
	}

	if !isPostExist {
		RenderError(w, "404 NOT FOUND: INVALID GIVEN URL")
		return
	}
}
