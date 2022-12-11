package forum

import (
	"net/http"
	"strconv"
)

func AddPostLike(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		postID, errPostID := strconv.Atoi(r.FormValue("postID"))
		userID := GetLoggedInUserID(w, r)
		isLike, errIsLike := strconv.Atoi(r.FormValue("isLike"))

		if userID != -1 && errPostID == nil && errIsLike == nil && (isLike == 1 || isLike == 0) {
			likeID, err := IsPostLiked(userID, postID, isLike)
			if err == nil {
				if likeID == 0 {
					err := AddLikeToPost(userID, postID, isLike)
					if err == nil {
						w.WriteHeader(http.StatusAccepted)
					}
				} else {
					err := RemoveLikeFromPost(likeID, userID, postID, isLike)
					if err == nil {
						w.WriteHeader(http.StatusAccepted)
					}
				}
			}
		}
	}
}

func AddCommentLike(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		commentID, errCommentID := strconv.Atoi(r.FormValue("commentID"))
		userID := GetLoggedInUserID(w, r)
		isLike, errIsLike := strconv.Atoi(r.FormValue("isLike"))

		if userID != -1 && errCommentID == nil && errIsLike == nil && (isLike == 1 || isLike == 0) {
			likeID, err := IsCommentLiked(userID, commentID, isLike)
			if err == nil {
				if likeID == 0 {
					err := AddLikeToComment(userID, commentID, isLike)
					if err == nil {
						w.WriteHeader(http.StatusAccepted)
					}
				} else {
					err := RemoveLikeFromComment(likeID, userID, commentID, isLike)
					if err == nil {
						w.WriteHeader(http.StatusAccepted)
					}
				}
			}
		}
	}
}
