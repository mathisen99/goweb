package forum

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Opening the database
func OpenDatabase() (*sql.DB, string) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: DATABASE CORRUPTED"
	}

	return db, "200 OK"
}

// Getting all the categories from the database
func GetAllCategories(db *sql.DB) ([]Category, string) {

	// Creating a variable to hold the categories
	var category Category
	var categories []Category

	// Selecting all fields from the category table
	rows, err := db.Query("select id,category_name,created_at from category")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
	}
	defer rows.Close()

	// Looping through each row and saving the returned category
	for rows.Next() {
		err := rows.Scan(&category.ID, &category.CategoryName, &category.CreatedAt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: CATEGORY DATA CORRUPTED"
		}

		// Appending the returned category to the categories list
		categories = append(categories, category)
	}

	return categories, "200 OK"
}

// Getting the total number of posts per category
func GetPostsNumberPerCategory(db *sql.DB, categories []Category) (map[Category]int, string) {

	// Creating a variable to hold the posts numbers
	posts_numbers := make(map[Category]int)

	// Looping through each category and counting the total number of posts
	for i := 0; i < len(categories); i++ {

		// Getting all the posts of this category
		rows, err := db.Query("select id,category_id,post_id from category_relation where category_id = ?", categories[i].ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}
		defer rows.Close()

		// Counting the total number of posts for this category
		posts_number := 0
		for rows.Next() {
			posts_number++
		}

		// Mapping this category to the total number of posts
		posts_numbers[categories[i]] = posts_number
	}

	return posts_numbers, "200 OK"
}

// Getting the last post per category
func GetLastPostPerCategory(db *sql.DB, categories []Category) (map[Category]Post, string) {

	// Creating a variable to hold the last posts
	var last_relation Relation
	var last_post Post
	last_posts := make(map[Category]Post)

	// Looping through each category and getting its last post
	for i := 0; i < len(categories); i++ {

		// Getting the last relation of this category
		rows, err := db.Query("select id,category_id,post_id from category_relation where category_id = ? order by post_id desc limit 1", categories[i].ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}
		defer rows.Close()

		// Reading the only row and saving the returned relation
		rows.Next()
		err = rows.Scan(&last_relation.ID, &last_relation.CategoryID, &last_relation.PostID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: RELATION DATA CORRUPTED"
		}

		// Getting the last post given the post_id from the last relation
		rows, err = db.Query("select id,user_id,title,content,created_at,liked_no,disliked_no from post where id = ?", last_relation.PostID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Reading the only row and saving the returned post
		rows.Next()
		err = rows.Scan(&last_post.ID, &last_post.UserID, &last_post.Title, &last_post.Content, &last_post.CreatedAt, &last_post.LikedNumber, &last_post.DislikedNumber)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: POST DATA CORRUPTED"
		}

		// Mapping this category to the returned post
		last_posts[categories[i]] = last_post
	}

	return last_posts, "200 OK"
}

// Getting the creating user per post
func GetCreatingUserPerPost(db *sql.DB, last_posts map[Category]Post) (map[Post]User, string) {

	// Creating a variable to hold the creating users
	var creating_user User
	creating_users := make(map[Post]User)

	// Looping through each post and getting its creating user
	for _, last_post := range last_posts {
		rows, err := db.Query("select id,username,passwrd,email,created_at from user where id = ?", last_post.UserID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
		}
		defer rows.Close()

		// Reading the only row and saving the returned user
		rows.Next()
		err = rows.Scan(&creating_user.ID, &creating_user.Username, &creating_user.Password, &creating_user.Email, &creating_user.CreatedAt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
		}

		// Mapping this post to the returned user
		creating_users[last_post] = creating_user
	}

	return creating_users, "200 OK"
}

// Getting the password of the given user from the database
func GetPassword4User(db *sql.DB, username string) (string, string) {

	// Creating a variable to hold the expected user
	var expected_user User

	rows, err := db.Query("select id,username,passwrd,email,created_at from user where username = ?", username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", "500 INTERNAL SERVER ERROR: USER DATA CORRUPTED"
	}
	defer rows.Close()

	// Reading the only row and saving the returned user
	rows.Next()
	err = rows.Scan(&expected_user.ID, &expected_user.Username, &expected_user.Password, &expected_user.Email, &expected_user.CreatedAt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", "401 UNAUTHORIZED: USER NOT FOUND"
	}

	return expected_user.Password, "200 OK"
}
