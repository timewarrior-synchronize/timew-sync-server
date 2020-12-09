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
	"reflect"
	"testing"
	"time"
)

func TestStringsToIntervals(t *testing.T) {
	testData := []string{
		"inc 20201125T093910Z - 20201125T093943Z",
		"inc 20201125T095240Z - 20201125T095253Z # test",
		"inc 20201209T140521Z - 20201209T140533Z # a b c",
	}
	loc, _ := time.LoadLocation("UTC")
	expected := make([]storage.Interval, 3, 3)
	expected[0] = storage.Interval{
		Start:        time.Date(2020, 11, 25, 9, 39, 10, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 39, 43, 0, loc),
		Tags:         []string{},
		LastModified: time.Time{},
		Deleted:      false,
	}
	expected[1] = storage.Interval{
		Start:        time.Date(2020, 11, 25, 9, 52, 40, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 52, 53, 0, loc),
		Tags:         []string{"test"},
		LastModified: time.Time{},
		Deleted:      false,
	}
	expected[2] = storage.Interval{
		Start:        time.Date(2020, 12, 9, 14, 5, 21, 0, loc),
		End:          time.Date(2020, 12, 9, 14, 5, 33, 0, loc),
		Tags:         []string{"a", "b", "c"},
		LastModified: time.Time{},
		Deleted:      false,
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
		if !reflect.DeepEqual(actualInterval.Tags, expected[i].Tags) {
			t.Errorf("wrong tags for interval %v: expected %v of type %v got %v of type %v", i, expected[i].Tags, reflect.TypeOf(expected[i].Tags), actualInterval.Tags, reflect.TypeOf(actualInterval.Tags))
		}
		if actualInterval.Deleted {
			t.Errorf("wrong value of deleted flag for interval %v: expected false got %v", i, actualInterval.Deleted)
		}
	}
}

func TestIntervalsToStrings(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	testData := make([]storage.Interval, 3, 3)
	testData[0] = storage.Interval{
		Start:        time.Date(2020, 11, 25, 9, 39, 10, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 39, 43, 0, loc),
		Tags:         []string{},
		LastModified: time.Time{},
		Deleted:      false,
	}
	testData[1] = storage.Interval{
		Start:        time.Date(2020, 11, 25, 9, 52, 40, 0, loc),
		End:          time.Date(2020, 11, 25, 9, 52, 53, 0, loc),
		Tags:         []string{"test"},
		LastModified: time.Time{},
		Deleted:      false,
	}
	testData[2] = storage.Interval{
		Start:        time.Date(2020, 12, 9, 14, 5, 21, 0, loc),
		End:          time.Date(2020, 12, 9, 14, 5, 33, 0, loc),
		Tags:         []string{"a", "b", "c"},
		LastModified: time.Time{},
		Deleted:      false,
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
			t.Errorf("wrong conversion for interval %v: expected \"%v\" got \"%v\"\n", i, a, expected[i])
		}
	}

}
