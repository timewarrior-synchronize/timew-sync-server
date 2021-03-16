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
	"encoding/json"
	"github.com/timewarrior-synchronize/timew-sync-server/data"
	"log"
	"time"
)

// A UserId represents a unique ID assigned to each user of the
// timewarrior sync server
type UserId int

type IntervalKey struct {
	Start      time.Time
	End        time.Time
	Tags       string
	Annotation string
}

// ConvertToKeys converts a slice of data.Interval to a slice of IntervalKey
func ConvertToKeys(data []data.Interval) []IntervalKey {
	result := make([]IntervalKey, len(data), len(data))
	for i, interval := range data {
		result[i] = IntervalToKey(interval)
	}
	return result
}

// ConvertToIntervals converts a slice of IntervalKey to a slice of data.Interval
func ConvertToIntervals(keys []IntervalKey) []data.Interval {
	result := make([]data.Interval, len(keys), len(keys))
	for i, key := range keys {
		result[i] = KeyToInterval(key)
	}
	return result
}

// IntervalToKey converts a data.Interval struct to an IntervalKey struct which can be used as key in maps
func IntervalToKey(data data.Interval) IntervalKey {
	result, err := json.Marshal(data.Tags)
	if err != nil {
		log.Printf("Error parsing tag Array %v to json string", data.Tags)
	}
	return IntervalKey{
		Start:      data.Start,
		End:        data.End,
		Tags:       string(result),
		Annotation: data.Annotation,
	}
}

// KeyToInterval converts an IntervalKey struct to a data.Interval struct.
func KeyToInterval(key IntervalKey) data.Interval {
	result := data.Interval{
		Start:      key.Start,
		End:        key.End,
		Tags:       nil,
		Annotation: key.Annotation,
	}
	err := json.Unmarshal([]byte(key.Tags), &result.Tags)
	if err != nil {
		log.Printf("Error parsing Tags json-String %v to slice of string", key.Tags)
	}
	return result
}

// Storage defines an interface for accessing stored intervals.
// Every User has a set of intervals, which can be accessed and modified independently.
type Storage interface {
	// Initialize runs all necessary setup for this Storage instance
	Initialize() error

	// Acquire the lock for this user id
	Lock(userId UserId)

	// Release the lock for this user id
	Unlock(userId UserId)

	// GetIntervals returns all intervals associated with a user
	GetIntervals(userId UserId) ([]data.Interval, error)

	// SetIntervals overrides all intervals of a user
	SetIntervals(userId UserId, intervals []data.Interval) error

	// ModifyIntervals atomically adds and deletes a specified set
	// of intervals
	ModifyIntervals(userId UserId, add []data.Interval, del []data.Interval) error

	// AddInterval adds an interval to a user's intervals
	AddInterval(userId UserId, interval data.Interval) error

	// RemoveInterval removes an interval from a user's intervals
	RemoveInterval(userId UserId, interval data.Interval) error
}

var GlobalStorage Storage
