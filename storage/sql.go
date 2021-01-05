package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Sql struct {
	DB *sql.DB
}

// GetIntervals returns all intervals stored for a user
// Returns an error, if there are problems while reading the data
func (s *Sql) GetIntervals(userId UserId) ([]Interval, error) {
	var intervals []Interval

	q := `
SELECT start_time, end_time, tags, annotation
FROM interval
WHERE user_id == $1
`
	rows, err := s.DB.Query(q, userId)
	if err != nil {
		return nil, fmt.Errorf("sql_storage: Error during SQL Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		interval := Interval{}
		err = rows.Scan(&interval.Start, &interval.End, &interval.Tags, &interval.Annotation)
		if err != nil {
			return nil, fmt.Errorf("sql_storage: Error while reading database row: %w", err)
		}

		intervals = append(intervals, interval)
	}

	return intervals, nil
}

// SetIntervals replaces all intervals stored for a user
// Returns an error if an error occurs while replacing the data
func (s *Sql) SetIntervals(userId UserId, intervals []Interval) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while starting transaction: %w", err)
	}

	q := `
DELETE FROM interval
WHERE user_id = $1
`
	_, err = tx.ExecContext(ctx, q, userId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("sql_storage: Unable to rollback: %v", rollbackErr)
			return err
		}
		return err
	}

	q = `
INSERT INTO interval (user_id, start_time, end_time, tags, annotation)
VALUES ($1, $2, $3, $4, $5)
`
	for _, interval := range intervals {
		_, err = tx.ExecContext(ctx, q, userId, interval.Start, interval.End, interval.Tags, interval.Annotation)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("sql_storage: Unable to rollback: %v", rollbackErr)
				return err
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql_storage: Error during commit: %w", err)
	}

	return nil
}

// AddInterval adds a single interval to the intervals stored for a user
// Returns an error if an error occurs while adding the interval
func (s *Sql) AddInterval(userId UserId, interval Interval) error {
	q := `
INSERT INTO interval (user_id, start_time, end_time, tags, annotation)
VALUES ($1, $2, $3, $4, $5)
`
	_, err := s.DB.Exec(q, userId, interval.Start, interval.End, interval.Tags, interval.Annotation)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while adding interval: %w", err)
	}

	return nil
}

// RemoveInterval removes a single interval from the intervals stored for a user
// Returns an error if an error occurs while deleting the interval
func (s *Sql) RemoveInterval(userId UserId, interval Interval) error {
	q := `
DELETE FROM interval
WHERE user_id = $1 AND start_time = $2 AND end_time = $3 AND tags = $4 AND annotation = $5
`
	_, err := s.DB.Exec(q, userId, interval.Start, interval.End, interval.Tags, interval.Annotation)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while removing interval: %w", err)
	}

	return nil
}

func (s *Sql) Setup() {
	q := `
CREATE TABLE IF NOT EXISTS interval (
    user_id integer NOT NULL,
    start_time datetime NOT NULL,
    end_time datetime NOT NULL,
    tags text,
    annotation text,
    PRIMARY KEY (user_id, start_time, end_time, tags, annotation)
);
`
	_, err := s.DB.Exec(q)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}
}
