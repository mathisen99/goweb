package forum

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt string
}

type Category struct {
	ID           int
	CategoryName string
	CreatedAt    string
}

type Comment struct {
	ID             int
	UserID         int
	PostID         int
	Content        string
	CreatedAt      string
	LikedNumber    int
	DislikedNumber int
}

type Post struct {
	ID             int
	UserID         int
	Title          string
	Content        string
	CreatedAt      string
	LikedNumber    int
	DislikedNumber int
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
	ExpiredTime time.Time
}

type IndexPage struct {
	Username      string
	Categories    []Category
	PostsNumbers  map[Category]int
	LastPosts     map[Category]Post
	CreatingUsers map[Post]User
}
