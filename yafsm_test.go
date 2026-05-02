package yafsm

import (
	"errors"
	"sync"
	"testing"
	"time"
)

const (
	stateA = "A"
	stateB = "B"
	stateC = "C"
	evAB   = "a->b"
	evBC   = "b->c"
	evCA   = "c->a"
)

func newAB() *FSM {
	fsm := NewFSM()
	a := fsm.Init(stateA)
	b := fsm.AddState(stateB)
	fsm.AddEvent(evAB, a, b)
	return fsm
}

func TestNewState(t *testing.T) {
	st := NewState(stateA)
	if st.State != stateA {
		t.Fatalf("got %q want %q", st.State, stateA)
	}
}

func TestStateHandlers(t *testing.T) {
	fsm := NewFSM()
	a := fsm.Init(stateA)
	b := fsm.AddState(stateB)
	if _, err := fsm.AddEvent(evAB, a, b); err != nil {
		t.Fatal(err)
	}

	var entered, left string
	a.AddLeft(func(st *State) { left = st.State })
	b.AddEnter(func(st *State) { entered = st.State })

	if err := fsm.EmitEvent(evAB); err != nil {
		t.Fatal(err)
	}
	if left != stateA || entered != stateB {
		t.Fatalf("hooks not fired: left=%q entered=%q", left, entered)
	}
}

func TestEventHandler(t *testing.T) {
	fsm := newAB()
	ets := fsm.GetEvents(evAB)
	if len(ets) != 1 {
		t.Fatalf("want 1 event, got %d", len(ets))
	}
	called := false
	ets[0].AddHandler(func(*Event) { called = true })

	if err := fsm.EmitEvent(evAB); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("event handler not invoked")
	}
}

func TestInitAndState(t *testing.T) {
	fsm := newAB()
	if got := fsm.State(); got != stateA {
		t.Fatalf("State() = %q, want %q", got, stateA)
	}
	if !fsm.InStates(stateA, stateB) {
		t.Fatal("InStates: expected match on A")
	}
	if fsm.InStates(stateC) {
		t.Fatal("InStates: unexpected match on C")
	}
}

func TestSetState(t *testing.T) {
	fsm := newAB()
	if !fsm.SetState(stateB) {
		t.Fatal("SetState(B) should succeed")
	}
	if fsm.State() != stateB {
		t.Fatalf("expected B, got %q", fsm.State())
	}
	if fsm.SetState("nonexistent") {
		t.Fatal("SetState on unknown state should fail")
	}
}

func TestAddGetDelState(t *testing.T) {
	fsm := newAB()
	if got := fsm.GetState(stateA); got == nil || got.State != stateA {
		t.Fatalf("GetState(A): %v", got)
	}
	if got := fsm.GetState("missing"); got != nil {
		t.Fatalf("GetState(missing) should be nil, got %v", got)
	}
	if !fsm.DelState(stateB) {
		t.Fatal("DelState(B) should succeed")
	}
	if fsm.DelState(stateB) {
		t.Fatal("second DelState(B) should fail")
	}
	// associated event must have been removed
	if got := fsm.GetEvents(evAB); got != nil {
		t.Fatalf("expected events for A->B to be gone, got %v", got)
	}
}

func TestDelStateRemovesAllAssociatedEvents(t *testing.T) {
	// regression: DelState used to skip subsequent matches due to
	// iterator invalidation after list.Remove.
	fsm := NewFSM()
	a := fsm.Init(stateA)
	b := fsm.AddState(stateB)
	c := fsm.AddState(stateC)
	if _, err := fsm.AddEvent(evAB, a, b); err != nil {
		t.Fatal(err)
	}
	if _, err := fsm.AddEvent(evCA, c, a); err != nil {
		t.Fatal(err)
	}
	// share a single event name with two transitions both touching A
	shared := "shared"
	if _, err := fsm.AddEvent(shared, b, c); err != nil {
		t.Fatal(err)
	}
	if _, err := fsm.AddEvent(shared, c, b); err != nil {
		t.Fatal(err)
	}
	if !fsm.DelState(stateA) {
		t.Fatal("DelState(A) should succeed")
	}
	if got := fsm.GetEvents(evAB); got != nil {
		t.Fatalf("evAB should be cleaned: %v", got)
	}
	if got := fsm.GetEvents(evCA); got != nil {
		t.Fatalf("evCA should be cleaned: %v", got)
	}
	// shared transitions don't touch A, must remain
	if got := fsm.GetEvents(shared); len(got) != 2 {
		t.Fatalf("shared events should remain (2), got %d", len(got))
	}
}

func TestAddEventDuplicate(t *testing.T) {
	fsm := newAB()
	a := fsm.GetState(stateA)
	b := fsm.GetState(stateB)
	if _, err := fsm.AddEvent(evAB, a, b); !errors.Is(err, ErrEventDuplicated) {
		t.Fatalf("want ErrEventDuplicated, got %v", err)
	}
}

