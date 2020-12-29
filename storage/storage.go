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

import "time"

// A UserId represents a unique ID assigned to each user of the
// timewarrior sync server
type UserId int

// An Interval represents a time from Start to End.
// The Tags field represents the interval's tags as string.
type Interval struct {
	Start      time.Time
	End        time.Time
	Tags       string
	Annotation string
}

// Storage defines an interface for accessing stored intervals.
// Every User has a set of intervals, which can be accessed and modified independently.
type Storage interface {
	// GetIntervals returns all intervals associated with a user
	GetIntervals(userId UserId) ([]Interval, error)

	// SetIntervals overrides all intervals of a user
	SetIntervals(userId UserId, intervals []Interval) error

	// AddInterval adds an interval to a user's intervals
	AddInterval(userId UserId, interval Interval) error

	// RemoveInterval removes an interval from a user's intervals
	RemoveInterval(userId UserId, interval *Interval) error
}

var GlobalStorage Storage
