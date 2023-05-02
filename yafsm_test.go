package yafsm

import (
	"container/list"
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/singchia/yafsm/pkg/prioqueue"
)

func TestNewState(t *testing.T) {
	type args struct {
		state string
	}
	tests := []struct {
		name string
		args args
		want *State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewState(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_AddEnter(t *testing.T) {
	type fields struct {
		State  string
		enters []StateHandler
		lefts  []StateHandler
	}
	type args struct {
		handler StateHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &State{
				State:  tt.fields.State,
				enters: tt.fields.enters,
				lefts:  tt.fields.lefts,
			}
			st.AddEnter(tt.args.handler)
		})
	}
}

func TestState_AddLeft(t *testing.T) {
	type fields struct {
		State  string
		enters []StateHandler
		lefts  []StateHandler
	}
	type args struct {
		handler StateHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &State{
				State:  tt.fields.State,
				enters: tt.fields.enters,
				lefts:  tt.fields.lefts,
			}
			st.AddLeft(tt.args.handler)
		})
	}
}

func TestEvent_AddHandler(t *testing.T) {
	type fields struct {
		Event    string
		From     *State
		To       *State
		handlers []EventHandler
		ch       chan error
	}
	type args struct {
		handler EventHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et := &Event{
				Event:    tt.fields.Event,
				From:     tt.fields.From,
				To:       tt.fields.To,
				handlers: tt.fields.handlers,
				ch:       tt.fields.ch,
			}
			et.AddHandler(tt.args.handler)
		})
	}
}

func TestNewFSM(t *testing.T) {
	tests := []struct {
		name string
		want *FSM
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFSM(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFSM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_Init(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.Init(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_Close(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			fsm.Close()
		})
	}
}

func TestFSM_SetState(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.SetState(tt.args.state); got != tt.want {
				t.Errorf("FSM.SetState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_State(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.State(); got != tt.want {
				t.Errorf("FSM.State() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_InStates(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		states []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.InStates(tt.args.states...); got != tt.want {
				t.Errorf("FSM.InStates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_AddState(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.AddState(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.AddState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_GetState(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.GetState(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.GetState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_DelState(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.DelState(tt.args.state); got != tt.want {
				t.Errorf("FSM.DelState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_AddEvent(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event    string
		from     *State
		to       *State
		handlers []EventHandler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			got, err := fsm.AddEvent(tt.args.event, tt.args.from, tt.args.to, tt.args.handlers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FSM.AddEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.AddEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_GetEvents(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.GetEvents(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.GetEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_GetEvent(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
		from  *State
		to    *State
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.GetEvent(tt.args.event, tt.args.from, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.GetEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_DelEvents(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.DelEvents(tt.args.event); got != tt.want {
				t.Errorf("FSM.DelEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_DelEvent(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
		from  *State
		to    *State
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.DelEvent(tt.args.event, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("FSM.DelEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_EmitEvent(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if err := fsm.EmitEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FSM.EmitEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFSM_EmitEventAsync(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		event string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.EmitEventAsync(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.EmitEventAsync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_EmitPrioEvent(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		prio  int
		event string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if err := fsm.EmitPrioEvent(tt.args.prio, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FSM.EmitPrioEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFSM_EmitPrioEventAsync(t *testing.T) {
	type fields struct {
		state  string
		states map[string]*State
		events map[string]*list.List
		mutex  sync.RWMutex
		pq     *prioqueue.PrioQueue
		cancel context.CancelFunc
	}
	type args struct {
		prio  int
		event string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := &FSM{
				state:  tt.fields.state,
				states: tt.fields.states,
				events: tt.fields.events,
				mutex:  tt.fields.mutex,
				pq:     tt.fields.pq,
				cancel: tt.fields.cancel,
			}
			if got := fsm.EmitPrioEventAsync(tt.args.prio, tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.EmitPrioEventAsync() = %v, want %v", got, tt.want)
			}
		})
	}
}
