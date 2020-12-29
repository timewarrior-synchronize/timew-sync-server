package storage

import (
	"database/sql"
	"fmt"
)

type Sql struct {
	DB *sql.DB
}

func (sql *Sql) GetIntervals(userId UserId) ([]Interval, error) {
	var intervals []Interval

	q := `
SELECT start_time, end_time, tags, annotation
FROM interval
WHERE user_id == ?
`
	rows, err := sql.DB.Query(q, userId)
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

func (sql *Sql) SetIntervals(userId UserId, intervals []Interval) {
	panic("implement me")
}

func (sql *Sql) AddInterval(userId UserId, interval Interval) {
	panic("implement me")
}

func (sql *Sql) RemoveInterval(userId UserId, interval *Interval) {
	panic("implement me")
}
