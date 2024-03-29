package yafsm

import (
	"container/list"
	"context"
	"sync"

	"github.com/singchia/yafsm/pkg/prioqueue"
)

type StateHandler func(st *State)

type State struct {
	State  string
	enters []StateHandler
	lefts  []StateHandler
}

func NewState(state string) *State {
	return &State{State: state}
}

func (st *State) AddEnter(handler StateHandler) {
	st.enters = append(st.enters, handler)
}

func (st *State) AddLeft(handler StateHandler) {
	st.lefts = append(st.lefts, handler)
}

type EventHandler func(event *Event)

type Event struct {
	Event    string
	From, To *State
	handlers []EventHandler
	ch       chan error
}

type dup byte

const (
	dupLegal       dup = 0x00
	dupEvent       dup = 0x01
	dupEventFrom   dup = 0x02
	dupEventFromTo dup = 0x03
)

func (et *Event) duplicate(event, from, to string) dup {
	if event == et.Event && from == et.From.State && to == et.To.State {
		return dupEventFromTo
	} else if event == et.Event && from == et.From.State {
		return dupEventFrom
	} else if event == et.Event {
		return dupEvent
	} else {
		return dupLegal
	}
}

func (et *Event) AddHandler(handler EventHandler) {
	et.handlers = append(et.handlers, handler)
}

type FSMOption func(*FSM)

func WithAsync() FSMOption {
	return func(fsm *FSM) {
		fsm.async = true
	}
}

func WithInSeq() FSMOption {
	return func(fsm *FSM) {
		fsm.inseq = true
	}
}

type FSM struct {
	state  string
	states map[string]*State
	events map[string]*list.List

	async, inseq bool
	mutex        sync.RWMutex
	pq           *prioqueue.PrioQueue
	cancel       context.CancelFunc
}

func NewFSM(opts ...FSMOption) *FSM {
	pq, _ := prioqueue.NewPrioQueue()
	ctx, cancel := context.WithCancel(context.Background())
	fsm := &FSM{
		pq:     pq,
		cancel: cancel,
		states: make(map[string]*State),
		events: make(map[string]*list.List),
	}
	for _, opt := range opts {
		opt(fsm)
	}
	if fsm.async {
		go fsm.emit(ctx)
	}
	return fsm
}

func (fsm *FSM) Init(state string) *State {
	fsm.state = state
	return fsm.AddState(state)
}

func (fsm *FSM) Close() {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	for k, _ := range fsm.states {
		delete(fsm.states, k)
	}
	for k, v := range fsm.events {
		for elem := v.Front(); elem != nil; elem = elem.Next() {
			event := elem.Value.(*Event)
			event.handlers, event.From, event.To = nil, nil, nil
		}
		delete(fsm.events, k)
	}
	fsm.pq.Close()
	fsm.cancel()
}

func (fsm *FSM) emit(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fsm.emitOne()
		}
	}
}

func (fsm *FSM) emitOne() {
	data := fsm.pq.PopSync()
	if data == nil {
		return
	}
	switch ec := data.(type) {
	case *eventchan:
		et := (*Event)(nil)

		fsm.mutex.RLock()
		etList, ok := fsm.events[ec.event]
		if !ok {
			ec.ch <- ErrEventNotExist
			close(ec.ch)
			fsm.mutex.RUnlock()
			return
		}
		fsm.mutex.RUnlock()

		fsm.mutex.Lock()
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			tmp := elem.Value.(*Event)
			if tmp.From.State == fsm.state {
				et = tmp
			}
		}
		if et == nil {
			ec.ch <- ErrIllegalStateForEvent
			close(ec.ch)
			fsm.mutex.Unlock()
			return
		}
		fsm.state = et.To.State
		fsm.mutex.Unlock()

		for _, left := range et.From.lefts {
			left(et.From)
		}
		for _, handler := range et.handlers {
			handler(et)
		}
		for _, enter := range et.To.enters {
			enter(et.To)
		}
		ec.ch <- nil
		close(ec.ch)
	}
}

