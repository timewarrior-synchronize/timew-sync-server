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
	"log"
)

// Ephemeral represents storage of user interval data.
// It contains the time intervals.
// Each interval is represented as a string in intervals.
// Data is not stored persistently.
type Ephemeral struct {
	LockerRoom
	intervals map[UserId]intervalSet
}

// intervalSet represents a set of intervals
type intervalSet map[IntervalKey]bool

// Initialize runs all necessary setup for this Storage instance
func (ep *Ephemeral) Initialize() error {
	ep.intervals = make(map[UserId]intervalSet)
	ep.InitializeLockerRoom()
	return nil
}

// GetIntervals returns all intervals stored for a specific user
func (ep *Ephemeral) GetIntervals(userId UserId) ([]data.Interval, error) {
	intervals := make([]IntervalKey, len(ep.intervals[userId]))

	i := 0
	for interval := range ep.intervals[userId] {
		intervals[i] = interval
		i++
	}

	return ConvertToIntervals(intervals), nil
}

// SetIntervals replaces all intervals of a specific user
func (ep *Ephemeral) SetIntervals(userId UserId, intervals []data.Interval) error {
	keys := ConvertToKeys(intervals)
	ep.intervals[userId] = make(intervalSet, len(keys))
	for _, key := range keys {
		ep.intervals[userId][key] = true
	}
	log.Printf("ephemeral: Set Intervals of User %v\n", userId)

	return nil
}

// AddInterval adds a single interval to the intervals stored for a user
func (ep *Ephemeral) AddInterval(userId UserId, interval data.Interval) error {
	if ep.intervals[userId] == nil {
		ep.intervals[userId] = make(intervalSet)
	}

	ep.intervals[userId][IntervalToKey(interval)] = true
	log.Printf("ephemeral: Added an Interval to User %v\n", userId)

	return nil
}

// RemoveInterval removes an interval from the intervals stored for a user
func (ep *Ephemeral) RemoveInterval(userId UserId, interval data.Interval) error {
	delete(ep.intervals[userId], IntervalToKey(interval))
	log.Printf("ephemeral: Removed an Interval of User %v\n", userId)

	return nil
}

// ModifyIntervals atomically adds and deletes a specified set
// of intervals
func (ep *Ephemeral) ModifyIntervals(userId UserId, add []data.Interval, del []data.Interval) error {
	for _, interval := range del {
		delete(ep.intervals[userId], IntervalToKey(interval))
	}

	if ep.intervals[userId] == nil {
		ep.intervals[userId] = make(intervalSet, len(add))
	}
	for _, interval := range add {
		ep.intervals[userId][IntervalToKey(interval)] = true
	}

	log.Printf("ephemeral: Modified Intervals of User %v\n", userId)

	return nil
}
