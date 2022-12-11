package forum

import "time"

type User struct {
	ID        int
	Privilege int
	Username  string
	Password  string
	Email     string
	CreatedAt string
}

type Category struct {
	ID           int
	CategoryName string
	Description  string
	CreatedAt    string
}

type Comment struct {
	ID             int
	UserID         int
	UserName       string
	PostID         int
	Content        string
	CreatedAt      time.Time
	Date           string
	LikedNumber    int
	DislikedNumber int
}

type Post struct {
	ID             int
	UserID         int
	UserName       string
	Title          string
	Content        string
	CreatedAt      time.Time
	Date           string
	LikedNumber    int
	DislikedNumber int
	URL            string
}

type Reaction struct {
	ID        int
	UserID    int
	PostID    int
	IsLiked   int
	CreatedAt string
}

type Relation struct {
	ID         int
	CategoryID int
	PostID     int
}

type Credentials struct {
	Username string
	Password string
}

type Session struct {
	Username    string
	Privilege   int
	Cookie      string
	ExpiredTime time.Time
}

type IndexPage struct {
	Username      string
	ErrorMessage  string
	Categories    []Category
	PostsNumbers  map[Category]int
	LastPosts     map[Category]Post
	CreatingUsers map[Post]User
}

type UserPage struct {
	Username      string
	ErrorMessage  string
	CreatedPosts  []Post
	LikedPosts    []Post
	CreatingUsers map[Post]User
}

type PostPage struct {
	Username       string
	ErrorMessage   string
	PostInfo       Post
	Comments       []Comment
	IsPostLiked    bool
	IsPostDisliked bool
	CommentLike    map[int]bool
	CommentDislike map[int]bool
}

type FormErrorPage struct {
	Username                string
	ErrorMessage            string
	PrivilegeErrorMessage   string
	UserErrorMessage        string
	EmailErrorMessage       string
	PasswordErrorMessage    string
	CategoryErrorMessage    string
	DescriptionErrorMessage string
}

type CategoryPage struct {
	Username       string
	ErrorMessage   string
	Posts          []Post
	CommentNumbers map[Post]int
	LastComment    map[Post]Comment
	CreatingUsers  map[Comment]User
}
