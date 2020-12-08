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
	"testing"
)

func TestEphemeralStorage(t *testing.T) {
	var s Storage
	s = &EphemeralStorage{}

	intervals := []string{
		"inc 20201202T080000Z - 20201202T10000Z",
		"inc 20201202T110000Z - 20201202T12000Z",
	}

	s.overwriteIntervals(intervals)
	result := s.getIntervals()

	if len(result) != len(intervals) {
		t.Errorf("length doesn't match, expected %v, got %v", len(intervals), len(result))
	}

	for i, x := range result {
		if x != intervals[i] {
			t.Errorf("interval data does not match, expected %v, got %v", intervals[i], x)
		}
	}
}
