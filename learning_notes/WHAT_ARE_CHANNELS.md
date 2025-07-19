# What Are Go Channels and Why Haven't We Used Them Yet?

This note explains the concept of Go channels and clarifies why we have used `goroutines` but not yet `channels` in our project.

## What Are Channels?

Think of a channel as a **conveyor belt** that connects two or more goroutines.

-   You can put things on the conveyor belt at one end (**sending** data).
-   Another goroutine can take things off the belt at the other end (**receiving** data).
-   The conveyor belt is typed: you can have a channel specifically for `int`s, `string`s, or any other type.

Channels are Go's primary tool for **communication** and **synchronization** between goroutines. They allow one goroutine to safely pass data to another without them directly accessing the same memory at the same time.

## Two Models of Concurrency

There are two main ways to handle concurrency:

1.  **Shared Memory with Locks (The `Mutex` approach):**
    -   This is what we did in our `InMemoryEventBus`.
    -   Multiple goroutines access the same piece of memory (our `subscribers` map).
    -   We use a `sync.Mutex` to act as a traffic cop, ensuring only one goroutine can *change* the map at a time to prevent chaos.

2.  **Communicating Sequential Processes (The `Channel` approach):**
    -   This is the model that channels enable.
    -   Each goroutine has its own memory and doesn't share it.
    -   When they need to coordinate, they pass messages to each other over channels.

There is a famous proverb in the Go community that captures this philosophy:

> "Do not communicate by sharing memory; instead, share memory by communicating."

## Why We Haven't Checked It Off Yet

In our `InMemoryEventBus`, we used goroutines to run our handlers concurrently, but they are "fire-and-forget." The `Publish` method launches them and immediately moves on. The handlers don't communicate back to the `Publish` method or to each other. They are independent workers.

We have used `goroutines` for **concurrency**, but we haven't yet used `channels` for **communication** or **synchronization** between them.

### A Simple Example Where We *Would* Use a Channel

Imagine if the `Publish` method needed to know when all the handlers were finished before it could proceed. We could create a channel and have each handler send a "done" signal back on the channel when it's finished. The `Publish` method would then wait until it received a "done" signal from every handler it started. This is a form of synchronization that channels make easy and safe.
