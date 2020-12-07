package sync

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseJSON(t *testing.T) {
	testMessage := parse`{
    userId: 1,
    clientId: 1,
    intervalData: [
        "2020-11.data",
        "2020-12.data"
    ]
}`

	parsedMessage, err := ParseJSON(testMessage)

	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}

	expectedMessage := RequestData{
		userId:   1,
		clientId: 1,
		intervalData: []string{"2020-11.data", "2020-12.data"},
	}

	if diff := cmp.Diff(expectedMessage, parsedMessage); diff != "" {
		t.Errorf("ParseJSON() missmatch (-want +got):\n%s", diff)
	}
}
