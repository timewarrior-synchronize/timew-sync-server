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
