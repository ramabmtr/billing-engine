package lib

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLockManager(t *testing.T) {
	lockManager := NewLockManager()
	assert.NotNil(t, lockManager, "LockManager should not be nil")
	assert.Implements(t, (*LockManager)(nil), lockManager, "Should implement LockManager interface")
}

func TestGetLock(t *testing.T) {
	lockManager := NewLockManager()

	// Test getting a lock for a key
	lock1 := lockManager.GetLock("key1")
	assert.NotNil(t, lock1, "Lock should not be nil")

	// Test getting the same lock for the same key
	lock1Again := lockManager.GetLock("key1")
	assert.Same(t, lock1, lock1Again, "Should return the same lock instance for the same key")

	// Test getting a different lock for a different key
	lock2 := lockManager.GetLock("key2")
	assert.NotNil(t, lock2, "Lock should not be nil")
	assert.NotSame(t, lock1, lock2, "Should return different lock instances for different keys")
}

func TestLockManagerConcurrency(t *testing.T) {
	lockManager := NewLockManager()
	const numGoroutines = 10
	const numIterations = 100

	// Create a counter that will be incremented by multiple goroutines
	counter := 0
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines that will try to increment the counter
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Get the lock for the counter
				lock := lockManager.GetLock("counter")
				lock.Lock()
				// Critical section
				counter++
				lock.Unlock()
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check that the counter has been incremented correctly
	assert.Equal(t, numGoroutines*numIterations, counter, 
		"Counter should be incremented correctly with lock protection")
}

func TestLockManagerMultipleKeys(t *testing.T) {
	lockManager := NewLockManager()
	const numKeys = 5
	const numGoroutines = 10
	const numIterations = 100

	// Create counters for each key
	counters := make([]int, numKeys)
	var wg sync.WaitGroup
	wg.Add(numGoroutines * numKeys)

	// Launch goroutines that will increment counters for different keys
	for k := 0; k < numKeys; k++ {
		key := "key" + string(rune('0'+k))
		for i := 0; i < numGoroutines; i++ {
			go func(keyIndex int, keyName string) {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					lock := lockManager.GetLock(keyName)
					lock.Lock()
					counters[keyIndex]++
					lock.Unlock()
				}
			}(k, key)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check that each counter has been incremented correctly
	for k := 0; k < numKeys; k++ {
		assert.Equal(t, numGoroutines*numIterations, counters[k], 
			"Counter for key%d should be incremented correctly", k)
	}
}