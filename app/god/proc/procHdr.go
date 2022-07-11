package proc

import (
	"encoding/binary"
	"time"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
)

//CnConn contains the connection to the proc connector socket
type CnConn struct {
	c        *netlink.Conn
	boottime int64
}

//Preserve a lot of the const/enums from cn_proc.h for the sake of documentation
//No one wants see 0x1 everywhere

const ( //proc_cn_mcast_op
	//ProcCnMcastListen registers a listen event with the kernel
	ProcCnMcastListen = iota + 1
	//ProcCnMcastIgnore registers an ignore event with the kernel
	ProcCnMcastIgnore = iota
)

//CnIdxProc is the Id used for proc/connector, and is a unique identifier which is used for message routing and must be registered in connector.h for in-kernel usage.
const CnIdxProc = 0x1

//CnValProc is the corrisponding value used by chID,  and is a unique identifier which is used for message routing and must be registered in connector.h for in-kernel usage.
const CnValProc = 0x1

//Various message structs from connector.h

/*
 * idx and val are unique identifiers which
 * are used for message routing and
 * must be registered in connector.h for in-kernel usage.
 */
type cbID struct {
	Idx uint32
	Val uint32
}

//The connector message struct
type cnMsg struct {
	ID    cbID
	Seq   uint32
	Ack   uint32
	Len   uint16
	Flags uint16
}

var cnMsgLen = binary.Size(cnMsg{})

//MarshallBinary converts the entire struct into a slice, along with the proc_cn_mcast_op body
func (hdr cnMsg) marshalBinaryAndBody(body uint32) []byte {

	bytes := make([]byte, binary.Size(hdr)+binary.Size(body))
	nlenc.PutUint32(bytes[0:4], hdr.ID.Idx)
	nlenc.PutUint32(bytes[4:8], hdr.ID.Val)
	nlenc.PutUint32(bytes[8:12], hdr.Seq)
	nlenc.PutUint32(bytes[12:16], hdr.Ack)
	nlenc.PutUint16(bytes[16:18], hdr.Len)
	nlenc.PutUint16(bytes[18:20], hdr.Flags)
	nlenc.PutUint32(bytes[20:24], body)

	return bytes
}

//This is just an internal  header that allows us to easily cast the raw binary data
type procEventHdr struct {
	What      EventType
	CPU       uint32
	Timestamp uint64
}

var procEventHdrLen = binary.Size(procEventHdr{})

//unmarshal the proc hdr
func unmarshalProcEventHdr(data []byte) procEventHdr {
	hdr := procEventHdr{}

	hdr.What = EventType(nlenc.Uint32(data[0:4]))
	hdr.CPU = nlenc.Uint32(data[4:8])
	hdr.Timestamp = nlenc.Uint64(data[8:16])

	return hdr
}

//EventType is a type for carrying around the valid list of event types
type EventType uint32

//These types are taken from cn_proc.h, and represent all the known types that the proc connector will notify on
const (

	//ProcEventNone is only used for ACK events
	ProcEventNone EventType = 0x00000000
	//ProcEventFork is a fork event
	ProcEventFork EventType = 0x00000001
	//ProcEventExec is a exec() event
	ProcEventExec EventType = 0x00000002
	//ProcEventUID is a user ID change
	ProcEventUID EventType = 0x00000004
	//ProcEventGID is a group ID change
	ProcEventGID EventType = 0x00000040
	//ProcEventSID is a session ID change
	ProcEventSID EventType = 0x00000080
	//ProcEventSID is a process trace event
	ProcEventPtrace EventType = 0x00000100
	//ProcEventComm is a comm(and) value change. Any value over 16 bytes will be truncated
	ProcEventComm EventType = 0x00000200
	//ProcEventCoredump is a core dump event
	ProcEventCoredump EventType = 0x40000000
	//ProcEventExit is an exit() event
	ProcEventExit EventType = 0x80000000
)

//ProcEvent is the struct representing all the event data that comes across the wire, in parsed form.
type ProcEvent struct {
	WhatString string    `json:"event_string"`
	What       EventType `json:"event"`
	CPU        uint32    `json:"cpu"`
	Timestamp  time.Time `json:"timestamp"`
	EventData  EventData `json:"event_data"`
}
