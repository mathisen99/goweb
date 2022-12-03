package web

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// global time variable
// NOT the best way to do this, but it works for now
var Current_time = time.Now().Format("2006-01-02 15:04:05")

// function to open database will remove loads of code from other functions
func Open_db() *sql.DB {
	//open database
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// function to get all users from database currently returns a slice of structs
// Not sure if we want this to return a slice or a map or something else, we can talk about it
func Get_users() []User {
	//creating variable to hold users
	var user User
	var users []User
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from user table
	rows, err := db.Query("select id, username,passwrd,email,created_at from user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		//append to slice
		users = appendUser(users, user)
	}
	return users
}

// function to create user
func Create_user(username string, password string, email string) {
	//open database
	//open database
	db := Open_db()
	defer db.Close()

	//insert into user table
	_, err := db.Exec("insert into user(username,passwrd,email,created_at) values(?,?,?,?)", username, password, email, Current_time)
	if err != nil {
		log.Fatal(err)
	}
}

// function to append struct to slice
func appendUser(users []User, user User) []User {
	users = append(users, user)
	return users
}

// function to delete user by id (not sure if we want to delete by id or username)
// we can talk about it, it will depend on how we handle backend
func Delete_user(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//delete from user table
	_, err := db.Exec("delete from user where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update user by id (not sure if we want to update by id or username)
func Update_user(id int, username string, password string, email string) {
	//open database
	db := Open_db()
	defer db.Close()

	//update user table
	_, err := db.Exec("update user set username = ?, passwrd = ?, email = ? where id = ?", username, password, email, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to get single user by id (not sure if we want to get by id or username)
func Get_user(id int) User {
	//creating variable to hold user
	var user User
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from user table
	rows, err := db.Query("select id, username,passwrd,email,created_at from user where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user
}

// function to get all categories from database currently returns a slice of structs
func Get_categories() []Category {
	//creating variable to hold categories
	var category Category
	var categories []Category
	//open database
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from category table
	rows, err := db.Query("select id, category,created_at from category")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&category.Id, &category.Category, &category.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		//append to slice
		categories = appendCategory(categories, category)
	}
	return categories
}

// function to append struct to slice
func appendCategory(categories []Category, category Category) []Category {
	categories = append(categories, category)
	return categories
}

// function to create a category
func Create_category(name string) {
	//open database
	db := Open_db()
	defer db.Close()

	//insert into category table
	_, err := db.Exec("insert into category(category,created_at) values(?,?)", name, Current_time)
	if err != nil {
		log.Fatal(err)
	}
}

// function to delete category by id
func Delete_category(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//delete from category table
	_, err := db.Exec("delete from category where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update category by id
func Update_category(id int, name string) {
	//open database
	db := Open_db()
	defer db.Close()

	//update category table
	_, err := db.Exec("update category set category = ? where id = ?", name, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to get single category by id
func Get_category(id int) Category {
	//creating variable to hold category
	var category Category
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from category table
	rows, err := db.Query("select id, category,created_at from category where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&category.Id, &category.Category, &category.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return category
}

// function to create a comment
func Create_comment(user_id int, post_id int, comment string, liked int, disliked int) {
	//setting liked and disliked to 0 as this is a new comment! no need for this to be passed in but its here for now
	//ll fix this better later il just want to keep track of all variables at the moment
	liked = 0
	disliked = 0
	//open database
	db := Open_db()
	defer db.Close()

	//insert into comment table
	_, err := db.Exec("insert into comment(user_id,post_id,content,liked_no,disliked_no,created_at) values(?,?,?,?,?,?)", user_id, post_id, comment, liked, disliked, Current_time)
	if err != nil {
		log.Fatal(err)
	}
}

// function to delete comment by id
func Delete_comment(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//delete from comment table
	_, err := db.Exec("delete from comment where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update comment by id
func Update_comment(id int, comment string) {
	//open database
	db := Open_db()
	defer db.Close()

	//update comment table
	_, err := db.Exec("update comment set content = ? where id = ?", comment, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to get single comment by id
func Get_comment(id int) Comment {
	//creating variable to hold comment
	var comment Comment
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from comment table
	rows, err := db.Query("select id,user_id,post_id,content,liked_no,disliked_no,created_at from comment where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&comment.Id, &comment.user_id, &comment.post_id, &comment.content, &comment.liked_no, &comment.disliked_no, &comment.createdAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return comment
}

// function to like a comment
func Like_comment(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update comment table
	_, err := db.Exec("update comment set liked_no = liked_no + 1 where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to dislike a comment
func Dislike_comment(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update comment table
	_, err := db.Exec("update comment set disliked_no = disliked_no + 1 where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to create a post
func Create_post(user_id int, category_id int, title string, content string, liked int, disliked int) {
	//setting liked and disliked to 0 as this is a new post! no need for this to be passed in but its here for now
	liked = 0
	disliked = 0
	//open database
	db := Open_db()
	defer db.Close()

	//insert into post table
	_, err := db.Exec("insert into post(user_id,category_id,title,content,liked_no,disliked_no,created_at) values(?,?,?,?,?,?,?)", user_id, category_id, title, content, liked, disliked, Current_time)
	if err != nil {
		log.Fatal(err)
	}
}

// function to delete post by id
func Delete_post(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//delete from post table
	_, err := db.Exec("delete from post where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update post by id
func Update_post(id int, title string, content string) {
	//open database
	db := Open_db()
	defer db.Close()

	//update post table
	_, err := db.Exec("update post set title = ?, content = ? where id = ?", title, content, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to get single post by id
func Get_post(id int) Post {
	//creating variable to hold post
	var post Post
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from post table
	rows, err := db.Query("select id,user_id,category_id,title,content,liked_no,disliked_no,created_at from post where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&post.Id, &post.user_id, &post.category_id, &post.title, &post.content, &post.liked_no, &post.disliked_no, &post.createdAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return post
}

// function to get last post
func Get_last_post() Post {
	//creating variable to hold post
	var post Post
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from post table
	rows, err := db.Query("select id,user_id,category_id,title,content,liked_no,disliked_no,created_at from post order by id desc limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//save into struct
		err := rows.Scan(&post.Id, &post.user_id, &post.category_id, &post.title, &post.content, &post.liked_no, &post.disliked_no, &post.createdAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return post
}

//function to get last post and check used id and get usenme from user table

// function to get all posts
func Get_all_posts() []Post {
	//creating variable to hold posts
	var posts []Post
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from post table
	rows, err := db.Query("select id,user_id,category_id,title,content,liked_no,disliked_no,created_at from post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//creating variable to hold post
		var post Post
		//save into struct
		err := rows.Scan(&post.Id, &post.user_id, &post.category_id, &post.title, &post.content, &post.liked_no, &post.disliked_no, &post.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		//append post to posts
		posts = append(posts, post)
	}
	return posts
}

// function to like a post
func Like_post(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update post table
	_, err := db.Exec("update post set liked_no = liked_no + 1 where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to dislike a post
func Dislike_post(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update post table
	_, err := db.Exec("update post set disliked_no = disliked_no + 1 where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// here comes the user reaction part.... we may want more here right????
// function to get all reactions by user id
func Get_user_reactions(id int) []Reaction {
	//creating variable to hold reactions
	var reactions []Reaction
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from reaction table
	rows, err := db.Query("select id,user_id,post_id,is_liked,created_at from user_reaction where user_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//creating variable to hold reaction
		var reaction Reaction
		//save into struct
		err := rows.Scan(&reaction.Id, &reaction.user_id, &reaction.post_id, &reaction.is_liked, &reaction.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		//append reaction to reactions
		reactions = append(reactions, reaction)
	}
	return reactions
}

// function to create a reaction???? not about this one...
func Create_reaction(user_id int, post_id int, is_liked int) {
	//open database
	db := Open_db()
	defer db.Close()

	//insert into reaction table
	_, err := db.Exec("insert into user_reaction(user_id,post_id,is_liked,created_at) values(?,?,?,?)", user_id, post_id, is_liked, Current_time)
	if err != nil {
		log.Fatal(err)
	}
}

// function to delete reaction by id
func Delete_reaction(id int) {
	//open database
	//open database
	db := Open_db()
	defer db.Close()

	//delete from reaction table
	_, err := db.Exec("delete from user_reaction where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update reaction by id
func Update_reaction(id int, is_liked int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update reaction table
	_, err := db.Exec("update user_reaction set is_liked = ? where id = ?", is_liked, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to create category_relation
func Create_category_relation(category_id int, post_id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//insert into category_relation table
	_, err := db.Exec("insert into category_relation(category_id,post_id) values(?,?)", category_id, post_id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to delete category_relation
func Delete_category_relation(id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//delete from category_relation table
	_, err := db.Exec("delete from category_relation where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to update category_relation
func Update_category_relation(id int, category_id int, post_id int) {
	//open database
	db := Open_db()
	defer db.Close()

	//update category_relation table
	_, err := db.Exec("update category_relation set category_id = ?, post_id = ? where id = ?", category_id, post_id, id)
	if err != nil {
		log.Fatal(err)
	}
}

// function to check if user exists
func Check_user_exists(username string) bool {
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from user table
	rows, err := db.Query("select id,username,passwrd,email,created_at from user where username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//creating variable to hold user
		var user User
		//save into struct
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		//check if user exists
		if user.Id != 0 {
			return true
		}
	}
	return false
}

// function to check if user password is correct
func Check_user_password(username string, password string) bool {
	//open database
	db := Open_db()
	defer db.Close()

	//selecting all fields from user table
	rows, err := db.Query("select id,username,passwrd,email,created_at from user where username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//loop through rows and save to struct
	for rows.Next() {
		//creating variable to hold user
		var user User
		//save into struct
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		//check if user password is correct
		fmt.Println(user.Password)
		if user.Password == password {
			return true
		}
	}
	return false
}
