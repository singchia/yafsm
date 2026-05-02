# yafsm

[![Go](https://github.com/singchia/yafsm/actions/workflows/go.yml/badge.svg)](https://github.com/singchia/yafsm/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/singchia/yafsm.svg)](https://pkg.go.dev/github.com/singchia/yafsm)
[![Go Report Card](https://goreportcard.com/badge/github.com/singchia/yafsm)](https://goreportcard.com/report/github.com/singchia/yafsm)
[![License](https://img.shields.io/github/license/singchia/yafsm)](LICENSE)

Yet another finite state machine for Go.

`yafsm` is a small, dependency-light FSM library focused on three things:

- **Concurrency-safe** transitions guarded by an internal `sync.RWMutex`.
- **Hooks** on states (enter/leave) and on events (transition handlers).
- **Three emission modes** — synchronous, async via a background dispatcher, or strictly serialized — plus an internal **priority queue** for prioritized events.

## Install

```bash
go get github.com/singchia/yafsm
```

Requires Go 1.20+.

## Quickstart

```go
package main

import (
    "fmt"

    "github.com/singchia/yafsm"
)

const (
    Closed = "closed"
    Open   = "open"

    OpenIt  = "open"
    CloseIt = "close"
)

func main() {
    fsm := yafsm.NewFSM()

    closed := fsm.Init(Closed)
    open := fsm.AddState(Open)

    closed.AddLeft(func(s *yafsm.State) { fmt.Println("leaving", s.State) })
    open.AddEnter(func(s *yafsm.State) { fmt.Println("entering", s.State) })

    if _, err := fsm.AddEvent(OpenIt, closed, open); err != nil {
        panic(err)
    }
    if _, err := fsm.AddEvent(CloseIt, open, closed); err != nil {
        panic(err)
    }

    fmt.Println("state:", fsm.State()) // closed
    if err := fsm.EmitEvent(OpenIt); err != nil {
        panic(err)
    }
    fmt.Println("state:", fsm.State()) // open
}
```

## Concepts

A yafsm `FSM` is a directed graph of:

- **States**, identified by their string name. Each state can carry any number of *enter* and *leave* handlers.
- **Events**, identified by name. An event is bound to exactly one `(from, to)` pair; the same event name may be reused for different `from` states. Each event can carry any number of handlers fired during the transition.

The library enforces two invariants when adding events:

| Condition | Result |
| --- | --- |
| Same `event` + `from` + `to` already registered | `ErrEventDuplicated` |
| Same `event` + `from`, different `to`            | `ErrEventIllegal` |
| Either `from` or `to` not registered yet         | `ErrStateNotExist` |

Emitting an event when the FSM is not in the event's `from` state returns `ErrIllegalStateForEvent`. Emitting an unknown event returns `ErrEventNotExist`.

## Emission modes

| Constructor | Behaviour |
| --- | --- |
| `NewFSM()` | Synchronous: `EmitEvent` runs the transition on the caller goroutine. Concurrent `EmitEvent` calls may interleave their lock windows but never corrupt state. |
| `NewFSM(WithInSeq())` | Strictly serialized: each transition (including all its enter/leave/event handlers) runs while holding the FSM lock, so no other transition can interleave. Use when handlers mutate shared state and you need linearizability. |
| `NewFSM(WithAsync())` | A background goroutine drains the queue and runs transitions one at a time. `EmitEvent` blocks until completion; `EmitEventAsync` returns a `<-chan error`. Call `Close()` to stop the worker. |

A built-in priority queue lets you push higher-priority events ahead of pending ones via `EmitPrioEvent` / `EmitPrioEventAsync`. Larger priority value = higher priority.

## API at a glance

```go
fsm := yafsm.NewFSM(opts ...FSMOption)            // WithAsync, WithInSeq

// states
state := fsm.Init("idle")                          // or fsm.AddState
state.AddEnter(func(*yafsm.State) {})
state.AddLeft(func(*yafsm.State) {})

// events
ev, err := fsm.AddEvent("go", from, to, handlers...)
ev.AddHandler(func(*yafsm.Event) {})

// inspection
fsm.State()                                        // current state
fsm.InStates("a", "b")                             // is current in any of these
fsm.GetState("idle"); fsm.GetEvent("go", from, to)
fsm.GetEvents("go")

// mutation
fsm.SetState("idle")                               // skip the transition pipeline
fsm.DelState("idle")                               // also drops events touching it
fsm.DelEvent("go", from, to); fsm.DelEvents("go")

// emission
err := fsm.EmitEvent("go")
ch  := fsm.EmitEventAsync("go")                    // <-chan error
err := fsm.EmitPrioEvent(prio, "go")
ch  := fsm.EmitPrioEventAsync(prio, "go")

// teardown (mandatory in async mode)
fsm.Close()
```

See [`pkg.go.dev`](https://pkg.go.dev/github.com/singchia/yafsm) for the full reference.

## Examples

A complete TCP connection state machine modelled after RFC 793 lives in
[`examples/tcp_state_transition`](examples/tcp_state_transition/tcp_state_transition.go) — it exercises every emission mode, including async event handlers and priority preemption.

## Benchmark

A simple stress benchmark that creates 100k FSMs and walks them through a connection lifecycle is under [`bench/`](bench/main.go):

```bash
go run ./bench
```

It exposes `pprof` on `:6061` so you can attach `go tool pprof` while it runs.

## License

[Apache-2.0](LICENSE)
