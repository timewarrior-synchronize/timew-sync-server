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
package storage

import (
	"github.com/timewarrior-synchronize/timew-sync-server/data"
	"reflect"
	"testing"
	"time"
)

func TestIntervalToKey(t *testing.T) {
	start := time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC)
	annotation := "test"
	testInput := data.Interval{
		Start:      start,
		End:        end,
		Tags:       []string{"prank", "laugh"},
		Annotation: annotation,
	}

	expected := IntervalKey{
		Start:      start,
		End:        end,
		Tags:       `["prank","laugh"]`,
		Annotation: annotation,
	}

	result := IntervalToKey(testInput)
	if !result.Start.Equal(expected.Start) {
		t.Errorf("Expected Start time to be %v got %v", expected.Start, result.Start)
	}
	if !result.End.Equal(expected.End) {
		t.Errorf("Expected End time to be %v got %v", expected.End, result.End)
	}
	if result.Annotation != expected.Annotation {
		t.Errorf("Expected Annotation to be %v got %v", expected.Annotation, result.Annotation)
	}
	if result.Tags != expected.Tags {
		t.Errorf("Expected Tags tp be %v got %v", expected.Tags, result.Tags)
	}
}

func TestKeyToInterval(t *testing.T) {
	start := time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(2020, time.April, 1, 15, 30, 0, 0, time.UTC)
	annotation := "test"
	testInput := IntervalKey{
		Start:      start,
		End:        end,
		Tags:       `["prank","laugh"]`,
		Annotation: annotation,
	}
	expected := data.Interval{
		Start:      start,
		End:        end,
		Tags:       []string{"prank", "laugh"},
		Annotation: annotation,
	}
	result := KeyToInterval(testInput)
	if !result.Start.Equal(expected.Start) {
		t.Errorf("Expected Start time to be %v got %v", expected.Start, result.Start)
	}
	if !result.End.Equal(expected.End) {
		t.Errorf("Expected End time to be %v got %v", expected.End, result.End)
	}
	if result.Annotation != expected.Annotation {
		t.Errorf("Expected Annotation to be %v got %v", expected.Annotation, result.Annotation)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected Tags tp be %v got %v", expected.Tags, result.Tags)
	}
}

func TestConvertToIntervals(t *testing.T) {
	// test empty slice
	if !reflect.DeepEqual(ConvertToIntervals([]IntervalKey{}), []data.Interval{}) {
		t.Errorf("Empty slice does not map to emtpy slice")
	}
}

func TestConvertToKeys(t *testing.T) {
	// test empty slice
	if !reflect.DeepEqual(ConvertToKeys([]data.Interval{}), []IntervalKey{}) {
		t.Errorf("Empty slice does not map to emtpy slice")
	}
}
