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

import (
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/storage"
	"strings"
	"time"
)

// StringsToIntervals converts a slice of Strings (each string encoding one time interval) to a slice of the
// corresponding interval structs
// The LastModified timestamps are initialized to time.Now()
func StringsToIntervals(data []string) []storage.Interval {
	now := time.Now()
	layout := "20060102T150405Z"
	result := make([]storage.Interval, len(data), len(data))

	for i, element := range data {
		tokens := strings.Fields(element)
		startString := tokens[1]
		endString := tokens[3]
		var tags []string
		if len(tokens) > 4 {
			tags = tokens[5:]
		} else {
			tags = []string{}
		}
		result[i] = storage.Interval{}
		result[i].Start, _ = time.Parse(layout, startString)
		result[i].End, _ = time.Parse(layout, endString)
		result[i].Tags = make([]string, len(tags), len(tags))
		copy(result[i].Tags, tags)
		result[i].LastModified = now
		result[i].Deleted = false
	}
	return result
}
