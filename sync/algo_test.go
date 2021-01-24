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
	"testing"
	"time"
)

func contains(slice []data.Interval, interval data.Interval) bool {
	keySlice := storage.ConvertToKeys(slice)
	for _, a := range keySlice {
		if a == storage.IntervalToKey(interval) {
			return true
		}
	}
	return false
}

func TestSync(t *testing.T) {
	store := storage.Ephemeral{}
	serverState := []data.Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "a",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "b",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "c",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "x",
		},
	}
	added := []data.Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "a",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "b",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "d",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "e",
		},
	}
	removed := []data.Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "c",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "e",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "f",
		},
	}
	expected := []data.Interval{
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "a",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "b",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "d",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "e",
		},
		{
			Start:      time.Time{},
			End:        time.Time{},
			Tags:       nil,
			Annotation: "x",
		},
	}

	req := data.SyncRequest{
		UserID:  0,
		Added:   added,
		Removed: removed,
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverState)
	result, _, err := Sync(req, &store)

	if err != nil {
		t.Errorf("Sync failed with error %v", err)
	}

	if len(result) != len(expected) {
		t.Errorf("Sync result wrong. Expected %v got %v", expected, result)
	}
	for _, interval := range expected {
		if !contains(result, interval) {
			t.Errorf("Sync result does not contain interval %v", interval)
		}
	}
}
