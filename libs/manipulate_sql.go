package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Opening the database
func OpenDatabase() (*sql.DB, string) {
	db, err := sql.Open("sqlite3", "./database/forum.db?parseTime=true")
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: DATABASE CORRUPTED"
	}

	return db, "200 OK"
}

// Getting all the posts created by the given user
func GetCreatedPostsOfGivenUser(username string) ([]Post, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the created posts
	var created_post Post
	var created_posts []Post

	// Selecting all posts (created by the given user) from the post table
	rows, err := db.Query("select id,user_id,title,content,created_at from post where user_id = ?", GetUserID(username))
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each row and saving the returned post
	for rows.Next() {
		err := rows.Scan(&created_post.ID, &created_post.UserID, &created_post.Title, &created_post.Content, &created_post.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Appending the returned post to the created posts list
		created_posts = append(created_posts, created_post)
	}

	return created_posts, "200 OK"
}

// Getting all the posts liked by the given user
func GetLikedPostsOfGivenUser(username string) ([]Post, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the liked posts
	var post_id int
	var liked_posts []Post

	// Selecting ids of all the posts (liked by the given user) from the user_post_reaction table
	rows, err := db.Query("select post_id from user_post_reaction where user_id = ? and is_liked = 1", GetUserID(username))
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: USER POST REACTION DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each post id and getting the corresponding post
	for rows.Next() {
		err := rows.Scan(&post_id)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER POST REACTION DATA CORRUPTED"
		}

		liked_post, err := GetPost(post_id)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Appending the returned post to the liked posts list
		liked_posts = append(liked_posts, liked_post)
	}

	return liked_posts, "200 OK"
}

// Getting all the categories from the database
func GetAllCategories() ([]Category, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the categories
	var category Category
	var categories []Category

	// Selecting all fields from the category table
	rows, err := db.Query("select id,category_name,descript,created_at from category")
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each row and saving the returned category
	for rows.Next() {
		err := rows.Scan(&category.ID, &category.CategoryName, &category.Description, &category.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
		}

		// Appending the returned category to the categories list
		categories = append(categories, category)
	}

	return categories, "200 OK"
}

// Getting the total number of posts per category
func GetPostsNumberPerCategory(categories []Category) (map[Category]int, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the posts numbers
	posts_numbers := make(map[Category]int)
	var posts_number int

	// Looping through each category and counting the total number of posts
	for i := 0; i < len(categories); i++ {

		// Getting the total number of posts for this category
		row := db.QueryRow("select count(*) from category_relation where category_id = ?", categories[i].ID)
		err := row.Scan(&posts_number)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}

		// Checking if there is at least one non-dummy post in the current category
		if posts_number > 1 &&
			strings.Compare(categories[i].CategoryName, "Cuisines") != 0 &&
			strings.Compare(categories[i].CategoryName, "Places") != 0 &&
			strings.Compare(categories[i].CategoryName, "Activities") != 0 {
			posts_number--
		}

		// Mapping this category to the total number of posts
		posts_numbers[categories[i]] = posts_number
	}

	return posts_numbers, "200 OK"
}

// Getting the last post per category
func GetLastPostPerCategory(categories []Category) (map[Category]Post, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the last posts
	var last_relation Relation
	var last_post Post
	last_posts := make(map[Category]Post)

	// Looping through each category and getting its last post
	for i := 0; i < len(categories); i++ {

		// Getting the last relation of this category
		row := db.QueryRow("select id,category_id,post_id from category_relation where category_id = ? order by post_id desc limit 1", categories[i].ID)
		err := row.Scan(&last_relation.ID, &last_relation.CategoryID, &last_relation.PostID)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}

		// Getting the last post given the post id from the last relation
		row = db.QueryRow("select id,user_id,title,content,created_at,liked_no,disliked_no from post where id = ?", last_relation.PostID)
		err = row.Scan(&last_post.ID, &last_post.UserID, &last_post.Title, &last_post.Content, &last_post.CreatedAt, &last_post.LikedNumber, &last_post.DislikedNumber)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		last_post.URL = "/post/" + strconv.Itoa(last_post.ID)
		// Mapping this category to the returned post
		last_posts[categories[i]] = last_post
	}

	return last_posts, "200 OK"
}

