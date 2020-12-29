package storage

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSql_GetIntervals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := []Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       "Tag1 Tag2",
			Annotation: "",
		},
		{
			Tags:       "Tag3 Tag4",
			Annotation: "Annotation",
		},
	}

	q := `
SELECT start_time, end_time, tags, annotation
FROM interval
WHERE user_id == ?
`
	columns := []string{"start_time", "end_time", "tags", "annotation"}
	mock.ExpectQuery(q).
		WithArgs(4).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(time.Time{}, time.Time{}, "Tag1 Tag2", "").
				AddRow(time.Time{}, time.Time{}, "Tag3 Tag4", "Annotation"))

	sql := Sql{DB: db}
	result, err := sql.GetIntervals(4)
	if err != nil {
		t.Errorf("Error '%s' during GetIntervals", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Results differ from expected:\n%s", diff)
	}
}

func TestSql_SetIntervals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer db.Close()

	testData := []Interval{
		{
			Start:      time.Date(2020, 12, 29, 20, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 12, 29, 23, 0, 0, 0, time.UTC),
			Tags:       "Tag1 Tag2",
			Annotation: "Annotation",
		},
		{
			Tags:       "Tag3 Tag4",
			Annotation: "Annotation2",
		},
	}

	mock.ExpectBegin()
	q := `
DELETE FROM interval
WHERE user_id = \$1
`
	mock.ExpectExec(q).WithArgs(42).WillReturnResult(sqlmock.NewResult(0, 0))

	q = `
INSERT INTO interval
`
	mock.ExpectExec(q).
		WithArgs(42, testData[0].Start, testData[0].End, testData[0].Tags, testData[0].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(q).
		WithArgs(42, testData[1].Start, testData[1].End, testData[1].Tags, testData[1].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	sql := Sql{DB: db}
	err = sql.SetIntervals(42, testData)
	if err != nil {
		t.Errorf("Error '%s' during SetIntervals", err)
	}
}