func (fsm *FSM) emitOneInSeq() {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	data := fsm.pq.PopSync()
	if data == nil {
		return
	}
	switch ec := data.(type) {
	case *eventchan:
		et := (*Event)(nil)
		etList, ok := fsm.events[ec.event]
		if !ok {
			ec.ch <- ErrEventNotExist
			close(ec.ch)
			return
		}
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			tmp := elem.Value.(*Event)
			if tmp.From.State == fsm.state {
				et = tmp
			}
		}
		if et == nil {
			ec.ch <- ErrIllegalStateForEvent
			close(ec.ch)
			return
		}
		fsm.state = et.To.State
		for _, left := range et.From.lefts {
			left(et.From)
		}
		for _, handler := range et.handlers {
			handler(et)
		}
		for _, enter := range et.To.enters {
			enter(et.To)
		}
		ec.ch <- nil
		close(ec.ch)
	}
}

func (fsm *FSM) SetState(state string) bool {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	_, ok := fsm.states[state]
	if !ok {
		return false
	}
	fsm.state = state
	return true
}

func (fsm *FSM) State() string {
	fsm.mutex.RLock()
	defer fsm.mutex.RUnlock()
	return fsm.state
}

func (fsm *FSM) InStates(states ...string) bool {
	fsm.mutex.RLock()
	defer fsm.mutex.RUnlock()

	for _, state := range states {
		if state == fsm.state {
			return true
		}
	}
	return false
}

func (fsm *FSM) AddState(state string) *State {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	st, ok := fsm.states[state]
	if !ok {
		st = &State{State: state}
		fsm.states[state] = st
	}
	return st
}

func (fsm *FSM) GetState(state string) *State {
	fsm.mutex.RLock()
	defer fsm.mutex.RUnlock()

	st, ok := fsm.states[state]
	if ok {
		return st
	}
	return nil
}

func (fsm *FSM) DelState(state string) bool {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	_, ok := fsm.states[state]
	if !ok {
		return false
	}
	delete(fsm.states, state)

	et := (*Event)(nil)
	for event, etList := range fsm.events {
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			et = elem.Value.(*Event)
			if et.From.State == state ||
				et.To.State == state {
				etList.Remove(elem)
			}
		}
		if etList.Len() == 0 {
			delete(fsm.events, event)
		}
	}
	return true
}

func (fsm *FSM) stateExists(state string) bool {
	_, ok := fsm.states[state]
	return ok
}

func (fsm *FSM) AddEvent(event string, from, to *State,
	handlers ...EventHandler) (*Event, error) {

	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	if !fsm.stateExists(from.State) || !fsm.stateExists(to.State) {
		return nil, ErrStateNotExist
	}
	et := (*Event)(nil)
	etList, ok := fsm.events[event]
	if ok {
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			et = elem.Value.(*Event)
			dup := et.duplicate(event, from.State, to.State)
			switch dup {
			case dupEventFromTo:
				// same event, same from, same to
				return nil, ErrEventDuplicated
			case dupEventFrom:
				// same event, same from, different to
				return nil, ErrEventIllegal
			}
		}
		et = &Event{
			Event:    event,
			From:     from,
			To:       to,
			handlers: handlers,
		}
		etList.PushBack(et)

	} else {
		etList := list.New()
		et = &Event{
			Event:    event,
			From:     from,
			To:       to,
			handlers: handlers,
		}
		etList.PushBack(et)
		fsm.events[event] = etList
	}
	return et, nil
}

func (fsm *FSM) GetEvents(event string) []*Event {
	fsm.mutex.RLock()
	defer fsm.mutex.RUnlock()

	et := (*Event)(nil)
	ets := []*Event{}
	etList, ok := fsm.events[event]
	if ok {
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			et = elem.Value.(*Event)
			ets = append(ets, et)
		}
		return ets
	}
	return nil
}

