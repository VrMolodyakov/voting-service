CREATE TABLE vote(
    vote_id SERIAL PRIMARY KEY NOT NULL,
    vote_title VARCHAR(200) NOT NULL
);

CREATE TABLE vote_choice(
    choice_title VARCHAR(200) PRIMARY KEY,
    count int,
    vote_id INT REFERENCES vote(vote_id) ON DELETE UPDATE 
);