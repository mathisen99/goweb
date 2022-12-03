package web

// User struct
type User struct {
	Id        int
	Username  string
	Password  string
	Email     string
	CreatedAt string
}

type Category struct {
	Id        int
	Category  string
	CreatedAt string
}

type Comment struct {
	Id          int
	user_id     int
	post_id     int
	content     string
	createdAt   string
	liked_no    int
	disliked_no int
}

type Post struct {
	Id          int
	user_id     int
	category_id int
	title       string
	content     string
	createdAt   string
	liked_no    int
	disliked_no int
}

type Reaction struct {
	Id        int
	user_id   int
	post_id   int
	is_liked  int
	createdAt string
}

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
