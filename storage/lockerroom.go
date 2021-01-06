package storage

import (
	"sync"
)

// A LockerRoom is a collection of Mutexes mapped to user ids
type LockerRoom struct {
	globalLock sync.Mutex
	locks map[UserId]*sync.Mutex
}

// Sets up this LockerRoom instance
func (lr *LockerRoom) InitializeLockerRoom() {
	lr.locks = make(map[UserId]*sync.Mutex)
}

// Creates an entry into the locks map if the user does not exist yet
func (lr *LockerRoom) createUserIfNotExists(userId UserId) {
	lr.globalLock.Lock()
	defer lr.globalLock.Unlock()

	if lr.locks[userId] == nil {
		lr.locks[userId] = &sync.Mutex{}
	}
}

// Acquire the lock for this user id
func (lr *LockerRoom) Lock(userId UserId) {
	lr.createUserIfNotExists(userId)

	lr.locks[userId].Lock()
}

// Release the lock for this user id
func (lr *LockerRoom) Unlock(userId UserId) {
	lr.locks[userId].Unlock()
}
