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
	"github.com/google/go-cmp/cmp"
	"github.com/timewarrior-synchronize/timew-sync-server/data"
	"testing"
	"time"
)

func TestEphemeralStorage(t *testing.T) {
	var s Storage
	intervals := []data.Interval{
		{
			Start:      time.Date(2020, time.December, 24, 18, 0, 0, 0, time.UTC),
			End:        time.Date(2020, time.December, 24, 22, 0, 0, 0, time.UTC),
			Tags:       []string{"Merry", "Christmas"},
			Annotation: "test",
		},
		{},
	}
	s = &Ephemeral{}
	_ = s.Initialize()
	_ = s.SetIntervals(0, intervals)
	result, _ := s.GetIntervals(0)

	if len(result) != len(intervals) {
		t.Errorf("length doesn't match, expected %v, got %v", len(intervals), len(result))
	}

	for _, x := range result {
		correct := false
		for i, _ := range intervals {
			if diff := cmp.Diff(intervals[i], x); diff == "" {
				correct = true
			}
		}
		if !correct {
			t.Errorf("result: %v not as expected: %v They do not contain exactly the same elements", result, intervals)
		}
	}
}

func TestEphemeralStorage_ModifyIntervals(t *testing.T) {
	var s Storage
	add := []data.Interval{
		{
			Start:      time.Date(2020, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag3", "Tag4"},
			Annotation: "Annotation2",
		},
	}
	del := []data.Interval{
		{
			Start:      time.Date(2021, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag1", "Tag2"},
			Annotation: "Annotation",
		},
	}

	s = &Ephemeral{}
	_ = s.Initialize()
	_ = s.SetIntervals(42, del)
	_ = s.ModifyIntervals(42, add, del)
	result, _ := s.GetIntervals(42)

	if len(result) != len(add) {
		t.Errorf("length doesn't match, expected %v, got %v", len(add), len(result))
	}

	for _, x := range result {
		correct := false
		for i, _ := range add {
			if diff := cmp.Diff(add[i], x); diff == "" {
				correct = true
			}
		}
		if !correct {
			t.Errorf("result: %v not as expected: %v They do not contain exactly the same elements", result, add)
		}
	}
}

func TestEphemeral_ModifyIntervals_add(t *testing.T) {
	var s Storage

	add := []data.Interval{
		{
			Start:      time.Date(2020, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag3", "Tag4"},
			Annotation: "Annotation2",
		},
		{
			Start:      time.Date(2021, 01, 01, 12, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 01, 01, 13, 0, 0, 0, time.UTC),
			Tags:       []string{"Tag1", "Tag2"},
			Annotation: "Annotation1",
		},
	}

	s = &Ephemeral{}
	_ = s.Initialize()
	_ = s.ModifyIntervals(0, add, []data.Interval{})
	result, _ := s.GetIntervals(0)

	if len(result) != len(add) {
		t.Errorf("length doesn't match, expected %v, got %v", len(add), len(result))
	}

	for _, x := range result {
		correct := false
		for i, _ := range add {
			if diff := cmp.Diff(add[i], x); diff == "" {
				correct = true
			}
		}
		if !correct {
			t.Errorf("result: %v not as expected: %v They do not contain exactly the same elements", result, add)
		}
	}
}

func TestEphemeral_AddInterval(t *testing.T) {
	var s Storage

	add := data.Interval{
		Start:      time.Date(2020, 01, 01, 12, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 01, 01, 13, 0, 0, 0, time.UTC),
		Tags:       []string{"Tag3", "Tag4"},
		Annotation: "Annotation2",
	}

	s = &Ephemeral{}
	_ = s.Initialize()
	_ = s.AddInterval(0, add)
	result, _ := s.GetIntervals(0)

	if len(result) != 1 {
		t.Errorf("length doesn't match, expected %v, got %v", 1, len(result))
	}

	if diff := cmp.Diff(add, result[0]); diff != "" {
		t.Errorf("result: %v not as expected: %v", result, add)
	}
}
