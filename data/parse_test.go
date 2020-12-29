package data

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestParseJSON(t *testing.T) {
	testInput := `
	{
		"userId": 1,
		"clientId": 1,
		"intervalData": [
			"inc 20200401T120000Z - 20200401T153000Z # prank laugh"
		]
	}`

	expected := SyncRequest{
		UserId:   1,
		ClientId: 1,
		Intervals: []Interval{
			{
				Start: time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC),
				Tags:  []string{"prank", "laugh"},
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
