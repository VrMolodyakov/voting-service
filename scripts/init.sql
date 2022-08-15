CREATE TABLE vote(
    vote_id SERIAL PRIMARY KEY,
    vote_title VARCHAR(200) NOT NULL
);
CREATE TABLE choice(
    choice_title VARCHAR(200) PRIMARY KEY,
    count int,
    vote_id INT REFERENCES vote(vote_id) ON DELETE CASCADE 
);