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

package sync

import "encoding/json"

// RequestData represents a sync request.
// It contains the unique client and user id's who are syncing
// and all their tracked intervals.
type RequestData struct {
	UserId       int      `json:"userID"`
	ClientId     int      `json:"clientId"`
	IntervalData []string `json:"intervalData"`
}

// ParseSyncRequest parses the JSON of a sync request into a
// RequestData struct.
func ParseSyncRequest(jsonInput string) (RequestData, error) {
	var requestData RequestData

	err := json.Unmarshal([]byte(jsonInput), &requestData)

	return requestData, err
}

type ResponseData struct {
	IntervalData []string `json:"intervalData"`
}

func ToJSON(data []string) (string, error) {
	var responseData ResponseData
	responseData.IntervalData = data
	result, err := json.Marshal(responseData)
	return string(result), err
}
