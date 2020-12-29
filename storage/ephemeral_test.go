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

package storage

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestEphemeralStorage(t *testing.T) {
	var s Storage
	s = &Ephemeral{}

	intervals := []Interval{
		{
			Start: time.Date(2020, time.December, 24, 18, 0, 0, 0, time.UTC),
			End:   time.Date(2020, time.December, 24, 22, 0, 0, 0, time.UTC),
			Tags:  "Christmas",
		},
		{},
	}

	_ = s.SetIntervals(0, intervals)
	result, _ := s.GetIntervals(0)

	if len(result) != len(intervals) {
		t.Errorf("length doesn't match, expected %v, got %v", len(intervals), len(result))
	}

	for i, x := range result {
		if diff := cmp.Diff(intervals[i], x); diff != "" {
			t.Errorf("interval data does not match, expected %v, got %v", intervals[i], x)
		}
	}
}
