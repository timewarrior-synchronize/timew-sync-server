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
	"time"
)

// JSONRequest represents the JSON structure of a sync request.
// It contains the unique client and user ids who are syncing
// and all their tracked intervals as strings.
// It is (and should) only be used for JSON parsing.
type JSONRequest struct {
	UserId       int      `json:"userID"`
	ClientId     int      `json:"clientId"`
	IntervalData []string `json:"intervalData"`
}

// SyncRequest represents a sync request.
// It contains the id of the user and client, who are syncing
// and the intervals tracked on the client.
type SyncRequest struct {
	UserId int
	ClientId int
	Intervals []Interval
}

// Interval represents a timewarrior interval.
// It contains a start and end time and the tags associated
// with the interval.
type Interval struct {
	Start time.Time
	End time.Time
	Tags []string
}

// ParseSyncRequest parses the JSON of a sync request into a
// JSONRequest struct.
func ParseSyncRequest(jsonInput string) (SyncRequest, error) {
	var requestData JSONRequest

	err := json.Unmarshal([]byte(jsonInput), &requestData)
	syncRequest := SyncRequest{
		UserId:    requestData.UserId,
		ClientId:  requestData.ClientId,
		Intervals: StringsToIntervals(requestData.IntervalData),
	}

	return syncRequest, err
}

// ResponseData represents a sync response
// It contains the new interval for the client
type ResponseData struct {
	IntervalData []string `json:"intervalData"`
}

// ToJSON creates JSON for response body from interval data and returns it as string
func ToJSON(data []Interval) (string, error) {
	var responseData ResponseData
	responseData.IntervalData = IntervalsToStrings(data)
	result, err := json.Marshal(responseData)
	return string(result), err
}