// Getting the creating user per post in a list
func GetCreatingUserPerPostInList(posts []Post) (map[Post]User, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the creating users
	var creating_user User
	creating_users := make(map[Post]User)

	// Looping through each post and getting its creating user
	for i := 0; i < len(posts); i++ {

		// Reading the only row and saving the returned user
		row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where id = ?", posts[i].UserID)
		err := row.Scan(&creating_user.ID, &creating_user.Privilege, &creating_user.Username, &creating_user.Password, &creating_user.Email, &creating_user.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
		}

		// Mapping this post to the returned user
		creating_users[posts[i]] = creating_user
	}

	return creating_users, "200 OK"
}

// Getting the creating user per post in a map
func GetCreatingUserPerPostInMap(posts map[Category]Post) (map[Post]User, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the creating users
	var creating_user User
	creating_users := make(map[Post]User)

	// Looping through each post and getting its creating user
	for _, post := range posts {

		// Reading the only row and saving the returned user
		row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where id = ?", post.UserID)
		err := row.Scan(&creating_user.ID, &creating_user.Privilege, &creating_user.Username, &creating_user.Password, &creating_user.Email, &creating_user.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
		}

		// Mapping this post to the returned user
		creating_users[post] = creating_user
	}

	return creating_users, "200 OK"
}

// Getting the password of the given user from the database
func GetPassword4User(username string) (string, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return "", mess
	}
	defer db.Close()

	// Creating a variable to hold the expected user
	var expected_user User

	// Reading the only row and saving the returned user
	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return "", "401 UNAUTHORIZED: USER NOT FOUND"
	}

	return expected_user.Password, "200 OK"
}

// Function to check if user exist in database
func GetUserID(username string) int {
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return -1
	}
	defer db.Close()

	// Creating a variable to hold the expected user
	var expected_user User

	// Reading the only row and saving the returned user
	row := db.QueryRow("select id,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return -1
	}

	return expected_user.ID
}

// Function to check if user exist in database
func GetUserName(id int) (string, error) {

	var err error
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return "", err
	}
	defer db.Close()

	// Creating a variable to hold the expected user
	var expected_user User

	// Reading the only row and saving the returned user
	row := db.QueryRow("select id,username,passwrd,email,created_at from user where id = ?", id)
	err = row.Scan(&expected_user.ID, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return "", err
	}

	return expected_user.Username, err
}

// function to create user
func CreateUser(username, password, email string, prev int) string {
	//open database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()
	//insert into user table
	_, err := db.Exec("insert into user(privilege,username,passwrd,email,created_at) values(?,?,?,?,?)", prev, username, password, email, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	return "200 OK"
}

func GetPost(id int) (Post, error) {

	var post Post
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		err = errors.New(mess)
		return post, err
	}
	defer db.Close()

	query := fmt.Sprintf("select id,user_id,title,content,created_at,liked_no, disliked_no from post where id = %v", id)
	rows, err := db.Query(query)
	if err != nil {
		return post, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.LikedNumber, &post.DislikedNumber)

		if err != nil {
			return post, err
		}

		post.Date = post.CreatedAt.Format("02.01.2006 03:04")
		post.UserName, err = GetUserName(post.UserID)

		if err != nil {
			return post, err
		}
	} else {
		err = errors.New("No post found")
	}

	return post, err
}

func GetComments(id int) ([]Comment, error) {

	var comment Comment
	var comments []Comment
	var users []int
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		err = errors.New(mess)
		return comments, err
	}
	defer db.Close()

	query := fmt.Sprintf("select id, user_id, content, created_at, liked_no, disliked_no from comment where post_id = %v", id)
	rows, err := db.Query(query)
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.LikedNumber, &comment.DislikedNumber)
		if err != nil {
			return comments, err
		}
		users = append(users, comment.UserID)
		comment.Date = comment.CreatedAt.Format("02.01.2006 03:04")
		//todo get from 1 request
		comment.UserName, err = GetUserName(comment.UserID)
		comments = append(comments, comment)
	}

	return comments, err
}

