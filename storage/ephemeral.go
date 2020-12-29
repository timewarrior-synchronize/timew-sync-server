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
	"errors"
	"log"
)

// Ephemeral represents storage of user interval data.
// It contains the time intervals.
// Each interval is represented as a string in intervals.
// Data is not stored persistently.
type Ephemeral struct {
	intervals map[storageClient][]IntervalWithMetadata
}

// storageClient represents a client, not only with it's non unique ClientId,
// but also with it's UserId. This makes storageClients uniquely identifiable.
type storageClient struct {
	userId   UserId
	clientId ClientId
}

func (ep *Ephemeral) GetIntervals(userId UserId, clientId ClientId) []IntervalWithMetadata {
	client := storageClient{
		userId:   userId,
		clientId: clientId,
	}

	return ep.intervals[client]
}

func (ep *Ephemeral) SetIntervals(userId UserId, clientId ClientId, intervals []IntervalWithMetadata) {
	if ep.intervals == nil {
		ep.intervals = make(map[storageClient][]IntervalWithMetadata)
	}

	client := storageClient{
		userId:   userId,
		clientId: clientId,
	}

	ep.intervals[client] = intervals
	log.Printf("ephemeral: Set Intervals of User, Client %v, %v\n", userId, client)
}

func (ep *Ephemeral) AddInterval(userId UserId, clientId ClientId, interval IntervalWithMetadata) {
	client := storageClient{
		userId:   userId,
		clientId: clientId,
	}

	if ep.intervals == nil {
		ep.intervals = make(map[storageClient][]IntervalWithMetadata)
	}

	ep.intervals[client] = append(ep.intervals[client], interval)
	log.Printf("ephemeral: Added an Interval to User, Client %v, %v\n", userId, clientId)
}

func (ep *Ephemeral) RemoveInterval(userId UserId, clientId ClientId, interval *IntervalWithMetadata) {
	client := storageClient{
		userId:   userId,
		clientId: clientId,
	}

	intervalIndex, err := findInterval(interval, ep.intervals[client])
	if err != nil {
		log.Printf("ephemeral: Couldn't find Interval to remove. Skipping")
		return
	}

	ep.intervals[client] = append(ep.intervals[client][:intervalIndex], ep.intervals[client][intervalIndex+1:]...)
	log.Printf("ephemeral: Removed an Interval of User, Client %v, %v\n", userId, clientId)
}

func findInterval(wanted *IntervalWithMetadata, intervals []IntervalWithMetadata) (int, error) {
	for i, interval := range intervals {
		if &interval == wanted {
			return i, nil
		}
	}

	return 0, errors.New("ephemeral: Couldn't find interval")
}
