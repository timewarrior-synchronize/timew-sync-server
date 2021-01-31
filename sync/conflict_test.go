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
	"reflect"
	"testing"
	"time"
)

func elementwiseEqual(aSlice []data.Interval, bSlice []data.Interval) bool {
	keyA := storage.ConvertToKeys(aSlice)
	keyB := storage.ConvertToKeys(bSlice)
	if len(keyA) != len(keyB) {
		return false
	}
	for _, a := range keyA {
		match := false
		for i, b := range keyB {
			if a == b {
				match = true
				keyB = append(keyB[:i], keyB[i+1:]...)
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}
func sliceString(s []data.Interval) {
	print("[\n")
	for _, i := range s {
		fmt.Printf("\n-------------------------------------\nStart = %v\nEnd = %v\nTags = %v\nAnnotation = %v", i.Start, i.End, i.Tags, i.Annotation)
	}
	fmt.Printf("\n]\n\n\n")
}

func TestSolveConflict_MultiConflict(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateMultiConflict := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"starting"},
			Annotation: "problemStart",
		},
		{
			Start:      time.Date(2000, 5, 11, 13, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 19, 30, 40, 0, time.UTC),
			Tags:       []string{"middle"},
			Annotation: "problemMiddle",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"ending"},
			Annotation: "problemEnd",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	multiConflictExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 13, 30, 40, 0, time.UTC), // changed end time
			Tags:       []string{"starting"},
			Annotation: "problemStart",
		},
		{
			Start:      time.Date(2000, 5, 11, 13, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			Tags:       []string{"middle", "problemMiddle", "problemStart", "starting"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"ending", "middle", "problemMiddle", "problemStart", "starting"},
			Annotation: "problemEnd",
		},
		{
			Start:      time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 19, 30, 40, 0, time.UTC),
			Tags:       []string{"ending", "middle", "problemEnd", "problemMiddle"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 11, 19, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"ending"},
			Annotation: "problemEnd",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateMultiConflict)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("MultiConflict: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("MultiConflict: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	sliceString(multiConflictExpected)
	sliceString(result)
	if !elementwiseEqual(multiConflictExpected, result) {
		t.Errorf("MultiConflict: State after solve wrong. Expected %v got %v", multiConflictExpected, result)
	}
}

func TestSolveConflict_InnerInterval(t *testing.T) {
	store := storage.Ephemeral{}
	serverInnerInterval := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"outer"},
			Annotation: "o",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"inner"},
			Annotation: "i",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	innerIntervalExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC), // changed end time
			Tags:       []string{"outer"},
			Annotation: "o",
		},
		{ // merged
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"i", "inner", "o", "outer"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC), // changed start time
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"outer"},
			Annotation: "o",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverInnerInterval)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("InnerInterval: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("InnerInterval: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(innerIntervalExpected, result) {
		t.Errorf("InnerInterval: State after solve wrong. Expected %v got %v", innerIntervalExpected, result)
	}
}

func TestSolveConflict_SameEnd(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateSameEnd := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"c", "a", "b"},
			Annotation: "problemOne",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"e", "d", "c"},
			Annotation: "problemTwo",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	sameEndExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			Tags:       []string{"c", "a", "b"},
			Annotation: "problemOne",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"a", "b", "c", "d", "e", "problemOne", "problemTwo"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateSameEnd)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("SameEnd: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("SameEnd: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(sameEndExpected, result) {
		t.Errorf("SameEnd: State after solve wrong. Expected %v got %v", sameEndExpected, result)
	}
}

func TestSolveConflict_SameStart(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateSameStart := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"c", "a", "b"},
			Annotation: "problemOne",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"e", "d", "c"},
			Annotation: "problemTwo",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	sameStartExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"a", "b", "c", "d", "e", "problemOne", "problemTwo"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"e", "d", "c"},
			Annotation: "problemTwo",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateSameStart)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("SameStart: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("SameStart: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(sameStartExpected, result) {
		t.Errorf("SameStart: State after solve wrong. Expected %v got %v", sameStartExpected, result)
	}
}

func TestSolveConflict_Congruent(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateCongruent := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"c", "a", "b"},
			Annotation: "problemOne",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"e", "d", "c"},
			Annotation: "problemTwo",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	congruentExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"a", "b", "c", "d", "e", "problemOne", "problemTwo"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateCongruent)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("Congruent: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("Congruent: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(congruentExpected, result) {
		t.Errorf("Congruent: State after solve wrong. Expected %v got %v", congruentExpected, result)
	}
}

