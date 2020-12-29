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

package sync

import (
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/storage"
	"log"
	"strings"
)

// Sync completely overrides the Storage with the new data and returns all stored intervals afterwards.
// This is a naive approach for testing and not the final sync algorithm.
func Sync(syncRequest data.SyncRequest) []data.Interval {
	intervals := make([]storage.Interval, len(syncRequest.Intervals))
	for i, interval := range syncRequest.Intervals {
		tags := strings.Join(interval.Tags, " ")

		intervals[i] = storage.Interval{
			Start:      interval.Start,
			End:        interval.End,
			Tags:       tags,
			Annotation: "",
		}
	}

	err := storage.GlobalStorage.SetIntervals(0, intervals)
	if err != nil {
		panic("Error while writing to storage. Aborting sync process.")
	}

	intervals, err = storage.GlobalStorage.GetIntervals(0)
	if err != nil {
		log.Fatalf("Error while reading from storage: %v", err)
	}

	syncedIntervals := make([]data.Interval, len(intervals))

	for i, interval := range intervals {
		syncedIntervals[i] = data.Interval{
			Start: interval.Start,
			End:   interval.End,
			// TODO: Replace with something useful, when either data.Interval is modified
			//		or the Tag parser is rewritten
			//		- Vincent Stollenwerk
			Tags: []string{interval.Tags},
		}
	}

	return syncedIntervals
}
