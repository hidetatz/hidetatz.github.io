What's the "sync.Cond"---2021-04-13 12:00:00

## Introduction

Go's `sync.Cond` is a synchronization primitive for multithreaded (goroutine, to be precise) programming.
Compared to mutexes, sync.Cond has limited use cases, and its usage seems to be rather complicated. Also, I don't see enough explanations or examples of it.
In this article, I will explain sync.Cond systematically to the best of my ability, not just as a reference for its usage. In addition, I will further explain and to get deep understanding of sync.Cond based on a subject that can be programmed in a better way using sync, and is likely to be found in the real world.

## What is sync.Cond?

In a nutshell, sync.Cond is a mechanism to use "**condition variables**" in Go programs.
First of all, [sync.Cond's godoc](https://golang.org/pkg/sync/#Cond) says the following.

```
Cond implements a condition variable, a rendezvous point for 
goroutines waiting for or announcing the occurrence of an event.
```

Again, the word `condition variable` appears.

Condition variables is not a Go-specific term, but (if I understand it correctly) a POSIX term. In fact, there is an interface to manipulate condition variables for [pthread](https://linux.die.net/man/3/pthread_cond_init).
There are also interfaces to manipulate condition variables for [c++](https://en.cppreference.com/w/cpp/thread/condition_variable), [Rust](https://doc.rust-lang.org/std/sync/struct.Condvar.html), [Ruby](https://docs.ruby-lang.org/ja/2.0.0/class/ConditionVariable.html), and various other languages.
Also, as you can see from the reference of each language, the interface is almost the same in all languages. ( `wait` , `signal (or notify_one)` , `broadcast (or notify_all)` )

"Condition variables" is a synchronization primitives in multithreaded programming, like semaphores and mutexes. In other words, it is a mechanism to avoid concurrent access to shared resources.
However, while semaphores and mutexes work as synchronous primitives on their own, condition variables are used in combination with mutexes. It never is used alone.
While mutexes are very simple and easy to understand as synchronization primitives, condition variables are a bit more complicated. In this article, let's deepen the understanding of condition variables by comparing it with mutexes.

### Difference between condition variables and mutexes

A mutex is a very simple mechanism to guarantee that only one thread can enter a critical section.
The situation where you want to use a mutex is when you want to perform exclusion control. In other words, you may want to use mutexes when you want to guarantee that only one thread has access to a critical section.

Condition variables appends a point to the above "I want to access a critical section" situation. The point is "I want to access a critical section, but I don't want to access it **until a condition xxx is met**". This might sound difficult to understand because it is written in a rather abstract way, but I will explain it later with a concrete example.

In order to achieve the above in mutex, it is generally necessary to implement it as a busy weight. Let's take a look at the following pseudo code.

```
mutex_lock(); // Lock the critical section
for (!condition) { // Spin until the condition gets true (busy wait)
    mutex_unlock();
    sched_yield();
    mutex_lock();
}
do_something();
mutex_unlock();
```

First, protect the critical section with a lock, then it checks the "condition" and continues to check it while looping infinitely until the condition is satisfied. In the loop, the lock is released once for other threads, and when the lock is obtained again, it returns to the head of the loop.

Busy waiting is basically an implementation that should be avoided in normal application programming because it wastes CPU resources.
With condition variables, this can be written without any spin.

```
mutex_lock();
for (!condition) {
  cond_wait(); // internally, it unlocks and locks the mutex
}
do_something();
cond_notify_all();
mutex_unlock();
```

It might look a bit confusing because it's pseudo-code (a working program in Go is described later), but `cond_wait` works internally as follows:

* Unlock locked mutex.
* Suspend its own thread and wait for notification on condition variables.
* Acquire a lock on the mutex when the notification arrives.

(See Go's [Implementation of `sync.Cond.Wait()`](https://golang.org/src/sync/cond.go?s=1353:1374#L42) for more clarification.)

When the lock is acquired, check the condition, and if true, exit the for-loop and do some processing in the critical section. Then it calls notify_all. This will wake up all other waiting threads.

Instead of `notify_all`, there is an interface `notify_one`, which wakes up not all threads but only one.

The difference between mutexes and condition variables is that mutexes need to be spun to "check conditions", while condition variables can wait for an explicit notification to be sent by an event-like mechanism.

This is similar to a situation like: a client is waiting for the success of some process on the server side, where **it is more efficient to have an event sent when the process succeeds, rather than polling the status of the process periodically**. If you take a look at the source code described below, you will see that condition variables can be used like publish-subscribe model in messaging models, as **"publish an event when a condition changes, and subscribe to it when threads want to observe the condition change"**.

In the mutex implementation, sched_yield is called to allow other threads to run, but with condition variables, the condition variable itself does this on its own.

Next, we will look at an example program that uses sync.Cond.

## sync.Cond in Real World

As an example, we will try to successfully implement a multi-threaded queue using sync.Cond. Since a simple queue does not require any condition variables, we will follow the special specification below:

* Only int type values are stored.
* There must be a maximum length, and no more than the maximum number of elements must be stored.
* Pushing an element beyond the maximum length will **block**.
* Popping an element on the empty queue will **block**.

Note that this example is based on [A simple condition variable primitive](http://joeduffyblog.com/2009/07/13/a-simple-condition-variable-primitive/).

### Implementation with mutexes

First, let's take a look at the implementation using mutexes. This source code can be found on GitHub. [dty1er/size-limited-queue/mutex_slqueue.go](https://github.com/dty1er/size-limited-queue/blob/381725020d4de089741743523ac2b032ee946767/mutex_slqueue.go)

```go
type MutexQueue struct {
	mu       sync.Mutex
	capacity int
	queue    []int
}

func NewMutexQueue(capacity int) *MutexQueue {
	return &MutexQueue{
		capacity: capacity,
		queue:    []int{},
	}
}

func (s *MutexQueue) Push(i int) {
	s.mu.Lock()
	for len(s.queue) == s.capacity {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}

	s.queue = append(s.queue, i)
	s.mu.Unlock()
}

func (s *MutexQueue) Pop() int {
	s.mu.Lock()
	for len(s.queue) == 0 {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}

	ret := s.queue[0]
	s.queue = s.queue[1:]
	s.mu.Unlock()

	return ret
}
```

Let's start with `Push()`. This program must work correctly in a multi-threaded (goroutine) environment. Also, as mentioned above, the queue has a maximum length. In other words, the confirmation of the queue length and the addition of elements must be atomic. The mutex is used to ensure the atomicity.

Because it must block if a threads try to push an element to the full queue, first make sure that the slice does not exceed the capacity, and if it does, then start spin releasing the lock until it gets a space in the queue.

`Pop()` is almost the same, atomically checking the queue length and popping the element. If the queue is empty, it blocks while spinning until an element is added (= until it can be popped).

This is exactly the "spin until a specific condition is met" described above, which is an inefficient implementation. Let's rewrite this using condition variables.

### Implementation using condition variables

The implementation using condition variables looks like this. You can find it on GitHub. [dty1er/size-limited-queue/slqueue.go](https://github.com/dty1er/size-limited-queue/blob/7482018a4aae723aebe80f0ff11f6b4f4fc265bc/slqueue.go)

```go
type SizeLimitedQueue struct {
	cond     *sync.Cond
	capacity int
	queue    []int
}

func New(capacity int) *SizeLimitedQueue {
	return &SizeLimitedQueue{
		cond:     sync.NewCond(&sync.Mutex{}),
		capacity: capacity,
		queue:    []int{},
	}
}

func (s *SizeLimitedQueue) Push(i int) {
	s.cond.L.Lock()
	for len(s.queue) == s.capacity {
		s.cond.Wait()
	}

	s.queue = append(s.queue, i)
	s.cond.Broadcast()
	s.cond.L.Unlock()
}

func (s *SizeLimitedQueue) Pop() int {
	s.cond.L.Lock()
	for len(s.queue) == 0 {
		s.cond.Wait()
	}

	ret := s.queue[0]
	s.queue = s.queue[1:]
	s.cond.Broadcast()
	s.cond.L.Unlock()

	return ret
}
```

For condition variables, [Cond in the sync package](https://golang.org/pkg/sync/#Cond) is used.
As mentioned above, condition variables are used in combination with mutexes, but since sync.Cond can have an internal mutex (`sync.Locker` more precisely) , it is used as the mutex.

Let's take a look at `Push()`. First, it explicitly locks the mutex, then checks the condition and `Wait()`. As mentioned above, in [Wait](https://golang.org/src/sync/cond.go?s=1353:1374#L42), the mutex is unlocked and locked by itself. There is no need to spin to wait for the conditions to change, the runtime will suspend and wake up the goroutine on its own.

After adding an element to the queue, call `Broadcast()`. This wakes up all goroutines that are `Wait()`, acquires locks, then rechecks the loop conditions.
 
It does almost the same thing with `Pop()`. It acquires the lock, then `Wait()`, then pops the element and calls `Broadcast`.

The use of condition variables eliminates unnecessary spins and makes the implementation more efficient. However, there are still inefficient parts of this program, and we can make improvements. We will discuss these issues to deepen our understanding of condition variables.

### Improving the implementation with condition variables

Let's summarize so far. What we can do with condition variables is to make a thread efficiently wait for some condition to be met. In the example of a queue with a maximum length, the "some condition" is currently "change of the number of elements in the queue". In both `Push()` and `Pop()`, the above program works as follows: "The number of elements in the queue changes -> wake up the thread waiting for the event of changing the number of elements in the queue -> check whether the number of elements in the queue matches the condition desired by the thread".

#### 1. Consideration of notification timing for condition variables

Now, the first point to improve is "to be able to narrow down the timing of notifications to condition variables".

Currently, `Push()` returns from `Wait()` when "the queue length is changed". However, `Push()` only needs to be notified when "the queue length was the same as the capacity (full), but it became less than the capacity". For example, if the capacity of the queue is 10 and the queue length is only 2 or 3, there is no need to send a notification because the push will never be blocked.

Similarly, for `Pop()`, it only needs to be sent the event "queue length was 0, but now it is 1".

Since we are broadcasting without considering these things now, we can make it more efficient by changing this to only do it when necessary.

#### 2. Make the condition variable meaningful

Currently, there is only one the condition variable `s.cond` for the queue, and both `Push()` and `Pop()` share the same one. However, we can find that they can be separated.

Semantically, the goroutine triggered by `Broadcast()` called in `Push()` is the one waiting to `Pop()`, not the one waiting to `Push()`. So, what is the goroutine in the `Push()` waiting for? They are waiting for the queue gets non-full; in other words, they are waiting for an element to be removed (= popped) from the queue that is full. Vice versa: "`Broadcast()` called in `Pop()`" is for goroutines waiting for an opportunity to `Push()`.

In other words, since `Push()` and `Pop()` are semantically different in what they want to notify, it makes no sense for them to share a single condition variable, `s.cond`.

Combining `1. Consideration of the timing of notifications to condition variables` and `2. Making condition variables meaningful`, we can organize them as follows.

* Change the condition variable `s.cond` : changes it to `s.nonEmptyCond` and modify it to "Notify when the queue that was empty is no longer empty". This is used for `Wait()` in `Pop()`.
* New condition variable `s.nonFullCond` : Newly defined. `s.nonFullCond` is for notifying when a queue that was full is no longer full. This is used during `Wait()` in `Push()`.

This fix changes the implementation from "all goroutines are woken up when the queue length just changes" to "each goroutine is woken up only when the queue length changes that it is interested in".

#### 3. Use `Signal()` instead of `Broadcast()`.

Changes 1 and 2 minimize the number of threads triggered, but there are still improvements to be made. For example, if the queue length changed to 1 from 0 by a `Push()`, all goroutines waiting for `Pop()` will be woken up. However, even if we wake up all the goroutines waiting for `Pop()`, only one goroutine can enter the next critical section, so it is more efficient to wake up only one thread instead of all of them. [`Signal()`](https://golang.org/pkg/sync/#Cond.Signal) can be used in this case. `Signal()` wakes up only one goroutine, not all the goroutines that are in `Wait()`. This makes the implementation more efficient by mitigating the Thundering Herd Problem when there are many waiting goroutines.

An implementation of the above 1, 2, and 3 is shown below, which can also be found on GitHub. [dty1er/sie-limited-queue/slqueue.go](https://github.com/dty1er/size-limited-queue/blob/e9dd8dc2c2937d6cffe7e784bc5c2d436632b758/slqueue.go) 
The diff is [here](https://github.com/dty1er/size-limited-queue/commit/e9dd8dc2c2937d6cffe7e784bc5c2d436632b758).

```go
type SizeLimitedQueue struct {
	nonFullCond *sync.Cond

	nonEmptyCond *sync.Cond

	capacity     int
	queue        []int
	mu           *sync.Mutex
}

func New(capacity int) *SizeLimitedQueue {
	mu := &sync.Mutex{}
	return &SizeLimitedQueue{
		nonFullCond:  sync.NewCond(mu),
		nonEmptyCond: sync.NewCond(mu),
		capacity:     capacity,
		queue:        []int{},
		mu:           mu,
	}
}

func (s *SizeLimitedQueue) Push(i int) {
	s.nonFullCond.L.Lock()
	for len(s.queue) == s.capacity {
		s.nonFullCond.Wait()
	}

	wasEmpty := len(s.queue) == 0
	s.queue = append(s.queue, i)

	if wasEmpty {
		s.nonEmptyCond.Signal()
	}
	s.nonFullCond.L.Unlock()
}

func (s *SizeLimitedQueue) Pop() int {
	s.nonEmptyCond.L.Lock()
	for len(s.queue) == 0 {
		s.nonEmptyCond.Wait()
	}

	wasFull := len(s.queue) == s.capacity
	ret := s.queue[0]
	s.queue = s.queue[1:]

	if wasFull {
		s.nonFullCond.Signal()
	}
	s.nonEmptyCond.L.Unlock()

	return ret
}
```

In order to improve the efficiency of implementation using condition variables, it is important to ensure that notifications are sent only when necessary and that only necessary threads are started.

In the next chapter, I will describe some supplementary points to keep in mind regarding condition variables.

## Condition Variables and Spurious Wakeups

First, let's talk about spurious wakeup, which you may have heard before if you have been doing multithreaded programming in Java.

Spurious wakeup refers to "the waking up of a thread that has been suspended without notification during a wait process when using a condition variable". This is usually caused by the OS or hardware, and cannot be controlled by the user of the condition variable.

In [Java documentation](https://docs.oracle.com/en/java/javase/11/docs/api/java.base/java/lang/Object.html#wait(long,int)) and [C++ documentation](https ://en.cppreference.com/w/cpp/thread/condition_variable) (and more), they explicitly state that spurious wake-ups can occur.

It is easy to deal with spurious wake-ups by always wrapping the wait process in a loop, as described in this article. Depending on the situation where you want to use a condition variable, you may not need to use looping actually, but even in such cases, you should always wait in a loop. For example, the [Developer's documentation on chrome condition variables](http://www.chromium.org/developers/lock-and-condition-variable) has a similar description.

Now, we want to know if Go's `sync.Cond` causes spurious wake-ups; This is not specified in the GoDoc as of [Go1.16.3](https://github.com/golang/go/blob/go1.16.3/src/sync/cond.go), which I referred to, and to be honest, I am not sure right now. I'll leave some of the articles and discussions I referenced below. However, I think that at least we, the users of synchronization primitives, should ensure that Wait is properly placed in the loop.

* [GopherCon and dotGo 2019 liveblogs - GopherCon 2018 - Rethinking Classical Concurrency Patterns](https://about.sourcegraph.com/go/gophercon-2018-rethinking-classical-concurrency-patterns/)
* [proposal: sync: mechanism to select on condition variables](https://github.com/golang/go/issues/16620)
* [sync.Cond and spurious wekeups](https://groups.google.com/g/golang-dev/c/Kc1nOjju3zk/discussion)
* [proposal: Go 2: sync: remove the Cond type](https://github.com/golang/go/issues/21165)

## Condition variables and the Two-step dance

About Two-step dance. If we use condition variables, one lock/unlock of a mutex that seems unnecessary can occur.
Let's see `Push()` `Pop()` in the above source again:

```go
func (s *SizeLimitedQueue) Push(i int) {
	s.nonFullCond.L.Lock()
	for len(s.queue) == s.capacity {
		s.nonFullCond.Wait()
	}

	wasEmpty := len(s.queue) == 0
	s.queue = append(s.queue, i)

	if wasEmpty {
		s.nonEmptyCond.Signal() // 1
	}
	s.nonFullCond.L.Unlock() // 3
}

func (s *SizeLimitedQueue) Pop() int {
	s.nonEmptyCond.L.Lock()
	for len(s.queue) == 0 {
		s.nonEmptyCond.Wait() // 2
	}

	wasFull := len(s.queue) == s.capacity
	ret := s.queue[0]
	s.queue = s.queue[1:]

	if wasFull {
		s.nonFullCond.Signal()
	}
	s.nonEmptyCond.L.Unlock()

	return ret
}
```

First of all, let's assume that there is a goroutine waiting to be notified by `nonEmptyCond` at `2` in the source code. Now let's say a goroutine is doing `Push()` (goroutineA). goroutineA performs a queue manipulation and sends a notification at `1`.

Here, when the notification is sent, one of the goroutines (goroutineB) that was `Wait()` at `2` is woken up from the `Wait()` state. As mentioned above (or more clearly if you refer to the source code of `sync.Cond.Wait`), in `Wait()` it tries to relock the mutex after being woken up from the wait (to enter the critical section). However, **at this point the mutex has been locked by goroutineA**. It will be unlocked at `3`. This means that the following unnecessary blocks are generated in these processes:

* goroutineB is first blocked by the condition variable `nonEmptyCond`.
* goroutineB is unblocked because the condition variable `nonEmptyCond` has been notified
* goroutineB tries to get a lock on the mutex with `2`, then **blocked**.
* goroutineA unlocks the mutex at `3`.
* goroutineB locks the mutex at `2`.

goroutineB will go through the state of unblock -> block -> unblock, which is called [Two-step dance](https://docs.microsoft.com/en-us/archive/msdn-magazine/2008/october/concurrency-hazards-solving-problems-in-your-multithreaded-code#two-step-dance).
The disadvantage of the two-step dance is that it causes unnecessary context switches.

There is a way to avoid the two-step dance; just place the `Signal()` at `1` after the `3`, i.e., to notify the condition variable after the critical section. However, this can sometimes be wrong.
The [developer documentation on condition variables in chromium](http://www.chromium.org/developers/lock-and-condition-variable) has the following.

```
In rare cases, it is incorrect to call Signal() after the critical section, 
so we recommend always using it inside the critical section. 
The following code can attempt to access the condition variable after it has been deleted, 
but could be safe if Signal() were called inside the critical section.

  struct Foo { Lock mu; ConditionVariable cv; bool del; ... };
  ...
  void Thread1(Foo *foo) {
    foo->mu.Acquire();
    while (!foo->del) {
      foo->cv.Wait();
    }
    foo->mu.Release();
    delete foo;
  }
  ...
  void Thread2(Foo *foo) {
    foo->mu.Acquire();
    foo->del = true;
                          // Signal() should be called here
    foo->mu.Release();
    foo->cv.Signal();     // BUG: foo may have been deleted
  }
```

The document also states that Chrome's implementation of condition variables detects two-step dances and delays the wake-up of the waiting thread to prevent context-switching (so notification of condition variables should always be done in the critical section).

## Summary

In Go, condition variables can be substituted with channels in some cases, so there are limited situations where they can be used. I sometimes hear people say "sync.Cond is difficult to understand", and this may be due to that.

It might be difficult to understand sync.Cond if you understand it from its behavior as "sync.Cond is used when you want to start multiple goroutines at once" because we can come up with a new question; "why do we want to do that?"

I think it would be easier to understand that "condition variables are an efficient way to wait for a certain condition to be satisfied in a multi-threaded environment" first, and then how to use them in Go.

To be honest, I don't fully understand the spurious wake-up, so I hope someone who knows more about it will write an article to explain it. By the way, when I was doing some research for this article, I found that there was a problem with [sync.WaitGroup causing spurious wake-up](https://github.com/golang/go/issues/7734) in the previous version of Go, which I found interesting.

Hope this article helps you. If you like this article, please leave it a star on the [size-limited-queue](https://github.com/dty1er/size-limited-queue) repo.

## Reference

* [A simple condition variable primitive](http://joeduffyblog.com/2009/07/13/a-simple-condition-variable-primitive/)
* [Chrome C++ Lock and ConditionVariable](http://www.chromium.org/developers/lock-and-condition-variable)
* [GopherCon 2018 Bryan C .Mills Rethinking Classical Concurrency Patterns](https://youtu.be/5zXAHh5tJqQ?t=801)
* [GopherCon and dotGo 2019 liveblogs - GopherCon 2018 - Rethinking Classical Concurrency Patterns](https://about.sourcegraph.com/go/gophercon-2018-rethinking-classical-concurrency-patterns/)
* [proposal: sync: mechanism to select on condition variables](https://github.com/golang/go/issues/16620)
* [sync.Cond and spurious wekeups](https://groups.google.com/g/golang-dev/c/Kc1nOjju3zk/discussion)
* [proposal: Go 2: sync: remove the Cond type](https://github.com/golang/go/issues/21165)
* [sync: spurious wakeup from WaitGroup.Wait](https://github.com/golang/go/issues/7734)
* [Ruby-Doc.org class ConditionVariable](https://ruby-doc.org/core-2.5.0/ConditionVariable.html)
* [Solving 11 Likely Problems In Your Multithreaded Code](https://docs.microsoft.com/en-us/archive/msdn-magazine/2008/october/concurrency-hazards-solving-problems-in-your-multithreaded-code)
* [ja - sync.Cond／コンディション変数についての解説](https://lestrrat.medium.com/sync-cond-%E3%82%B3%E3%83%B3%E3%83%87%E3%82%A3%E3%82%B7%E3%83%A7%E3%83%B3%E5%A4%89%E6%95%B0%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6%E3%81%AE%E8%A7%A3%E8%AA%AC-dd2050cdfab7)
* [ja - 条件変数 Step-by-Step入門](https://yohhoy.hatenablog.jp/entry/2014/09/23/193617)
* [ja - 条件変数とダンス(Two-Step Dance)を](https://yohhoy.hatenadiary.jp/entry/20120504/p1)
* [ja - 条件変数とデッドロック・パズル（出題編）](https://yohhoy.hatenadiary.jp/entry/20140926/p1)
* [ja - 条件変数とspurious wakeup](https://yohhoy.hatenadiary.jp/entry/20120326/p1)
* [ja - マルチスレッド・プログラミングの道具箱](https://zenn.dev/yohhoy/articles/multithreading-toolbox)