func TestSolveConflict_Overlap(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateOverlap := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"starting"},
			Annotation: "problemStart",
		},
		{
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"ending"},
			Annotation: "problemEnd",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}

	overlapExpected := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "all normal here",
		},
		{
			Start:      time.Date(2000, 5, 11, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC), // changed end time
			Tags:       []string{"starting"},
			Annotation: "problemStart",
		},
		{ // merged
			Start:      time.Date(2000, 5, 11, 14, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"ending", "problemEnd", "problemStart", "starting"},
			Annotation: "",
		},
		{
			Start:      time.Date(2000, 5, 11, 18, 30, 40, 0, time.UTC), // changed start time
			End:        time.Date(2000, 5, 11, 21, 30, 40, 0, time.UTC),
			Tags:       []string{"ending"},
			Annotation: "problemEnd",
		},
		{
			Start:      time.Date(2000, 5, 12, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 12, 18, 30, 40, 0, time.UTC),
			Tags:       []string{"tag2"},
			Annotation: "all normal here",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateOverlap)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("Overlap: Solve failed with error %v", err)
	}
	if !conflict {
		t.Errorf("Overlap: Solve did not detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(overlapExpected, result) {
		t.Errorf("Overlap: State after solve wrong. Expected %v got %v", overlapExpected, result)
	}
}

func TestSolveConflict_NoConflicts(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateNoConflicts := []data.Interval{
		{
			Start:      time.Date(2000, 5, 10, 12, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 5, 10, 13, 30, 40, 0, time.UTC), // +1h
			Tags:       []string{"tag1", "tag2"},
			Annotation: "a",
		},
		{
			Start:      time.Date(2000, 4, 10, 13, 30, 40, 0, time.UTC),
			End:        time.Date(2000, 4, 10, 13, 30, 40, 0, time.UTC),
			Tags:       []string{"tag1"},
			Annotation: "b",
		},
	}
	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateNoConflicts)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("NoConflicts: Solve failed with error %v", err)
	}
	if conflict {
		t.Errorf("NoConflicts: Solve falsely detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(serverStateNoConflicts, result) {
		t.Errorf("NoConflicts: State after solve wrong. Expected %v got %v", serverStateNoConflicts, result)
	}
}

func TestSolveConflict_NoIntervals(t *testing.T) {
	store := storage.Ephemeral{}
	serverStateNoIntervals := []data.Interval{}

	store.Initialize()
	store.SetIntervals(storage.UserId(0), serverStateNoIntervals)

	conflict, err := SolveConflict(0, &store)
	if err != nil {
		t.Errorf("NoIntervals: Solve failed with error %v", err)
	}
	if conflict {
		t.Errorf("NoIntervals: Solve falsely detected a conflict")
	}
	result, _ := store.GetIntervals(storage.UserId(0))
	if !elementwiseEqual(serverStateNoIntervals, result) {
		t.Errorf("NoIntervals: State after solve wrong. Expected %v got %v", serverStateNoIntervals, result)
	}
}

func TestUniteTagsAndAnnotation_AllEmpty(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{},
		Annotation: "",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{},
		Annotation: "",
	}
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, []string{}) {
		t.Errorf("AllEmpty: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != "" {
		t.Errorf("AllEmpty: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_DifferentAnnotations(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_a"},
		Annotation: "a",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_b"},
		Annotation: "b",
	}
	tagsExpected := []string{"a", "b", "tag_a", "tag_b"}
	annotationExpected := ""
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("DifferentAnnotations: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("DifferentAnnotations: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_AnnotationAPresent(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_a"},
		Annotation: "a",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_b"},
		Annotation: "",
	}
	tagsExpected := []string{"tag_a", "tag_b"}
	annotationExpected := "a"
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("AnnotationAPresent: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("AnnotationAPresent: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_AnnotationBPresent(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_a"},
		Annotation: "",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_b"},
		Annotation: "b",
	}
	tagsExpected := []string{"tag_a", "tag_b"}
	annotationExpected := "b"
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("AnnotationBPresent: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("AnnotationBPresent: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_SameAnnotation(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_a"},
		Annotation: "same",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_b"},
		Annotation: "same",
	}
	tagsExpected := []string{"tag_a", "tag_b"}
	annotationExpected := "same"
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("SameAnnotation: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("SameAnnotation: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_TagOverlap(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_a", "tag_same"},
		Annotation: "a",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"tag_b", "tag_same"},
		Annotation: "b",
	}
	tagsExpected := []string{"a", "b", "tag_a", "tag_b", "tag_same"}
	annotationExpected := ""
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("TagOverlap: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("TagOverlap: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_NoTagsDifferentAnnotation(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{},
		Annotation: "a",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{},
		Annotation: "b",
	}
	tagsExpected := []string{"a", "b"}
	annotationExpected := ""
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("NoTagsDifferentAnnotation: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("NoTagsDifferentAnnotation: Annotation does not match. Expected %v got %v", "", annotation)
	}
}

func TestUniteTagsAndAnnotation_DifferentAnnotationsPresentInTags(t *testing.T) {
	a := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"b"},
		Annotation: "a",
	}
	b := data.Interval{
		Start:      time.Time{},
		End:        time.Time{},
		Tags:       []string{"x", "y", "z"},
		Annotation: "b",
	}
	tagsExpected := []string{"a", "b", "x", "y", "z"}
	annotationExpected := ""
	tags, annotation := UniteTagsAndAnnotation(a, b)
	if !reflect.DeepEqual(tags, tagsExpected) {
		t.Errorf("NoTagsDifferentAnnotation: Tags do not match. Expected %v got %v", []string{}, tags)
	}
	if annotation != annotationExpected {
		t.Errorf("NoTagsDifferentAnnotation: Annotation does not match. Expected %v got %v", "", annotation)
	}
}
