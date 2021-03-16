/*
Copyright 2020 - 2021, Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/timewarrior-synchronize/timew-sync-server/data"
	"log"
)

type Sql struct {
	LockerRoom
	DB *sql.DB
}

// Initialize runs all necessary setup for this Storage instance
func (s *Sql) Initialize() error {
	s.InitializeLockerRoom()

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
	return nil
}

// GetIntervals returns all intervals stored for a user
// Returns an error, if there are problems while reading the data
func (s *Sql) GetIntervals(userId UserId) ([]data.Interval, error) {
	var intervals []IntervalKey

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
		interval := IntervalKey{}
		err = rows.Scan(&interval.Start, &interval.End, &interval.Tags, &interval.Annotation)
		if err != nil {
			return nil, fmt.Errorf("sql_storage: Error while reading database row: %w", err)
		}

		intervals = append(intervals, interval)
	}

	return ConvertToIntervals(intervals), nil
}

// SetIntervals replaces all intervals stored for a user
// Returns an error if an error occurs while replacing the data
func (s *Sql) SetIntervals(userId UserId, intervals []data.Interval) error {
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
	keys := ConvertToKeys(intervals)
	for _, key := range keys {
		_, err = tx.ExecContext(ctx, q, userId, key.Start, key.End, key.Tags, key.Annotation)
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
func (s *Sql) AddInterval(userId UserId, interval data.Interval) error {
	q := `
INSERT OR IGNORE INTO interval (user_id, start_time, end_time, tags, annotation)
VALUES ($1, $2, $3, $4, $5)
`
	key := IntervalToKey(interval)
	_, err := s.DB.Exec(q, userId, key.Start, key.End, key.Tags, key.Annotation)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while adding interval: %w", err)
	}

	return nil
}

// RemoveInterval removes a single interval from the intervals stored for a user
// Returns an error if an error occurs while deleting the interval
func (s *Sql) RemoveInterval(userId UserId, interval data.Interval) error {
	q := `
DELETE FROM interval
WHERE user_id = $1 AND start_time = $2 AND end_time = $3 AND tags = $4 AND annotation = $5
`
	key := IntervalToKey(interval)
	_, err := s.DB.Exec(q, userId, key.Start, key.End, key.Tags, key.Annotation)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while removing interval: %w", err)
	}

	return nil
}

// ModifyIntervals atomically adds and deletes a specified set of
// intervals. Returns an error if an error occurs while modifying the
// data
func (s *Sql) ModifyIntervals(userId UserId, add []data.Interval, del []data.Interval) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql_storage: Error while starting transaction: %w", err)
	}

	// Delete the specified intervals
	q := `
DELETE FROM interval
WHERE user_id = $1 AND start_time = $2 AND end_time = $3 AND tags = $4 AND annotation = $5
`
	keysToDelete := ConvertToKeys(del)
	for _, key := range keysToDelete {
		_, err = tx.ExecContext(ctx, q, userId, key.Start, key.End, key.Tags, key.Annotation)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("sql_storage: Unable to rollback: %v", rollbackErr)
				return err
			}
			return err
		}
	}

	// Add the specified intervals
	q = `
INSERT OR IGNORE INTO interval (user_id, start_time, end_time, tags, annotation)
VALUES ($1, $2, $3, $4, $5)
`
	keysToAdd := ConvertToKeys(add)
	for _, key := range keysToAdd {
		_, err = tx.ExecContext(ctx, q, userId, key.Start, key.End, key.Tags, key.Annotation)
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
