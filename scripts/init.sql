CREATE TABLE vote(
    pool_id SERIAL PRIMARY KEY NOT NULL,
    vote_title VARCHAR(200) NOT NULL
);

CREATE TABLE vote_choice(
    choice_title VARCHAR(200) PRIMARY KEY,
    count int,
    pool_id INT REFERENCES vote(pool_id) ON DELETE UPDATE 
);