func AddLikeToPost(userID int, postID int, isLike int) error {
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return errors.New(mess)
	}
	defer db.Close()

	_, err = db.Exec("insert into user_post_reaction(user_id, post_id, is_liked, created_at) values(?,?,?,?)", userID, postID, isLike, time.Now())
	if err != nil {
		return err
	}

	if isLike == 1 {
		_, err = db.Exec("update post set liked_no = liked_no + 1 where id = ?", postID)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("update post set disliked_no = disliked_no + 1  where id = ?", postID)
		if err != nil {
			return err
		}
	}

	return err
}

func RemoveLikeFromPost(likeID int, userID int, postID int, isLike int) error {
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return errors.New(mess)
	}
	defer db.Close()

	_, err = db.Exec("delete from user_post_reaction where id = ?", likeID)
	if err != nil {
		return err
	}

	if isLike == 1 {
		_, err = db.Exec("update post set liked_no = liked_no - 1 where id = ?", postID)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("update post set disliked_no = disliked_no - 1  where id = ?", postID)
		if err != nil {
			return err
		}
	}

	return err
}

func IsPostLiked(userID int, postID int, isLike int) (int, error) {

	var err error
	var likeID int

	isLiked := 0

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return isLiked, errors.New(mess)
	}
	defer db.Close()

	row := db.QueryRow("select id from user_post_reaction where user_id = ? and post_id = ? and is_liked = ?", userID, postID, isLike)

	rowErr := row.Scan(&likeID)
	if rowErr == sql.ErrNoRows {
		return isLiked, err
	} else {
		if rowErr != nil {
			return isLiked, rowErr
		}
		return likeID, err
	}
}

func AddLikeToComment(userID int, commentID int, isLike int) error {
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return errors.New(mess)
	}
	defer db.Close()

	_, err = db.Exec("insert into user_comment_reaction(user_id, comment_id, is_liked, created_at) values(?,?,?,?)", userID, commentID, isLike, time.Now())
	if err != nil {
		return err
	}

	if isLike == 1 {
		_, err = db.Exec("update comment set liked_no = liked_no + 1 where id = ?", commentID)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("update comment set disliked_no = disliked_no + 1  where id = ?", commentID)
		if err != nil {
			return err
		}
	}

	return err
}

func RemoveLikeFromComment(likeID int, userID int, commentID int, isLike int) error {
	var err error

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return errors.New(mess)
	}
	defer db.Close()

	_, err = db.Exec("delete from user_comment_reaction where id = ?", likeID)
	if err != nil {
		return err
	}

	if isLike == 1 {
		_, err = db.Exec("update comment set liked_no = liked_no - 1 where id = ?", commentID)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("update comment set disliked_no = disliked_no - 1  where id = ?", commentID)
		if err != nil {
			return err
		}
	}

	return err
}

func IsCommentLiked(userID int, commentID int, isLike int) (int, error) {

	var err error
	var likeID int

	isLiked := 0

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return isLiked, errors.New(mess)
	}
	defer db.Close()

	row := db.QueryRow("select id from user_comment_reaction where user_id = ? and comment_id = ? and is_liked = ?", userID, commentID, isLike)

	rowErr := row.Scan(&likeID)
	if rowErr == sql.ErrNoRows {
		return isLiked, err
	} else {
		if rowErr != nil {
			return isLiked, rowErr
		}
		return likeID, err
	}
}

