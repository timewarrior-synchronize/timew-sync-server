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
			{
				"start": "20200301T120000Z",
				"end": "20200301T153000Z",
				"tags": ["prank", "add"],
				"annotation": "annotation 1"
			}
		],
		"removed": [
			{
				"start": "20200401T120000Z",
				"end": "20200401T153000Z",
				"tags": ["prank", "remove"],
				"annotation": "all your codebase are belong to us"
			}
		]
	}`

	expected := SyncRequest{
		UserID: 1,
		Added: []Interval{
			{
				Start:      time.Date(2020, time.March, 1, 12, 0, 0, 0, time.UTC),
				End:        time.Date(2020, time.March, 1, 15, 30, 0, 0, time.UTC),
				Tags:       []string{"prank", "add"},
				Annotation: "annotation 1",
			},
		},
		Removed: []Interval{
			{
				Start:      time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
				End:        time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC),
				Tags:       []string{"prank", "remove"},
				Annotation: "all your codebase are belong to us",
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
			Start:      time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC),
			Tags:       []string{"prank", "laugh"},
			Annotation: "Sample Annotation",
		},
	}

	expected := `{"conflictsOccurred":false,"intervals":[{"start":"20200401T120000Z","end":"20200401T153000Z","tags":["prank","laugh"],"annotation":"Sample Annotation"}]}`

	result, err := ToJSON(testInput)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Result differs from expected: \n%s", diff)
	}

}
