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
	"reflect"
	"testing"
	"time"
)

func TestParseInterval(t *testing.T) {
	testData := "inc 20201215T225656Z - 20201215T230000Z # tag1 tag2"

	expected := Interval{
		Start: time.Date(2020, time.December, 15, 22, 56, 56, 0, time.UTC),
		End:   time.Date(2020, time.December, 15, 23, 0, 0, 0, time.UTC),
		Tags:  []string{"tag1", "tag2"},
	}

	result := ParseInterval(testData)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("Result differs from expected: \n%s", diff)
	}
}

func TestInterval_Serialize(t *testing.T) {
	testData := Interval{
		Start: time.Date(1991, time.March, 13, 3, 45, 45, 0, time.UTC),
		End:   time.Date(1991, time.March, 14, 7, 32, 56, 0, time.UTC),
		Tags:  []string{"tag1", "tag2"},
	}

	expected := "inc 19910313T034545Z - 19910314T073256Z # tag1 tag2"

	result := testData.Serialize()

	if expected != result {
		t.Errorf("Wrong interval format. Expected: \"%s\", got: \"%s\"", expected, result)
	}
}

func TestStringsToIntervals(t *testing.T) {
	testData := []string{
		"inc 20201125T093910Z - 20201125T093943Z",
		"inc 20201125T095240Z - 20201125T095253Z # test",
		"inc 20201209T140521Z - 20201209T140533Z # a b c",
	}
	loc, _ := time.LoadLocation("UTC")
	expected := make([]Interval, 3, 3)
	expected[0] = Interval{
		Start:        time.Date(2020, 11, 25, 9, 39, 10, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 39, 43, 0, loc),
		Tags:         []string{},
	}
	expected[1] = Interval{
		Start:        time.Date(2020, 11, 25, 9, 52, 40, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 52, 53, 0, loc),
		Tags:         []string{"test"},
	}
	expected[2] = Interval{
		Start:        time.Date(2020, 12, 9, 14, 5, 21, 0, loc),
		End:          time.Date(2020, 12, 9, 14, 5, 33, 0, loc),
		Tags:         []string{"a", "b", "c"},
	}
	actual := StringsToIntervals(testData)
	if len(actual) != 3 {
		t.Errorf("wrong number of intervals returned: expected 3 got %v\n", len(actual))
	}
	for i, actualInterval := range actual {
		if !actualInterval.Start.Equal(expected[i].Start) {
			t.Errorf("wrong start time for interval %v: expected %v got %v", i, expected[i].Start.String(), actualInterval.Start.String())
		}
		if !actualInterval.End.Equal(expected[i].End) {
			t.Errorf("wrong end time for interval %v: expected %v got %v", i, expected[i].End.String(), actualInterval.End.String())
		}
		if !cmp.Equal(actualInterval.Tags, expected[i].Tags) {
			t.Errorf("wrong tags for interval %v: expected %v of type %v got %v of type %v", i, expected[i].Tags, reflect.TypeOf(expected[i].Tags), actualInterval.Tags, reflect.TypeOf(actualInterval.Tags))
		}
	}
}

func TestIntervalsToStrings(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	testData := make([]Interval, 3, 3)
	testData[0] = Interval{
		Start:        time.Date(2020, 11, 25, 9, 39, 10, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 39, 43, 0, loc),
		Tags:         []string{},
	}
	testData[1] = Interval{
		Start:        time.Date(2020, 11, 25, 9, 52, 40, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 52, 53, 0, loc),
		Tags:         []string{"test"},
	}
	testData[2] = Interval{
		Start:        time.Date(2020, 12, 9, 14, 5, 21, 0, loc),
		End:          time.Date(2020, 12, 9, 14, 5, 33, 0, loc),
		Tags:         []string{"a", "b", "c"},
	}

	expected := []string{
		"inc 20201125T093910Z - 20201125T093943Z",
		"inc 20201125T095240Z - 20201125T095253Z # test",
		"inc 20201209T140521Z - 20201209T140533Z # a b c",
	}

	actual := IntervalsToStrings(testData)
	if len(actual) != 3 {
		t.Errorf("wrong number of strings returned: expected 3 got %v\n", len(actual))
	}

	for i, a := range actual {
		if a != expected[i] {
			t.Errorf("wrong conversion for interval %v: expected \"%v\" got \"%v\"\n", i, expected[i], a)
		}
	}

}
