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
)

// Sync updates the stored state in passed storage.Storage for the user issuing the sync request. If something fails it
//tries to restore the state
// prior to the syncRequest. This is not always possible though. The error message denotes whether restoring state was
//successful.
// Later atomicity should be guaranteed by storage.
// Iff no errors occur Sync returns the synced interval data of the user issuing the sync request.
func Sync(syncRequest data.SyncRequest, store storage.Storage) ([]data.Interval, bool, error) {
	// acquire lock and release it after syncing
	store.Lock(storage.UserId(syncRequest.UserID))
	defer store.Unlock(storage.UserId(syncRequest.UserID))

	// First, remove all intervals the client removed in its diff
	backup, err := store.GetIntervals(storage.UserId(syncRequest.UserID))
	if err != nil {
		return nil, false, fmt.Errorf("fatal error: Could not retrieve stored intervals for backup. " +
			"Stored state did not change")
	}
	for _, removedInterval := range syncRequest.Removed {
		err = store.RemoveInterval(storage.UserId(syncRequest.UserID), removedInterval)
		if err != nil {
			restoreError := store.SetIntervals(storage.UserId(syncRequest.UserID), backup) // trying to restore backup
			if restoreError != nil {
				return nil, false, fmt.Errorf("fatal error: Failed to remove interval %v from storage. "+
					"Also could not restore server state", removedInterval)
			} else {
				return nil, false, fmt.Errorf("fatal error: Failed to remove interval %v from storage. "+
					"Stored state did not change", removedInterval)
			}

		}
	}

	// Then add all intervals the client added in its diff
	for _, addedInterval := range syncRequest.Added {
		err = store.AddInterval(storage.UserId(syncRequest.UserID), addedInterval)
		if err != nil {
			restoreError := store.SetIntervals(storage.UserId(syncRequest.UserID), backup) // trying to restore backup
			if restoreError != nil {
				return nil, false, fmt.Errorf("fatal error: Failed to add interval %v to storage. "+
					"Also could not restore server state", addedInterval)
			} else {
				return nil, false, fmt.Errorf("fatal error: Failed to add interval %v to storage. "+
					"Stored state did not change", addedInterval)
			}

		}
	}

	conflict, solveErr := SolveConflict(syncRequest.UserID, store)
	if solveErr != nil {
		restoreError := store.SetIntervals(storage.UserId(syncRequest.UserID), backup) // try to restore backup
		if restoreError != nil {
			return nil, conflict, fmt.Errorf("fatal error: Failed to solve conflicts %v. "+
				"Also could not restore server state", solveErr)
		} else {
			return nil, conflict, fmt.Errorf("fatal error: Failed to solve conflicts %v. "+
				"Stored state unchanged", solveErr)
		}
	}

	result, err2 := store.GetIntervals(storage.UserId(syncRequest.UserID))
	if err2 != nil {
		restoreError := store.SetIntervals(storage.UserId(syncRequest.UserID), backup) // trying to restore backup
		if restoreError != nil {
			return nil, conflict, fmt.Errorf("fatal error: Failed to retrieve intervals from storage. " +
				"Also could not restore server state")
		} else {
			return nil, conflict, fmt.Errorf("fatal error: Failed to retrieve intervals from storage. " +
				"Stored state did not change")
		}
	}
	return result, conflict, nil
}
