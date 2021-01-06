package data

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestParseJSON(t *testing.T) {
	testInput := `
	{
		"userID": 1,
		"added": [
			"inc 20200301T120000Z - 20200301T153000Z # prank add"
		],
		"removed": [
			"inc 20200401T120000Z - 20200401T153000Z # prank remove"
		]
	}`

	expected := SyncRequest{
		UserID: 1,
		Added: []Interval{
			{
				Start: time.Date(2020, time.March, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.March, 1, 15, 30, 0, 0, time.UTC),
				Tags:  []string{"prank", "add"},
			},
		},
		Removed: []Interval{
			{
				Start: time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC),
				Tags:  []string{"prank", "remove"},
			},
		},
	}

	result, err := ParseSyncRequest(testInput)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Result differs from expected: \n%s", diff)
	}

}

func TestToJSON(t *testing.T) {
	testInput := []Interval{
		{
			Start: time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
			End:   time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC),
			Tags:  []string{"prank", "laugh"},
		},
	}

	expected := `{"intervalData":["inc 20200401T120000Z - 20200401T153000Z # prank laugh"]}`

	result, err := ToJSON(testInput)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Result differs from expected: \n%s", diff)
	}

}
