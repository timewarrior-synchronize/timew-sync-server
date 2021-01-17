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
	"fmt"
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/storage"
	"sort"
)

// SolveConflict merges overlapping intervals of given user.
// It then updates userId's state in store accordingly
// SolveConflict returns true iff a conflict was detected
func SolveConflict(userId int, store storage.Storage) (bool, error) {
	conflictDetected := false
	intervals, err := store.GetIntervals(storage.UserId(userId))
	var removed []data.Interval
	var added []data.Interval
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve User Data for UserId %v from Storage:\n%v", userId, err)
	}
	// Sort intervals by ascending start time (in place)
	sort.SliceStable(intervals, func(i, j int) bool {
		return intervals[i].Start.Before(intervals[j].Start)
	})
	fmt.Printf("\n UserId: %v   len: %v\n", userId, len(intervals))
	if len(intervals) == 0 {
		return false, nil
	}
	openInterval := intervals[0]
	var nextInterval data.Interval
	intervals = intervals[1:]
	var addedThisIteration []data.Interval
	// loop invariant:
	// openInterval.Start <= intervals[i].Start for all 0 <= i < len(intervals)
	// and intervals[i].Start <= intervals[i+1].Start for all 0 <= i < len(intervals) - 1
	// in short: append([]data.Interval{openInterval}, intervals) is always sorted by start time
	for len(intervals) > 0 {
		//fmt.Printf("Loop len: %v\n", len(intervals))
		fmt.Printf("\n UserId: %v   len: %v\n", userId, len(intervals))
		fmt.Printf("open: %v\n\nintervals: %v\n\n", openInterval, intervals)
		interval := intervals[0]
		intervals = intervals[1:]
		addedThisIteration = []data.Interval{}
		if interval.Start.Equal(openInterval.End) || interval.Start.After(openInterval.End) { // standard case - no overlapping intervals
			openInterval = interval
		} else {
			conflictDetected = true
			removed = append(removed, openInterval, interval)
			fmt.Print("New Intervals: \n")
			// end section (if exists)
			if !openInterval.End.Equal(interval.End) {
				if openInterval.End.After(interval.End) {

					nextInterval = data.Interval{
						Start:      interval.End,
						End:        openInterval.End,
						Tags:       openInterval.Tags,
						Annotation: openInterval.Annotation,
					}
				} else {
					nextInterval = data.Interval{
						Start:      openInterval.End,
						End:        interval.End,
						Tags:       interval.Tags,
						Annotation: interval.Annotation,
					}
				}
				fmt.Printf("End: %v\n", nextInterval)
				addedThisIteration = append(addedThisIteration, nextInterval)
			}

			// middle section
			tags, annotation := uniteTagsAndAnnotation(openInterval, interval)
			if openInterval.End.After(interval.End) {
				nextInterval = data.Interval{
					Start: interval.Start, // We have to use this start time since this is the middle section and
					// and interval.Start >= openInterval.Start by loop invariant
					End:        interval.End,
					Tags:       tags,
					Annotation: annotation,
				}

			} else {
				nextInterval = data.Interval{
					Start:      interval.Start,
					End:        openInterval.End,
					Tags:       tags,
					Annotation: annotation,
				}
			}
			fmt.Printf("Middle: %v\n", nextInterval)
			addedThisIteration = append(addedThisIteration, nextInterval)

			// start section
			if !openInterval.Start.Equal(interval.Start) {
				nextInterval = data.Interval{
					Start:      openInterval.Start,
					End:        interval.Start,
					Tags:       openInterval.Tags,
					Annotation: openInterval.Annotation,
				}
				addedThisIteration = append(addedThisIteration, nextInterval)
				fmt.Printf("Start: %v\n", nextInterval)
			}

			//getting ready for next iteration
			openInterval = nextInterval
			added = append(added, addedThisIteration...)
			// reinsert newly created intervals
			intervals = append(intervals, addedThisIteration[:len(addedThisIteration)-1]...)
			sort.SliceStable(intervals, func(i, j int) bool {
				return intervals[i].Start.Before(intervals[j].Start)
			}) // Maybe just iterating from left to right over intervals and inserting at the correct time is faster,
			// since intervals from addedThisIteration will probably have an "early" start time
		}

		fmt.Print(" -------------------------------------------\n\n")
	}
	store.Lock(storage.UserId(userId))
	defer store.Unlock(storage.UserId(userId))
	for _, a := range added {
		err = store.AddInterval(storage.UserId(userId), a)
		if err != nil {
			return conflictDetected, fmt.Errorf("Unable to change User Data for UserId %v in Storage:\n%v", userId, err)
		}
	}
	for _, r := range removed {
		err = store.RemoveInterval(storage.UserId(userId), r)
		if err != nil {
			return conflictDetected, fmt.Errorf("Unable to change User Data for UserId %v in Storage:\n%v", userId, err)
		}
	}

	return conflictDetected, nil
}

// uniteTagsAndAnnotation computes the new tags and annotation for overlapping intervals and returns tags, annotation.
// Case 1: Iff only one interval has an Annotation, we use this annotation. Case 2: Iff no interval has an annotation,
// we use "" as annotation. Case 3: Iff both intervals have different annotatio, we use "" as annotation, and add both
// annotation to tags. Case 4: Iff both intervals have the same annotation, we just use that annotation
// As tags we return the alphabetically sorted union of both intervals' tags (and both annotations in Case 3)
// without duplicates.
func uniteTagsAndAnnotation(a data.Interval, b data.Interval) ([]string, string) {
	tags := make([]string, len(a.Tags), len(a.Tags)+len(b.Tags))
	tmp := make([]string, len(b.Tags))
	copy(tags, a.Tags)
	copy(tmp, b.Tags)
	annotation := ""
	tags = append(tags, tmp...)
	if a.Annotation != "" && b.Annotation != "" && a.Annotation != b.Annotation {
		tags = append(tags, a.Annotation, b.Annotation)
	} else if a.Annotation == "" {
		annotation = b.Annotation
	} else {
		annotation = a.Annotation
	}
	sort.Strings(tags)
	i := 1
	for i < len(tags) {
		if tags[i] == tags[i-1] {
			tags = append(tags[:i], tags[i+1:]...)
		} else {
			i++
		}
	}
	return tags, annotation
}
