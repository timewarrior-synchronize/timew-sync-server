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
	"time"
)

// A UserId represents a unique ID assigned to each user of the
// timewarrior sync server
type UserId int

// A ClientId represents an ID assigned to each client of a user. The
// client IDs are not globally unique, instead they are only unique
// for a given user. A user always has at least one client.
type ClientId int

// An Interval represents a time from Start to End.
// It also contains LastModified timestamp and Deleted flag needed for synchronization
type Interval struct {
	Start time.Time
	End   time.Time

	LastModified time.Time
	Deleted      bool
}

// Storage defines an interface for accessing stored intervals.
type Storage interface {
	GetIntervals() []string
	OverwriteIntervals(intervals []string)
}

var GlobalStorage Storage