func TestAddEventIllegalSameFromDifferentTo(t *testing.T) {
	fsm := newAB()
	a := fsm.GetState(stateA)
	c := fsm.AddState(stateC)
	if _, err := fsm.AddEvent(evAB, a, c); !errors.Is(err, ErrEventIllegal) {
		t.Fatalf("want ErrEventIllegal, got %v", err)
	}
}

func TestAddEventStateNotExist(t *testing.T) {
	fsm := NewFSM()
	a := fsm.Init(stateA)
	missing := &State{State: "ghost"}
	if _, err := fsm.AddEvent(evAB, a, missing); !errors.Is(err, ErrStateNotExist) {
		t.Fatalf("want ErrStateNotExist, got %v", err)
	}
}

func TestEmitEventErrors(t *testing.T) {
	fsm := newAB()
	if err := fsm.EmitEvent("missing"); !errors.Is(err, ErrEventNotExist) {
		t.Fatalf("want ErrEventNotExist, got %v", err)
	}
	// fsm starts in A; evBC needs B
	b := fsm.AddState(stateB)
	c := fsm.AddState(stateC)
	if _, err := fsm.AddEvent(evBC, b, c); err != nil {
		t.Fatal(err)
	}
	if err := fsm.EmitEvent(evBC); !errors.Is(err, ErrIllegalStateForEvent) {
		t.Fatalf("want ErrIllegalStateForEvent, got %v", err)
	}
}

func TestEmitEventAsync(t *testing.T) {
	fsm := newAB()
	ch := fsm.EmitEventAsync(evAB)
	if err := <-ch; err != nil {
		t.Fatal(err)
	}
	if fsm.State() != stateB {
		t.Fatalf("expected B, got %q", fsm.State())
	}
}

func TestEmitEventAsyncEventNotExist(t *testing.T) {
	fsm := newAB()
	if err := <-fsm.EmitEventAsync("missing"); !errors.Is(err, ErrEventNotExist) {
		t.Fatalf("want ErrEventNotExist, got %v", err)
	}
}

func TestDelEvents(t *testing.T) {
	fsm := newAB()
	if !fsm.DelEvents(evAB) {
		t.Fatal("DelEvents should succeed")
	}
	if fsm.DelEvents(evAB) {
		t.Fatal("second DelEvents should fail")
	}
	if err := fsm.EmitEvent(evAB); !errors.Is(err, ErrEventNotExist) {
		t.Fatalf("want ErrEventNotExist, got %v", err)
	}
}

func TestDelEvent(t *testing.T) {
	fsm := newAB()
	a := fsm.GetState(stateA)
	b := fsm.GetState(stateB)
	if !fsm.DelEvent(evAB, a, b) {
		t.Fatal("DelEvent should succeed")
	}
	if fsm.DelEvent(evAB, a, b) {
		t.Fatal("second DelEvent should fail")
	}
}

func TestGetEvent(t *testing.T) {
	fsm := newAB()
	a := fsm.GetState(stateA)
	b := fsm.GetState(stateB)
	if got := fsm.GetEvent(evAB, a, b); got == nil {
		t.Fatal("GetEvent should find event")
	}
	if got := fsm.GetEvent("missing", a, b); got != nil {
		t.Fatal("GetEvent for unknown event should be nil")
	}
}

func TestEmitPrioEvent(t *testing.T) {
	fsm := newAB()
	if err := fsm.EmitPrioEvent(2, evAB); err != nil {
		t.Fatal(err)
	}
	if fsm.State() != stateB {
		t.Fatalf("expected B, got %q", fsm.State())
	}
}

func TestEmitPrioEventAsync(t *testing.T) {
	fsm := newAB()
	if err := <-fsm.EmitPrioEventAsync(5, evAB); err != nil {
		t.Fatal(err)
	}
	if fsm.State() != stateB {
		t.Fatalf("expected B, got %q", fsm.State())
	}
}

func TestAsyncMode(t *testing.T) {
	fsm := NewFSM(WithAsync())
	defer fsm.Close()
	a := fsm.Init(stateA)
	b := fsm.AddState(stateB)
	if _, err := fsm.AddEvent(evAB, a, b); err != nil {
		t.Fatal(err)
	}
	if err := fsm.EmitEvent(evAB); err != nil {
		t.Fatal(err)
	}
	if fsm.State() != stateB {
		t.Fatalf("expected B, got %q", fsm.State())
	}
}

func TestInSeqModeConcurrentEmits(t *testing.T) {
	fsm := NewFSM(WithInSeq())
	a := fsm.Init(stateA)
	b := fsm.AddState(stateB)
	c := fsm.AddState(stateC)
	if _, err := fsm.AddEvent(evAB, a, b); err != nil {
		t.Fatal(err)
	}
	if _, err := fsm.AddEvent(evBC, b, c); err != nil {
		t.Fatal(err)
	}
	if _, err := fsm.AddEvent(evCA, c, a); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); fsm.EmitEvent(evAB) }()
		go func() { defer wg.Done(); fsm.EmitEvent(evBC) }()
		go func() { defer wg.Done(); fsm.EmitEvent(evCA) }()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("emit deadlock under InSeq mode")
	}
}

func TestClose(t *testing.T) {
	fsm := newAB()
	fsm.Close()
}
