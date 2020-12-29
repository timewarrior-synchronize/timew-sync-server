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
	"time"
)

// Sync completely overrides the Storage with the new data and returns all stored intervals afterwards.
// This is a naive approach for testing and not the final sync algorithm.
func Sync(syncRequest data.SyncRequest) []data.Interval {
	intervalsWithMetadata := make([]storage.IntervalWithMetadata, len(syncRequest.Intervals))

	for i, interval := range syncRequest.Intervals {
		intervalsWithMetadata[i] = storage.IntervalWithMetadata{
			Interval:     interval,
			LastModified: time.Now(),
			Deleted:      false,
		}
	}

	storage.GlobalStorage.SetIntervals(0, 0, intervalsWithMetadata)
	intervalsWithMetadata = storage.GlobalStorage.GetIntervals(0, 0)

	syncedIntervals := make([]data.Interval, len(intervalsWithMetadata))

	for i, intervalWithMetadata := range intervalsWithMetadata {
		syncedIntervals[i] = intervalWithMetadata.Interval
	}

	return syncedIntervals
}
