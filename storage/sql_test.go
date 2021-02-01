/*
Copyright 2020 - Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

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
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
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

	expected := []data.Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       []string{"Tag1", "Tag2"},
			Annotation: "",
		},
		{
			Tags:       []string{"Tag3", "Tag4"},
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
				AddRow(time.Time{}, time.Time{}, IntervalToKey(expected[0]).Tags, "").
				AddRow(time.Time{}, time.Time{}, IntervalToKey(expected[1]).Tags, "Annotation"))

	sql := Sql{DB: db}
	result, err := sql.GetIntervals(4)
	if err != nil {
		t.Errorf("Error '%s' during GetIntervals", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Results differ from expected:\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSql_SetIntervals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testData := []data.Interval{
		{
			Start:      time.Date(2020, 12, 29, 20, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 12, 29, 23, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag1", "Tag2"},
			Annotation: "Annotation",
		},
		{
			Tags:       []string{"Tag3", "Tag4"},
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
		WithArgs(42, testData[0].Start, testData[0].End, IntervalToKey(testData[0]).Tags, testData[0].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(q).
		WithArgs(42, testData[1].Start, testData[1].End, IntervalToKey(testData[1]).Tags, testData[1].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	sql := Sql{DB: db}
	err = sql.SetIntervals(42, testData)
	if err != nil {
		t.Errorf("Error '%s' during SetIntervals", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSql_AddInterval(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testData := data.Interval{
		Start:      time.Date(2003, 3, 12, 7, 20, 15, 0, time.UTC),
		End:        time.Date(2004, 2, 4, 16, 30, 43, 0, time.UTC),
		Tags:       []string{"TestTag", "TestTag2"},
		Annotation: "TestAnnotation",
	}

	q := `
INSERT INTO interval \(user_id, start_time, end_time, tags, annotation\)
VALUES \(\$1, \$2, \$3, \$4, \$5\)
`
	mock.ExpectExec(q).
		WithArgs(3, testData.Start, testData.End, IntervalToKey(testData).Tags, testData.Annotation).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sql := Sql{DB: db}
	err = sql.AddInterval(3, testData)
	if err != nil {
		t.Errorf("Error '%s' during AddInterval", err)
	}
}

func TestSql_RemoveInterval(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testData := data.Interval{
		Start:      time.Date(2030, 2, 24, 14, 23, 42, 0, time.UTC),
		End:        time.Date(2030, 2, 24, 17, 24, 0, 0, time.UTC),
		Tags:       []string{"Tag1", "Tag2", "Tag3"},
		Annotation: "Annotation",
	}

	q := `
DELETE FROM interval
WHERE user_id = \$1 AND start_time = \$2 AND end_time = \$3 AND tags = \$4 AND annotation = \$5
`
	mock.ExpectExec(q).
		WithArgs(0, testData.Start, testData.End, IntervalToKey(testData).Tags, testData.Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))

	sql := Sql{DB: db}
	err = sql.RemoveInterval(0, testData)
	if err != nil {
		t.Errorf("Error '%s' during RemoveInterval", err)
	}
}

func TestSql_ModifyIntervals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	add := []data.Interval{
		{
			Start:      time.Date(2020, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag3", "Tag4"},
			Annotation: "Annotation2",
		},
	}

	del := []data.Interval{
		{
			Start:      time.Date(2021, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag1", "Tag2"},
			Annotation: "Annotation",
		},
	}

	mock.ExpectBegin()
	q := `
DELETE FROM interval
WHERE user_id = \$1 AND start_time = \$2 AND end_time = \$3 AND tags = \$4 AND annotation = \$5
`
	mock.ExpectExec(q).
		WithArgs(123, del[0].Start, del[0].End, IntervalToKey(del[0]).Tags, del[0].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))

	q = `
INSERT INTO interval \(user_id, start_time, end_time, tags, annotation\)
VALUES \(\$1, \$2, \$3, \$4, \$5\)
`
	mock.ExpectExec(q).
		WithArgs(123, add[0].Start, add[0].End, IntervalToKey(add[0]).Tags, add[0].Annotation).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	sql := Sql{DB: db}
	err = sql.ModifyIntervals(123, add, del)
	if err != nil {
		t.Errorf("Error '%s' during SetIntervals", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
