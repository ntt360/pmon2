//Package garlic is a simple proc connector interface for golang.
package proc

/*
GArLIC: GolAng LInux Connector: Linux Processor Connector library
*/

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"time"

	"github.com/mdlayher/netlink"
)

//parse and handle the event Interface
func getEvent(hdr procEventHdr, data []byte) (EventData, error) {
	switch hdr.What {
	case ProcEventNone:
		//We should only see this when we're getting an ACK back from the kernel
		return nil, fmt.Errorf("Got ProcEventNone")
	case ProcEventFork:
		ev := Fork{}
		ev.ParentPid, ev.ParentTgid, ev.ChildPid, ev.ChildTgid = return4Uint32(data)
		return ev, nil
	case ProcEventExec:
		ev := Exec{}
		ev.ProcessPid, ev.ProcessTgid = return2Uint32(data)
		return ev, nil
	case ProcEventUID, ProcEventGID:
		ev := ID{}
		ev.ProcessPid, ev.ProcessTgid, ev.RealID, ev.EffectiveID = return4Uint32(data)
		return ev, nil
	case ProcEventSID:
		ev := Sid{}
		ev.ProcessPid, ev.ProcessTgid = return2Uint32(data)
		return ev, nil
	case ProcEventPtrace:
		ev := Ptrace{}
		ev.ProcessPid, ev.ProcessTgid, ev.TracerPid, ev.TracerTgid = return4Uint32(data)
		return ev, nil
	case ProcEventComm:
		ev := Comm{}
		ev.ProcessPid, ev.ProcessTgid = return2Uint32(data)
		ev.Comm = string(data[8:bytes.IndexByte(data[8:], 0)])
		//copy(ev.Comm[:], data[8:])
		return ev, nil
	case ProcEventCoredump:
		ev := Coredump{}
		ev.ProcessPid, ev.ProcessTgid = return2Uint32(data)
		return ev, nil
	case ProcEventExit:
		ev := Exit{}
		ev.ProcessPid, ev.ProcessTgid, ev.ExitCode, ev.ExitSignal = return4Uint32(data)
		return ev, nil
	}

	return Exit{}, fmt.Errorf("Unknown What: %x", hdr.What)
}

func (c CnConn) parseCn(data []byte) (ProcEvent, error) {

	hdr := unmarshalProcEventHdr(data[cnMsgLen:])
	//buf := bytes.NewBuffer(data[cnMsgLen+procEventHdrLen:])

	ev, err := getEvent(hdr, data[cnMsgLen+procEventHdrLen:])
	if err != nil {
		return ProcEvent{}, err
	}

	ts := time.Unix(0, int64(hdr.Timestamp)+(c.boottime*1000000000))

	return ProcEvent{What: hdr.What, CPU: hdr.CPU, Timestamp: ts, EventData: ev, WhatString: evtType2Str(hdr.What)}, nil
}

//ClosePCN closes the netlink  connection
func (c CnConn) ClosePCN() error {
	return c.c.Close()
}

//ReadPCN reads waits for a Proc connector event to come across the nl socket, and returns an event struct
//This is a blocking operation
func (c CnConn) ReadPCN() ([]ProcEvent, error) {
	retMsg, err := c.c.Receive()
	if err != nil {
		return nil, fmt.Errorf("Receive error: %s", err)
	}

	//I've never seen these underlying libs return more than one proc event, but lets not make assumptions
	evList := make([]ProcEvent, len(retMsg))
	for iter, value := range retMsg {
		parsedEv, err := c.parseCn(value.Data)
		if err != nil {
			return nil, fmt.Errorf("Bad parseCn: %s", err)
		}
		evList[iter] = parsedEv

	}

	return evList, nil
}

func dialPCN() (CnConn, error) {

	//DialPCN Config
	cCfg := netlink.Config{Groups: 0x1}

	//Bind
	c, err := netlink.Dial(syscall.NETLINK_CONNECTOR, &cCfg)
	//fmt.Println("Finished dial.")

	if err != nil {
		return CnConn{}, fmt.Errorf("Error in netlink: %s", err)
	}

	//setup process connector hdr
	cbHdr := cbID{Idx: CnIdxProc, Val: CnValProc}
	var connBody uint32 = ProcCnMcastListen
	cnHdr := cnMsg{ID: cbHdr, Len: uint16(binary.Size(connBody))}

	binHdr := cnHdr.marshalBinaryAndBody(connBody)

	reqMsg := netlink.Message{
		Header: netlink.Header{
			Type: syscall.NLMSG_DONE,
		},
		Data: binHdr,
	}

	//Send request message
	msgs, err := c.Send(reqMsg)
	if err != nil {
		return CnConn{}, fmt.Errorf("Execute error: %s\n %#v", err, msgs)
	}

	//Wait for our ack msg
	//ack, err := c.Receive()
	//if err != nil {
	//	return CnConn{}, fmt.Errorf("could not recv ack: %v", err)
	//}

	//check to make sure out ack valid
	/*if !isAck(ack[0].Data) {
		return CnConn{}, fmt.Errorf("Packet not a valid ACK: %+v", ack)
	}*/

	//get the system boot time to calculate the ktime timestamps
	bt, err := getBoottime()
	if err != nil {
		return CnConn{}, err
	}

	return CnConn{c: c, boottime: bt}, nil
}

//DialPCN connects to the proc connector socket, and returns a connection that will listens for all available event types:
//None, Fork, Execm UID, GID, SID, Ptrace, Comm, Coredump and Exit
func DialPCN() (CnConn, error) {
	return dialPCN()
}

//DialPCNWithEvents is the same as DialPCN(), but with a filter that allows you select a particular proc event.
//It uses bitmasks and PBF to filter for the given events
func DialPCNWithEvents(events []EventType) (CnConn, error) {

	c, err := dialPCN()
	if err != nil {
		return CnConn{}, err
	}

	filters, err := loadBPF(events)
	if err != nil {
		return CnConn{}, err
	}

	err = c.c.SetBPF(filters)
	if err != nil {
		return CnConn{}, err
	}

	return c, nil

}
