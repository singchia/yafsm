package yafsm

import "testing"

/*
                              +---------+ ---------\      active OPEN
                              |  CLOSED |            \    -----------
                              +---------+<---------\   \   create TCB
                                |     ^              \   \  snd SYN
                   passive OPEN |     |   CLOSE        \   \
                   ------------ |     | ----------       \   \
                    create TCB  |     | delete TCB         \   \
                                V     |                      \   \
                              +---------+            CLOSE    |    \
                              |  LISTEN |          ---------- |     |
                              +---------+          delete TCB |     |
                   rcv SYN      |     |     SEND              |     |
                  -----------   |     |    -------            |     V
 +---------+      snd SYN,ACK  /       \   snd SYN          +---------+
 |         |<-----------------           ------------------>|         |
 |   SYN   |                    rcv SYN                     |   SYN   |
 |   RCVD  |<-----------------------------------------------|   SENT  |
 |         |                    snd ACK                     |         |
 |         |------------------           -------------------|         |
 +---------+   rcv ACK of SYN  \       /  rcv SYN,ACK       +---------+
   |           --------------   |     |   -----------
   |                  x         |     |     snd ACK
   |                            V     V
   |  CLOSE                   +---------+
   | -------                  |  ESTAB  |
   | snd FIN                  +---------+
   |                   CLOSE    |     |    rcv FIN
   V                  -------   |     |    -------
 +---------+          snd FIN  /       \   snd ACK          +---------+
 |  FIN    |<-----------------           ------------------>|  CLOSE  |
 | WAIT-1  |------------------                              |   WAIT  |
 +---------+          rcv FIN  \                            +---------+
   | rcv ACK of FIN   -------   |                            CLOSE  |
   | --------------   snd ACK   |                           ------- |
   V        x                   V                           snd FIN V
 +---------+                  +---------+                   +---------+
 |FINWAIT-2|                  | CLOSING |                   | LAST-ACK|
 +---------+                  +---------+                   +---------+
   |                rcv ACK of FIN |                 rcv ACK of FIN |
   |  rcv FIN       -------------- |    Timeout=2MSL -------------- |
   |  -------              x       V    ------------        x       V
    \ snd ACK                 +---------+delete TCB         +---------+
     ------------------------>|TIME WAIT|------------------>| CLOSED  |
                              +---------+                   +---------+

                      TCP Connection State Diagram
                              RFC793 Figure 6.
*/

const (
	// states
	CLOSED      = "closed"
	SYN_SENT    = "syn_sent"
	ESTABLISHED = "established"
	FIN_WAIT1   = "fin_wait_1"
	FIN_WAIT2   = "fin_wait_2"
	CLOSING     = "closing"
	TIME_WAIT   = "time_wait"
	CLOSE_WAIT  = "close_wait"
	LAST_ACK    = "last_ack"

	// events
	ET_CLOSE       = "close"       // SYN_SENT -> CLOSED & ESTABLISHED -> FIN_WAIT1
	ET_SENDSYN     = "sendsyn"     // CLOSED -> SYN_SENT
	ET_RECVSYNACK  = "recvsynack"  // SYN_SENT -> ESTABLISHED
	ET_SYNTIMEOUT  = "syntimeout"  // SYN_SENT -> CLOSED
	ET_SENDFIN3    = "sendfin3"    // CLOSE_WAIT -> LAST_ACK
	ET_RECVFIN1    = "recvfin1"    // FIN_WAIT2 -> TIME_WAIT
	ET_RECVFIN2    = "recvfin2"    // FIN_WAIT1 -> CLOSING
	ET_RECVFIN3    = "recvfin3"    // ESTABLISHED -> CLOSE_WAIT
	ET_RECVFINACK1 = "recvfinack1" // FIN_WAIT1 -> FIN_WAIT2
	ET_RECVFINACK2 = "recvfinack2" // CLOSING -> TIME_WAIT
	ET_RECVFINACK3 = "recvfinack3" // LAST_ACK -> CLOSED
	ET_TIMEWAITOUT = "timewaitout" // TIME_WAIT -> CLOSED
)

func initFSM() (*FSM, error) {
	fsm := NewFSM()
	// states
	closed := fsm.Init(CLOSED)
	synsent := fsm.AddState(SYN_SENT)
	established := fsm.AddState(ESTABLISHED)
	finwait1 := fsm.AddState(FIN_WAIT1)
	finwait2 := fsm.AddState(FIN_WAIT2)
	closing := fsm.AddState(CLOSING)
	timewait := fsm.AddState(TIME_WAIT)
	closewait := fsm.AddState(CLOSE_WAIT)
	lastack := fsm.AddState(LAST_ACK)

	{
		_, err := fsm.AddEvent(ET_SENDSYN, closed, synsent)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_SYNTIMEOUT, synsent, closed)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_CLOSE, synsent, closed)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVSYNACK, synsent, established)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_CLOSE, established, finwait1)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFINACK1, finwait1, finwait2)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFIN1, finwait2, timewait)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_TIMEWAITOUT, timewait, closed)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFIN2, finwait1, closing)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFINACK2, closing, timewait)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFIN3, established, closewait)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_SENDFIN3, closewait, lastack)
		if err != nil {
			return nil, err
		}
		_, err = fsm.AddEvent(ET_RECVFINACK3, lastack, closed)
		if err != nil {
			return nil, err
		}
	}
	return fsm, nil
}

func emitNormal(t *testing.T, fsm *FSM) error {
	// CLOSED -> SYN_SENT -> ESTAB -> FIN_WAIT1 -> TIME_WAIT -> CLOSED
	err := fsm.EmitEvent(ET_SENDSYN)
	if err != nil {
		return err
	}
	t.Log(fsm.State())

	err = fsm.EmitEvent(ET_RECVSYNACK)
	if err != nil {
		return err
	}
	t.Log(fsm.State())

	err = fsm.EmitEvent(ET_CLOSE)
	if err != nil {
		return err
	}
	t.Log(fsm.State())

	err = fsm.EmitEvent(ET_RECVFINACK1)
	if err != nil {
		return err
	}
	t.Log(fsm.State())

	err = fsm.EmitEvent(ET_RECVFIN1)
	if err != nil {
		return err
	}
	t.Log(fsm.State())

	err = fsm.EmitEvent(ET_TIMEWAITOUT)
	if err != nil {
		return err
	}
	t.Log(fsm.State())
	return nil
}

func TestFSM(t *testing.T) {
	fsm, err := initFSM()
	if err != nil {
		t.Error(err)
		return
	}
	err = emitNormal(t, fsm)
	if err != nil {
		t.Error(err)
		return
	}
}
