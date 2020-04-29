---
tags: ccns
---

# 讀書會 Week 7
:::info
Content: The Go Programming Language
         Chapter 9. Concurrency with Shared Variables
:::

## 9.1 Race Conditions
* Avoid concurrent access to most variables by...
    * Confining them to single goroutine
    * Maintaining a higher-level invariant of mutual exclusion
* Deposit example
    * Try [sync/atomic](https://golang.org/pkg/sync/atomic/) maybe?
* To avoid data race...
    * Try **not** to write the variable
    * Avoid accessing the variable **from multiple goroutines**
        > Then how to access variable?
        > Use **channel**
    * Allow multiple goroutines to access the variable, but one at a time
* *monitor goroutine*
    * A goroutine that brokers access to a confined variable using channel requests.
* *serial confinement*
    * Share a variable between goroutines in a pipeline by passing address from one stage to another over a channel
## 9.2 Mutual Exclusion: sync.Mutex
* *binary semaphore*
    * A semaphore counts to 1
* Mutex
  ```golang
  import "sync"
  var(
      mu    sync.Mutex
  )
  
  func foo(){
      mu.Lock()
      defer mu.Unlock()
      // DO SOMETHING
  }
  ```
    * **critical section**
        * The region of code between **Lock()** and **Unlock()**
    * It's important to unlock the mutex on all routes, including error route
* There can only be one mutex at a time, only one function can hold the lock
* If there's function that requires lock to be held, remember to create comment
## 9.3 Read/Write Mutexes: sync.RWMutex
* **sync.RWMutex**
    * Multiple readers, single writer lock
* **mu.RLock(), mu.RUnlock()**
    * A lock mechanism that allows multiple read-only ops to run in parallel
    * Can only be used when there's no write ops in critical section
    * When not sure, use exclusive **Lock()**
  ```golang
  var mu sync.RWMutex
  
  func boo(){
      mu.RLock() // Read-only
      defer mu.RUnlock()
      // DO SOMETHING
  }
  ```
* RWMutex is slower than regular Mutex
## 9.4 Memory Synchronization
* Just confine variables to a single goroutine
* Or use mutual exclusion
## 9.5 Lazy Initialization: sync.Once
* Use multiple steps and RWMutex to load variables
    * Sample in the book
* or use **sync.Once**
    * Combines a mutex and a boolean var
    * If the boolean is false, call function and change the boolean var
    * Then the function won't be called again
  ```golang
  var mu sync.Once
  func foo() {
      mu.Do(bar())
      return
  }
  ```
## 9.6 The Race Detector
* `-race` flag
    * Can be added to `build`, `run`, `test` param
* [Go Race Detector](https://golang.org/doc/articles/race_detector.html)
## 9.7 Example: Concurrent Non-Blocking Cache
* *memo1*
  ```golang
  func (memo *Memo) Get(key string) (interface{}, error){
      res, ok := memo.cache[key]
      if !ok {
          res.value, res.err = memo.f(key)
          memo.cache[key] = res
      }
      return res.value, res.err
  }
  ```
  ```lang=
  WARNING: DATA RACE
  Write at 0x00c00011cd80 by goroutine 14:
  runtime.mapassign_faststr()
      /usr/lib/go/src/runtime/map_faststr.go:202 +0x0
  gopl.io/ch9/memo1.(*Memo).Get()
      /home/tsundere/go/src/gopl.io/ch9/memo1/memo.go:35 +0x1ce
  gopl.io/ch9/memotest.Concurrent.func1()
      /home/tsundere/go/src/gopl.io/ch9/memotest/memotest.go:93 +0xde
  ```
* *memo2*
  ```golang
  func (memo *Memo) Get(key string) (interface{}, error){
      memo.mu.Lock()
      res, ok := memo.cache[key]
      if !ok {
          res.value, res.err = memo.f(key)
          memo.cache[key] = res
      }
      memo.mu.Unlock()
      return res.value, res.err
  }
  ```
* *memo3*
  ```golang
  func (memo *Memo) Get(key string) (interface{}, error){
      memo.mu.Lock()
      res, ok := memo.cache[key]
      memo.mu.Unlock()
      if !ok {
          res.value, res.err = memo.f(key)
          
          memo.mu.Lock()
          memo.cache[key] = res
          memo.mu.Unlock()
      }
      return res.value, res.err
  }
  ```
* *memo4*
  ```golang
  // Func is the type of the function to memoize.
  type Func func(string) (interface{}, error)

  type result struct {
      value interface{}
      err   error
  }

  //!+
  type entry struct {
      res   result
      ready chan struct{} // closed when res is ready
  }

  func New(f Func) *Memo {
      return &Memo{f: f, cache: make(map\[string\]*entry)}
  }

  type Memo struct {
      f     Func
      mu    sync.Mutex // guards cache
      cache map\[string\]*entry
  }

  func (memo *Memo) Get(key string) (value interface{}, err error) {
      memo.mu.Lock()
      e := memo.cache[key]
      if e == nil {
          // This is the first request for this key.
          // This goroutine becomes responsible for computing
          // the value and broadcasting the ready condition.
          e = &entry{ready: make(chan struct{})}
          memo.cache\[key\] = e
          memo.mu.Unlock()

          e.res.value, e.res.err = memo.f(key)

          close(e.ready) // broadcast ready condition
      } else {
          // This is a repeat request for this key.
          memo.mu.Unlock()

          <-e.ready // wait for ready condition
      }
      return e.res.value, e.res.err
  }
  ```
* *memo5*
  ```golang
  // A request is a message requesting that the Func be applied to key.
  type request struct {
      key      string
      response chan<- result // the client wants a single result
  }

  type Memo struct{ requests chan request }

  // New returns a memoization of f.  Clients must subsequently call Close.
  func New(f Func) *Memo {
      memo := &Memo{requests: make(chan request)}
      go memo.server(f)
      return memo
  }

  func (memo *Memo) Get(key string) (interface{}, error) {
      response := make(chan result)
      memo.requests <- request{key, response}
      res := <-response
      return res.value, res.err
  }

  func (memo *Memo) Close() { close(memo.requests) }
  
  func (memo *Memo) server(f Func) {
	cache := make(map\[string\]*entry)
	for req := range memo.requests {
		e := cache\[req.key\]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache\[req.key\] = e
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response)
	}
  }

  func (e *entry) call(f Func, key string) {
	  // Evaluate the function.
	  e.res.value, e.res.err = f(key)
	  // Broadcast the ready condition.
	  close(e.ready)
  }

  func (e *entry) deliver(response chan<- result) {
	  // Wait for the ready condition.
	  <-e.ready
	  // Send the result to the client.
	  response <- e.res
  }
  ```
## 9.8 Goroutines and Threads
* A goroutine starts with a small stack, like 2KB
* Then it grows and shrinks as needed
* **GOMAXPROCS**
    * Determine how many OS threads may be actively executing Go code simultaneously
  ```golang
  for {
    go fmt.Print(0)
    fmt.Print(1)
  }
  $ GOMAXPROCS=1 go run hacker-cliché.go
  111111111111111111110000000000000000000011111...
  $ GOMAXPROCS=2 go run hacker-cliché.go
  010101010101010101011001100101011010010100110..
  ```
* Goroutines have no identity
    * Using *thread-local storage* on goroutines may have unexpected result