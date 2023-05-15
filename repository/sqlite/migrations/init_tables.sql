CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(50) UNIQUE NOT NULL,
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
);


CREATE TABLE IF NOT EXISTS session(
    user_id INTEGER UNIQUE NOT NULL,
    cookie VARCHAR UNIQUE NOT NULL,
    expire_at DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS posts(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    category VARCHAR,
    FOREIGN KEY (user_id)REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS votes(
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote BOOLEAN NOT NULL,
    UNIQUE(post_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY(post_id)REFERENCES posts(id)
);

CREATE TABLE IF NOT EXISTS comments(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content VARCHAR NOT NULL,
    FOREIGN KEY (user_id)REFERENCES users(id),
    FOREIGN key (post_id)REFERENCES posts(id)
); 

CREATE TABLE IF NOT EXISTS votes_comm(
    comm_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote BOOLEAN NOT NULL,
    UNIQUE(comm_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY(comm_id)REFERENCES comments(id)
);