func GetUserLikedComments(comments []Comment, userID int) (map[int]bool, map[int]bool, error) {

	var err error
	var commentsId []string
	var isLiked, commentId int

	commentsLikes := make(map[int]bool)
	commentsDislikes := make(map[int]bool)

	if len(comments) > 0 {
		for _, comment := range comments {
			commentsId = append(commentsId, strconv.Itoa(comment.ID))
		}
		// Opening the database
		db, mess := OpenDatabase()
		if strings.Compare(mess, "200 OK") != 0 {
			return commentsLikes, commentsDislikes, errors.New(mess)
		}
		defer db.Close()

		rows, err := db.Query("select comment_id, is_liked from user_comment_reaction where user_id = ? and comment_id in (2,3)", userID, strings.Join(commentsId, ","))

		if err != nil {
			return commentsLikes, commentsDislikes, err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&commentId, &isLiked)

			if err == nil {
				if isLiked == 1 {
					commentsLikes[commentId] = true
				} else {
					commentsDislikes[commentId] = true
				}
			}
		}
	}

	return commentsLikes, commentsDislikes, err
}

// function to create a comment
func CreateComment(userID int, postID int, comment string) string {
	liked := 0
	disliked := 0
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()

	//insert into comment table
	_, err := db.Exec("insert into comment(user_id,post_id,content,liked_no,disliked_no,created_at) values(?,?,?,?,?,?)", userID, postID, comment, liked, disliked, time.Now())
	if err != nil {
		//todo change
		log.Fatal(err)
	}

	return "200 OK"
}

// Fetching email if it exist, else return empty string
func GetEmail(mail string) string {
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return ""
	}
	defer db.Close()
	var expected_user User

	// Reading the only row and saving the returned user
	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where email = ?", mail)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return ""
	}
	return expected_user.Email
}

// Checking if given emailaddress is valid
func ValidateMail(mail string) bool {
	if strings.Contains(mail[1:], "@") && len(mail) > 5 {
		for i, symb := range mail {
			if symb == '@' {
				if strings.Contains(mail[i+2:], ".") && mail[len(mail)-1] != '.' {
					return true
				}
			}
		}
	}
	return false
}

// Checking if given input only include characters a-z, A-Z and 0-9. For password all differentneed to be used.
func ValidatePasswordUsername(pwd string, stat bool) bool {
	A := false
	a := false
	numb := false
	validate := true
	for _, symb := range pwd {
		if symb >= 'A' && symb <= 'Z' {
			A = true
		} else if symb >= 'a' && symb <= 'z' {
			a = true
		} else if symb >= '0' && symb <= '9' {
			numb = true
		} else {
			validate = false
		}
	}
	if stat {
		if !A || !a || !numb {
			validate = false
		}
	}
	return validate
}

// Function to check if user exist in database, return user privilege
func GetUserPrivilege(username string) int {
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return 0
	}
	defer db.Close()

	var expected_user User

	row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where username = ?", username)
	err := row.Scan(&expected_user.ID, &expected_user.Privilege, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		return 0
	}

	return expected_user.Privilege
}

// function to make a post
func CreatePost(userID int, title string, content string) string {
	liked := 0
	disliked := 0
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()

	//insert into post table
	_, err := db.Exec("insert into post(user_id,title,content,created_at ,liked_no,disliked_no) values(?,?,?,?,?,?)", userID, title, content, time.Now().Format("2006-01-02 15:04:05"), liked, disliked)
	if err != nil {
		log.Fatal(err)
	}

	return "200 OK"
}

// Create dummy Post Need category name and user id for dummy
func CreateDummyPost(Category string, userID int) string {
	content := "Be the first to post in this category!"
	title := "Welcome to the " + Category + " category!"
	liked := 0
	disliked := 0
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()

	//insert into post table
	_, err := db.Exec("insert into post(user_id,title,content,created_at ,liked_no,disliked_no) values(?,?,?,?,?,?)", userID, title, content, time.Now().Format("2006-01-02 15:04:05"), liked, disliked)
	if err != nil {
		log.Fatal(err)
	}

	return "200 OK"
}

// Set relation between post and category
func SetCategoryRelation(postID int, categoryID int) string {
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()

	//insert into post table
	_, err := db.Exec("insert into category_relation(post_id,category_id) values(?,?)", postID, categoryID)
	if err != nil {
		log.Fatal(err)
	}

	return "200 OK"
}

