package sync

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseJSON(t *testing.T) {
	testInput := `
	{
		"userId": 1,
		"clientId": 1,
		"intervalData": [
			"2020-11.data",
			"2020-12.data"
		]
	}`

	expected := RequestData{
		UserId:       1,
		ClientId:     1,
		IntervalData: []string{"2020-11.data", "2020-12.data"},
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
	testInput := []string{"2020-11.data", "2020-12.data"}

	expected := `{"intervalData":["2020-11.data","2020-12.data"]}`

	result, err := ToJSON(testInput)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Result differs from expected: \n%s", diff)
	}

}