func (fsm *FSM) GetEvent(event string, from, to *State) *Event {
	fsm.mutex.RLock()
	defer fsm.mutex.RUnlock()

	et := (*Event)(nil)
	etList, ok := fsm.events[event]
	if ok {
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			et = elem.Value.(*Event)
			if et.Event == event &&
				et.From == from &&
				et.To == to {
				return et
			}
		}
		return nil
	}
	return nil
}

func (fsm *FSM) DelEvents(event string) bool {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	_, ok := fsm.events[event]
	if ok {
		delete(fsm.events, event)
		return true
	}
	return false
}

func (fsm *FSM) DelEvent(event string, from, to *State) bool {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	et := (*Event)(nil)
	etList, ok := fsm.events[event]
	if ok {
		deleted := false
		for elem := etList.Front(); elem != nil; elem = elem.Next() {
			et = elem.Value.(*Event)
			if et.Event == event &&
				et.From == from &&
				et.To == to {
				etList.Remove(elem)
				deleted = true
				break
			}
		}
		if deleted {
			if etList.Len() == 0 {
				delete(fsm.events, event)
			}
			return true
		}
	}
	return false
}

type eventchan struct {
	event string
	ch    chan error
}

func (fsm *FSM) EmitEvent(event string) error {
	fsm.mutex.RLock()
	ch := make(chan error, 1)
	_, ok := fsm.events[event]
	fsm.mutex.RUnlock()
	if !ok {
		return ErrEventNotExist
	}

	eventchan := &eventchan{
		event: event,
		ch:    ch,
	}
	err := fsm.pq.Push(eventchan)
	if err != nil {
		return err
	}
	if !fsm.async {
		if fsm.inseq {
			fsm.emitOneInSeq()
		} else {
			fsm.emitOne()
		}
	}
	err = <-ch
	return err
}

func (fsm *FSM) EmitEventAsync(event string) <-chan error {
	fsm.mutex.RLock()
	ch := make(chan error, 1)
	_, ok := fsm.events[event]
	fsm.mutex.RUnlock()
	if !ok {
		ch <- ErrEventNotExist
		return ch
	}

	eventchan := &eventchan{
		event: event,
		ch:    ch,
	}

	err := fsm.pq.Push(eventchan)
	if err != nil {
		ch <- err
		return ch
	}
	if !fsm.async {
		if fsm.inseq {
			fsm.emitOneInSeq()
		} else {
			fsm.emitOne()
		}
	}
	return ch
}

func (fsm *FSM) EmitPrioEvent(prio int, event string) error {
	fsm.mutex.RLock()
	ch := make(chan error, 1)
	_, ok := fsm.events[event]
	fsm.mutex.RUnlock()
	if !ok {
		return ErrEventNotExist
	}

	eventchan := &eventchan{
		event: event,
		ch:    ch,
	}

	err := fsm.pq.PrioPush(prio, eventchan)
	if err != nil {
		return err
	}
	if !fsm.async {
		if fsm.inseq {
			fsm.emitOneInSeq()
		} else {
			fsm.emitOne()
		}
	}
	err = <-ch
	return err
}

func (fsm *FSM) EmitPrioEventAsync(prio int, event string) <-chan error {
	fsm.mutex.RLock()
	ch := make(chan error, 1)
	_, ok := fsm.events[event]
	fsm.mutex.RUnlock()
	if !ok {
		ch <- ErrEventNotExist
		return ch
	}

	eventchan := &eventchan{
		event: event,
		ch:    ch,
	}

	err := fsm.pq.PrioPush(prio, eventchan)
	if err != nil {
		ch <- err
		return ch
	}
	if !fsm.async {
		if fsm.inseq {
			fsm.emitOneInSeq()
		} else {
			fsm.emitOne()
		}
	}
	return ch
}
