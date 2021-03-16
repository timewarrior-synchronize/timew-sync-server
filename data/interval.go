/*
Copyright 2020 - 2021, Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

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
	"strings"
	"time"
)

// layout used by timewarrior. Needed for time conversion. See e.g. time.Parse
const timeLayout = "20060102T150405Z"

// Interval represents a timewarrior interval.
// It contains a start and end time and the tags associated
// with the interval.
type Interval struct {
	Start      time.Time
	End        time.Time
	Tags       []string
	Annotation string
}

// JSONInterval represents the JSON structure of an interval. As we
// represent times in strings, these need to be converted to and from
// time.Time instances.
type JSONInterval struct {
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Tags       []string `json:"tags"`
	Annotation string   `json:"annotation"`
}

// Converts this JSON interval representation into our internal
// representation of an interval. An error might occur during parsing
// of either the start or end time.
func (json JSONInterval) ToInterval() (Interval, error) {
	start, err := time.Parse(timeLayout, json.Start)
	if err != nil {
		return Interval{}, fmt.Errorf("Error while start time: %v", err)
	}

	end, err := time.Parse(timeLayout, json.End)
	if err != nil {
		return Interval{}, fmt.Errorf("Error while end time: %v", err)
	}

	return Interval{
		Start:      start,
		End:        end,
		Tags:       json.Tags,
		Annotation: json.Annotation,
	}, nil
}

// Converts this interval into an instance of a JSONInterval struct
// which can be marshalled to JSON.
func (interval Interval) ToJSONInterval() JSONInterval {
	start := interval.Start.Format(timeLayout)
	end := interval.End.Format(timeLayout)

	return JSONInterval{
		Start:      start,
		End:        end,
		Tags:       interval.Tags,
		Annotation: interval.Annotation,
	}
}

// Convenience wrapper around ToInterval() which batch processes a
// slice of JSONInterval
func FromJSONIntervals(intervals []JSONInterval) ([]Interval, error) {
	result := make([]Interval, len(intervals))

	for i, x := range intervals {
		interval, err := x.ToInterval()

		if err != nil {
			return nil, err
		}

		result[i] = interval
	}

	return result, nil
}

// Convenience wrapper around ToJSONInterval() which batch processes a
// slice of Interval
func ToJSONIntervals(intervals []Interval) []JSONInterval {
	result := make([]JSONInterval, len(intervals))

	for i, x := range intervals {
		result[i] = x.ToJSONInterval()
	}

	return result
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

// IntervalsToStrings converts a slice of Interval structs to a slice of the corresponding timewarrior interval strings
// Important: the LastModified information is not contained in the string representation
func IntervalsToStrings(intervals []Interval) []string {
	result := make([]string, len(intervals), len(intervals))

	for i, element := range intervals {
		result[i] = element.Serialize()
	}
	return result
}
