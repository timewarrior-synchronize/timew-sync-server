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

// Interval represents a timewarrior interval.
// It contains a start and end time and the tags associated
// with the interval.
type Interval struct {
	Start time.Time
	End time.Time
	Tags []string
}

// ParseInterval parses a string in timewarrior format into an Interval struct.
func ParseInterval(intervalString string) Interval {
	tokens := strings.Fields(intervalString) // tokens should be ["inc", startTime, "-", endTime] for no tags and
	// ["inc", startTime, "-", endTime, "#", tag1, tag2, ..., tagN] for N > 0 tags
	startString := tokens[1]
	endString := tokens[3]
	var tags []string
	if len(tokens) > 4 { // prepare tags, iff data[i] contains tags
		tags = tokens[5:]
	} else { // initialize to empty slice else
		tags = []string{}
	}

	interval := Interval{}

	startTime, err := time.Parse(timeLayout, startString) // time.Parse uses UTC as default
	if err != nil {
		log.Printf("Error parsing start time of interval %v", intervalString)
	}
	interval.Start = startTime
	endTime, err := time.Parse(timeLayout, endString)
	if err != nil {
		log.Printf("Error parsing end time of interval %v", interval)
	}
	interval.End = endTime
	interval.Tags = make([]string, len(tags), len(tags))
	copy(interval.Tags, tags)

	return interval
}

// Serialize converts an Interval struct into a string in timewarrior format.
func (interval Interval) Serialize() string {
	startTime := interval.Start.Format(timeLayout)
	endTime := interval.End.Format(timeLayout)

	intervalTime := fmt.Sprintf("inc %v - %v", startTime, endTime)

	if len(interval.Tags) > 0 {
		intervalTime = fmt.Sprintf("%v # %v", intervalTime, strings.Join(interval.Tags, " "))
	}

	return intervalTime
}

// A String method for the interval struct.
// This convenience method implements the fmt.Stringer interface
// and converts the Interval into a human readable format.
func (interval Interval) String() string {
	return interval.Serialize()
}

// StringsToIntervals converts a slice of strings (each string encoding one time interval) to a slice of the
// corresponding interval structs
// The LastModified timestamps are initialized to time.Now()
func StringsToIntervals(data []string) []Interval {
	result := make([]Interval, len(data), len(data))
	for i, intervalString := range data {
		result[i] = ParseInterval(intervalString)
	}

	return result
}

// IntervalsToStrings converts a slice of Interval structs to a slice of the corresponding timewarrior interval strings
// Important: the LastModified information is not contained in the string representation
func IntervalsToStrings(intervals []Interval) []string {
	result := make([]string, len(intervals), len(intervals))

	for i, element := range intervals {
		result[i] = element.Serialize()
	}
	return result
}
