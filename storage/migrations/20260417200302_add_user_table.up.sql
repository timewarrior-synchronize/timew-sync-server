CREATE TABLE user (
    user_id integer PRIMARY KEY AUTOINCREMENT
);

INSERT INTO user SELECT DISTINCT user_id FROM interval;

-- SQLite doesn't support adding foreign keys, so we create a new temporary table and move the data
-- https://sqlite.org/lang_altertable.html

CREATE TABLE new_interval (
    user_id integer NOT NULL,
    start_time datetime NOT NULL,
    end_time datetime NOT NULL,
    tags text,
    annotation text,
    PRIMARY KEY (user_id, start_time, end_time, tags, annotation),
    FOREIGN KEY (user_id) REFERENCES user (user_id) ON DELETE CASCADE
);

INSERT INTO new_interval (user_id, start_time, end_time, tags, annotation)
SELECT user_id, start_time, end_time, tags, annotation FROM interval;

DROP TABLE interval;

ALTER TABLE new_interval RENAME TO interval;
