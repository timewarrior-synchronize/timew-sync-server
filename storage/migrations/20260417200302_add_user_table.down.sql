-- see up migration

CREATE TABLE old_interval (
    user_id integer NOT NULL,
    start_time datetime NOT NULL,
    end_time datetime NOT NULL,
    tags text,
    annotation text,
    PRIMARY KEY (user_id, start_time, end_time, tags, annotation)
);

INSERT INTO old_interval (user_id, start_time, end_time, tags, annotation)
SELECT user_id, start_time, end_time, tags, annotation FROM interval;

DROP TABLE interval;

ALTER TABLE old_interval RENAME TO interval;

DROP TABLE user;
