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
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
)

// A UserId represents a unique ID assigned to each user of the
// timewarrior sync server
type UserId int

// An Interval represents a time from Start to End.
// It also contains LastModified timestamp and Deleted flag needed for synchronization
// The Tags field represents the intervals tags as a slice of string. If there are no tags associated with this
// particular interval, tags should be a slice of length 0
type Interval data.Interval

// Storage defines an interface for accessing stored intervals.
// Every User has a set of intervals, which can be accessed and modified independently.
type Storage interface {
	// GetIntervals returns all intervals associated with a user
	GetIntervals(userId UserId) []Interval

	// SetIntervals overrides all intervals of a user
	SetIntervals(userId UserId, intervals []Interval)

	// AddInterval adds an interval to a user's intervals
	AddInterval(userId UserId, interval Interval)

	// RemoveInterval removes
	RemoveInterval(userId UserId, interval *Interval)
}

var GlobalStorage Storage
