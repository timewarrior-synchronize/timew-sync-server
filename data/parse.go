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
	"encoding/json"
	"fmt"
)

// JSONRequest represents the JSON structure of a sync request.
// It contains the unique client id and an interval diff, stating added and removed intervals as strings.
// It is (and should) only be used for JSON parsing.
type JSONRequest struct {
	UserID  int            `json:"userID"`
	Added   []JSONInterval `json:"added"`
	Removed []JSONInterval `json:"removed"`
}

// SyncRequest represents a sync request.
// It contains the id of the user who is syncing
// and its interval diff.
type SyncRequest struct {
	UserID  int
	Added   []Interval
	Removed []Interval
}

// ParseSyncRequest parses the JSON of a sync request into a
// JSONRequest struct.
func ParseSyncRequest(jsonInput string) (SyncRequest, error) {
	var requestData JSONRequest

	err := json.Unmarshal([]byte(jsonInput), &requestData)
	if err != nil {
		return SyncRequest{}, fmt.Errorf("Error occured during JSON parse: %v", err)
	}

	added, err := FromJSONIntervals(requestData.Added)
	if err != nil {
		return SyncRequest{}, fmt.Errorf("Error occured during parsing of added intervals: %v", err)
	}

	removed, err := FromJSONIntervals(requestData.Removed)
	if err != nil {
		return SyncRequest{}, fmt.Errorf("Error occured during parsing of removed intervals: %v", err)
	}

	syncRequest := SyncRequest{
		UserID:  requestData.UserID,
		Added:   added,
		Removed: removed,
	}

	return syncRequest, err
}

// ResponseData represents a sync response
// It contains the new interval for the client
type ResponseData struct {
	ConflictsOccurred bool           `json:"conflictsOccurred"`
	Intervals         []JSONInterval `json:"intervals"`
}

// ToJSON creates JSON for response body from interval data and returns it as string
func ToJSON(data []Interval) (string, error) {
	response := ResponseData{
		// TODO: Return the real value when the conflict solving layer is implemented.
		ConflictsOccurred: false,
		Intervals:         ToJSONIntervals(data),
	}

	jsonResult, err := json.Marshal(response)

	return string(jsonResult), err
}