// Get category id from category name
func GetCategoryID(category string) int {
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return 0
	}
	defer db.Close()

	var expected_category Category

	row := db.QueryRow("select id,category_name from category where category_name = ?", category)
	err := row.Scan(&expected_category.ID, &expected_category.CategoryName)
	if err != nil {
		return 0
	}

	return expected_category.ID
}

// function get last post id
func GetLastPostID() int {
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return 0
	}
	defer db.Close()

	var expected_post Post

	row := db.QueryRow("select id,user_id,title,content,created_at,liked_no,disliked_no from post order by id desc limit 1")
	err := row.Scan(&expected_post.ID, &expected_post.UserID, &expected_post.Title, &expected_post.Content, &expected_post.CreatedAt, &expected_post.LikedNumber, &expected_post.DislikedNumber)
	if err != nil {
		return 0
	}

	return expected_post.ID
}

// create a new category and add it to database
func CreateCategory(new_category_name, description string) string {
	//open database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return mess
	}
	defer db.Close()
	//insert into category table
	_, err := db.Exec("insert into category(category_name,descript,created_at) values(?,?,?)", new_category_name, description, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	return "200 OK"
}

func GetAllPosts(categoryID int) ([]Post, string) {
	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the posts
	var post_id int
	var post Post
	var posts []Post

	// Selecting all fields from the post table
	rows, err := db.Query("select post_id from category_relation where category_id = ?", categoryID)
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each row and saving the returned post
	for rows.Next() {
		err := rows.Scan(&post_id)
		row := db.QueryRow("select id,user_id,title,content,created_at from post where id = ?", post_id)
		err = row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Appending the returned post to the posts list
		posts = append(posts, post)
	}

	return posts, "200 OK"
}

func GetCommentNumberPerPost(posts []Post) (map[Post]int, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the comments numbers
	comments_numbers := make(map[Post]int)
	comments_number := 0

	// Looping through each post and counting the total number of comments
	for i := 0; i < len(posts); i++ {

		// Getting the number of comments of this post
		row := db.QueryRow("select count(*) from relation where post_id = ?", posts[i].ID)
		err := row.Scan(&comments_number)
		if err != nil {
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}

		// Mapping this category to the total number of posts
		comments_numbers[posts[i]] = comments_number
	}

	return comments_numbers, "200 OK"
}

func GetLastCommentPerPost(posts []Post) (map[Post]Comment, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the last comment
	var last_comment Comment
	last_comments := make(map[Post]Comment)

	// Looping through each post and getting its last comment
	for i := 0; i < len(posts); i++ {

		// Getting the last comment given the comment id from the last post
		row := db.QueryRow("select id,user_id,content,created_at,liked_no,disliked_no from comment where post_id = ? order by id desc limit 1", posts[i].ID)
		err := row.Scan(&last_comment.ID, &last_comment.UserID, &last_comment.Content, &last_comment.CreatedAt, &last_comment.LikedNumber, &last_comment.DislikedNumber)
		if err != nil {
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Mapping this post to the returned comment
		last_comments[posts[i]] = last_comment
	}

	return last_comments, "200 OK"
}

func GetCreatingUserPerComment(last_comments map[Post]Comment) (map[Comment]User, string) {

	// Opening the database
	db, mess := OpenDatabase()
	if strings.Compare(mess, "200 OK") != 0 {
		return nil, mess
	}
	defer db.Close()

	// Creating a variable to hold the creating users
	var creating_user User
	creating_users := make(map[Comment]User)

	// Looping through each comment and getting its creating user
	for _, last_comment := range last_comments {

		// Reading the only row and saving the returned user
		row := db.QueryRow("select id,privilege,username,passwrd,email,created_at from user where id = ?", last_comment.UserID)
		err := row.Scan(&creating_user.ID, &creating_user.Privilege, &creating_user.Username, &creating_user.Password, &creating_user.Email, &creating_user.CreatedAt)
		if err != nil {
			//fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
		}

		// Mapping this post to the returned user
		creating_users[last_comment] = creating_user
	}

	return creating_users, "200 OK"
}
