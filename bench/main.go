package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"syscall"

	"github.com/jumboframes/armorigo/sigaction"
	"github.com/singchia/yafsm"
)

const (
	INIT       = "init"
	CONN_SENT  = "conn_sent"
	CONN_RECV  = "conn_recv"
	CONNED     = "conned"
	ABNORMAL   = "abnormal"
	CLOSE_SENT = "close_sent"
	CLOSE_RECV = "close_recv"
	CLOSE_HALF = "close_half"
	CLOSED     = "closed"
	FINI       = "fini"

	ET_CONNSENT  = "connsent"
	ET_CONNRECV  = "connrecv"
	ET_CONNACK   = "connack"
	ET_ERROR     = "error"
	ET_EOF       = "eof"
	ET_CLOSESENT = "closesent"
	ET_CLOSERECV = "closerecv"
	ET_CLOSEACK  = "closeack"
	ET_FINI      = "fini"
)

func initFSM(fsm *yafsm.FSM) {
	init := fsm.AddState(INIT)
	connrecv := fsm.AddState(CONN_RECV)
	conned := fsm.AddState(CONNED)
	closesent := fsm.AddState(CLOSE_SENT)
	closerecv := fsm.AddState(CLOSE_RECV)
	closehalf := fsm.AddState(CLOSE_HALF)
	closed := fsm.AddState(FINI)
	fini := fsm.AddState(FINI)
	fsm.SetState(INIT)

	// events
	fsm.AddEvent(ET_CONNRECV, init, connrecv)
	fsm.AddEvent(ET_CONNACK, connrecv, conned)

	fsm.AddEvent(ET_ERROR, init, closed)
	fsm.AddEvent(ET_ERROR, connrecv, closed)
	fsm.AddEvent(ET_ERROR, conned, closed)
	fsm.AddEvent(ET_ERROR, closesent, closed)
	fsm.AddEvent(ET_ERROR, closerecv, closed)

	fsm.AddEvent(ET_EOF, connrecv, closed)
	fsm.AddEvent(ET_EOF, conned, closed)

	fsm.AddEvent(ET_CLOSESENT, conned, closesent)
	fsm.AddEvent(ET_CLOSESENT, closerecv, closesent) // close and been closed at same time
	fsm.AddEvent(ET_CLOSESENT, closehalf, closehalf) // close and been closed at same time

	fsm.AddEvent(ET_CLOSERECV, conned, closerecv)
	fsm.AddEvent(ET_CLOSERECV, closesent, closerecv) // close and been closed at same time
	fsm.AddEvent(ET_CLOSERECV, closehalf, closehalf) // close and been closed at same time

	fsm.AddEvent(ET_CLOSEACK, closesent, closehalf)
	fsm.AddEvent(ET_CLOSEACK, closerecv, closehalf)
	fsm.AddEvent(ET_CLOSEACK, closehalf, closed)
	// fini
	fsm.AddEvent(ET_FINI, init, fini)
	fsm.AddEvent(ET_FINI, connrecv, fini)
	fsm.AddEvent(ET_FINI, conned, fini)
	fsm.AddEvent(ET_FINI, closesent, fini)
	fsm.AddEvent(ET_FINI, closerecv, fini)
	fsm.AddEvent(ET_FINI, closehalf, fini)
	fsm.AddEvent(ET_FINI, closed, fini)
}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:6061", nil)
	}()

	fsms := []*yafsm.FSM{}
	count := 100000
	for i := 0; i < count; i++ {
		fsm := yafsm.NewFSM()
		fsms = append(fsms, fsm)
		initFSM(fsm)
		fsm.EmitEvent(ET_CONNRECV)
		fsm.EmitEvent(CLOSE_SENT)
		fsm.EmitEvent(ET_ERROR)
		fsm.EmitEvent(ET_FINI)
		fsm.Close()
	}
	fmt.Println("done")

	sig := sigaction.NewSignal()
	sig.Add(syscall.SIGINT)
	sig.Wait(context.Background())
}
