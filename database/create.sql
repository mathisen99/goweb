CREATE TABLE user (
 id INTEGER NOT NULL PRIMARY KEY,
 username VARCHAR(30) NOT NULL,
 passwrd VARCHAR(30) NOT NULL,
 email VARCHAR(30) NOT NULL,
 created_at DATETIME NOT NULL
);

CREATE TABLE category (
 id INTEGER NOT NULL PRIMARY KEY,
 category_name VARCHAR(30) NOT NULL,
 created_at DATETIME NOT NULL
);

CREATE TABLE post (
 id INTEGER NOT NULL PRIMARY KEY,
 user_id INTEGER NOT NULL,
 title VARCHAR(30) NOT NULL,
 content TEXT NOT NULL,
 created_at DATETIME NOT NULL,
 liked_no INTEGER,
 disliked_no INTEGER,
 FOREIGN KEY(user_id) REFERENCES user(id)
);

CREATE TABLE comment (
 id INTEGER NOT NULL PRIMARY KEY,
 user_id INTEGER NOT NULL,
 post_id INTEGER NOT NULL,
 content TEXT NOT NULL,
 created_at DATETIME NOT NULL,
 liked_no INTEGER,
 disliked_no INTEGER,
 FOREIGN KEY(user_id) REFERENCES user(id),
 FOREIGN KEY(post_id) REFERENCES post(id)
);

CREATE TABLE user_reaction (
 id INTEGER NOT NULL PRIMARY KEY,
 user_id INTEGER NOT NULL,
 post_id INTEGER NOT NULL,
 is_liked TINYINT(1) NOT NULL,
 created_at DATETIME NOT NULL,
 FOREIGN KEY(user_id) REFERENCES user(id),
 FOREIGN KEY(post_id) REFERENCES post(id)
);

CREATE TABLE category_relation (
 id INTEGER NOT NULL PRIMARY KEY,
 category_id INTEGER NOT NULL,
 post_id INTEGER NOT NULL,
 FOREIGN KEY(category_id) REFERENCES category(id),
 FOREIGN KEY(post_id) REFERENCES post(id)
);
