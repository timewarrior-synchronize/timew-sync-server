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
	"fmt"
	"log"
	"strings"
	"time"
)

// layout used by timewarrior. Needed for time conversion. See e.g. time.Parse
const timeLayout = "20060102T150405Z"

// StringsToIntervals converts a slice of strings (each string encoding one time interval) to a slice of the
// corresponding interval structs
// The LastModified timestamps are initialized to time.Now()
func StringsToIntervals(data []string) []Interval {
	result := make([]Interval, len(data), len(data))
	for i, element := range data {
		tokens := strings.Fields(element) // tokens should be ["inc", startTime, "-", endTime] for no tags and
		// ["inc", startTime, "-", endTime, "#", tag1, tag2, ..., tagN] for N > 0 tags
		startString := tokens[1]
		endString := tokens[3]
		var tags []string
		if len(tokens) > 4 { // prepare tags, iff data[i] contains tags
			tags = tokens[5:]
		} else { // initialize to empty slice else
			tags = []string{}
		}
		result[i] = Interval{}
		startTime, err := time.Parse(timeLayout, startString) // time.Parse uses UTC as default
		if err != nil {
			log.Printf("Error parsing start time of interval %v", element)
		}
		result[i].Start = startTime
		endTime, err := time.Parse(timeLayout, endString)
		if err != nil {
			log.Printf("Error parsing end time of interval %v", element)
		}
		result[i].End = endTime
		result[i].Tags = make([]string, len(tags), len(tags))
		copy(result[i].Tags, tags)
	}

	return result
}

// IntervalsToStrings converts a slice of Interval structs to a slice of the corresponding timewarrior interval strings
// Important: the LastModified information is not contained in the string representation
func IntervalsToStrings(intervals []Interval) []string {
	result := make([]string, len(intervals), len(intervals))

	for i, element := range intervals {
		intervalTime := fmt.Sprintf("inc %v - %v", element.Start.Format(timeLayout), element.End.Format(timeLayout))
		if len(element.Tags) > 0 {
			intervalTime = fmt.Sprintf("%v # %v", intervalTime, strings.Join(element.Tags, " "))
		}
		result[i] = intervalTime
	}
	return result
}
