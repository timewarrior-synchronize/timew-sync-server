CREATE TABLE IF NOT EXISTS interval (
    user_id integer NOT NULL,
    start_time datetime NOT NULL,
    end_time datetime NOT NULL,
    tags text,
    annotation text,
    PRIMARY KEY (user_id, start_time, end_time, tags, annotation)
);
