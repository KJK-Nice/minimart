What is Mutex?
That is an excellent and very important question, especially in Go.


  A Mutex (which stands for Mutual Exclusion) is a synchronization tool used to protect shared data from
   being accessed and modified by multiple goroutines (concurrent threads of execution) at the same
  time.


  Imagine you have a simple counter that is being incremented by two different goroutines
  simultaneously.


   1 Shared Counter = 0



  Goroutine 1:
   1. Reads the value of Counter (gets 0).
   2. Increments the value (0 + 1 = 1).
   3. Writes the new value (1) back to Counter.


  Goroutine 2:
   1. Reads the value of Counter (gets 0).
   2. Increments the value (0 + 1 = 1).
   3. Writes the new value (1) back to Counter.


  If these two goroutines run at almost the exact same time, they might both read the value 0 before
  either has a chance to write its new value back. The final result would be 1, when it should be 2.
  This is called a race condition, and it's a classic and dangerous bug in concurrent programming.


  How a Mutex Solves This


  A mutex acts like a lock or a "talking stick." Before a goroutine can access the shared data, it must
  first acquire the lock. Only one goroutine can hold the lock at a time. All other goroutines that try
  to acquire the lock will be forced to wait until the current holder releases it.

  Here's the same scenario with a mutex:


  Goroutine 1:
   1. Acquires the lock. (Goroutine 2 must now wait).
   2. Reads the value of Counter (gets 0).
   3. Increments the value (1).
   4. Writes the value back (Counter is now 1).
   5. Releases the lock.


  Goroutine 2:
   1. (Was waiting) Acquires the lock.
   2. Reads the value of Counter (gets 1).
   3. Increments the value (2).
   4. Writes the value back (Counter is now 2).
   5. Releases the lock.

  The final result is 2, which is correct. The mutex ensures that the read-modify-write operation is
  atomic (it happens as a single, indivisible unit).


  sync.Mutex vs. sync.RWMutex

  Go's sync package provides two main types of mutexes:


   1. `sync.Mutex`:
       * Has two methods: Lock() and Unlock().
       * It's a simple, exclusive lock. It doesn't care if you're reading or writing; only one goroutine
         can have the lock at a time for any reason.


   2. `sync.RWMutex` (Read/Write Mutex):
       * This is the one I used in the InMemoryEventBus. It's a more sophisticated lock that can improve
         performance in situations where you have many more reads than writes.
       * It has two pairs of methods:
           * Lock() and Unlock() for writing.
           * RLock() and RUnlock() for reading.
       * The Rule:
           * Any number of goroutines can hold a read lock (`RLock`) at the same time.
           * Only one goroutine can hold a write lock (`Lock`), and no other goroutines (neither readers
             nor writers) can acquire any lock until it's released.


  In our InMemoryEventBus, subscribing is a "write" operation (it modifies the subscribers map), so we
  use b.mu.Lock(). Publishing is a "read" operation (it only reads from the map), so we use
  b.mu.RLock(). This allows multiple events to be published concurrently, which is more efficient.
