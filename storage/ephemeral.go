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
	"log"
	"sync"
)

// Ephemeral represents storage of user interval data.
// It contains the time intervals.
// Each interval is represented as a string in intervals.
// Data is not stored persistently.
type Ephemeral struct {
	globalLock sync.Mutex
	locks map[UserId]*sync.Mutex
	intervals map[UserId]intervalSet
}

// intervalSet represents a set of intervals
type intervalSet map[Interval]bool

// Initialize runs all necessary setup for this Storage instance
func (ep *Ephemeral) Initialize() error {
	ep.intervals = make(map[UserId]intervalSet)
	ep.locks = make(map[UserId]*sync.Mutex)
	return nil
}

// Creates an entry into the locks map if the user does not exist yet
func (ep *Ephemeral) createUserIfNotExists(userId UserId) {
	ep.globalLock.Lock()
	defer ep.globalLock.Unlock()

	if ep.locks[userId] == nil {
		ep.locks[userId] = &sync.Mutex{}
	}
}

// GetIntervals returns all intervals stored for a specific user
func (ep *Ephemeral) GetIntervals(userId UserId) ([]Interval, error) {
	ep.createUserIfNotExists(userId)

	ep.locks[userId].Lock()
	defer ep.locks[userId].Unlock()

	intervals := make([]Interval, len(ep.intervals[userId]))

	i := 0
	for interval := range ep.intervals[userId] {
		intervals[i] = interval
		i++
	}

	return intervals, nil
}

// SetIntervals replaces all intervals of a specific user
func (ep *Ephemeral) SetIntervals(userId UserId, intervals []Interval) error {
	ep.createUserIfNotExists(userId)

	ep.locks[userId].Lock()
	defer ep.locks[userId].Unlock()

	ep.intervals[userId] = make(map[Interval]bool)
	for _, interval := range intervals {
		ep.intervals[userId][interval] = true
	}
	log.Printf("ephemeral: Set Intervals of User %v\n", userId)

	return nil
}

// AddInterval adds a single interval to the intervals stored for a user
func (ep *Ephemeral) AddInterval(userId UserId, interval Interval) error {
	ep.createUserIfNotExists(userId)

	ep.locks[userId].Lock()
	defer ep.locks[userId].Unlock()

	ep.intervals[userId][interval] = true
	log.Printf("ephemeral: Added an Interval to User %v\n", userId)

	return nil
}

// RemoveInterval removes an interval from the intervals stored for a user
func (ep *Ephemeral) RemoveInterval(userId UserId, interval Interval) error {
	ep.createUserIfNotExists(userId)

	ep.locks[userId].Lock()
	defer ep.locks[userId].Unlock()

	delete(ep.intervals[userId], interval)
	log.Printf("ephemeral: Removed an Interval of User %v\n", userId)

	return nil
}
