package proc

import "github.com/mdlayher/netlink/nlenc"

//convenience methods
func return4Uint32(data []byte) (uint32, uint32, uint32, uint32) {
	return nlenc.Uint32(data[0:4]),
		nlenc.Uint32(data[4:8]),
		nlenc.Uint32(data[8:12]),
		nlenc.Uint32(data[12:16])
}

func return2Uint32(data []byte) (uint32, uint32) {
	return nlenc.Uint32(data[0:4]),
		nlenc.Uint32(data[4:8])

}

//turn the int event type to a string
func evtType2Str(evt EventType) string {

	var evtStr string

	switch evt {
	case ProcEventNone:
		evtStr = "None"
	case ProcEventFork:
		evtStr = "Fork"
	case ProcEventExec:
		evtStr = "Exec"
	case ProcEventUID:
		evtStr = "UID"
	case ProcEventGID:
		evtStr = "GID"
	case ProcEventSID:
		evtStr = "SID"
	case ProcEventPtrace:
		evtStr = "Ptrace"
	case ProcEventComm:
		evtStr = "Command"
	case ProcEventCoredump:
		evtStr = "Core Dump"
	case ProcEventExit:
		evtStr = "Exit"

	}

	return evtStr
}

/*
===============================================================================
These are the struct defs used in cn_proc.h

*/

//EventData is an interface that encapsulates the union type used in cn_proc
//The PID and TGID fields are the only attributes shared by all the event types.
//Go get all the other fields, cast to a concrete type
type EventData interface {
	Pid() uint32

	Tgid() uint32
}

//Fork is the event for process forks
type Fork struct {
	ParentPid  uint32 `json:"parent_pid" pretty:"Parent PID"`
	ParentTgid uint32 `josn:"parent_tgid" pretty:"Parent TGID"`
	ChildPid   uint32 `json:"child_pid" pretty:"Child PID"`
	ChildTgid  uint32 `json:"child_tgid" pretty:"Child TGID"`
}

//Pid returns the event Process ID
func (f Fork) Pid() uint32 {
	return f.ChildPid
}

//Tgid returns the event thread group ID
func (f Fork) Tgid() uint32 {
	return f.ChildTgid
}

//Exec is the event for process exec()s
type Exec struct {
	ProcessPid  uint32 `json:"proces_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
}

//Pid returns the event Process ID
func (e Exec) Pid() uint32 {
	return e.ProcessPid
}

//Tgid returns the event thread group ID
func (e Exec) Tgid() uint32 {
	return e.ProcessTgid
}

//ID represents UID/GID changes for a process.
//in cn_proc.h, the real/effective GID/UID is a series of union types, which Go does not have.
//creating a super-special interface for this would be overkill,
//So we're going to rename the vars and just use two.
//Consumers should use `what` to distinguish between the two.
type ID struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
	RealID      uint32 `json:"real_id" pretty:"Real ID"`
	EffectiveID uint32 `json:"effective_id" pretty:"Effective ID"`
}

//Pid returns the event Process ID
func (i ID) Pid() uint32 {
	return i.ProcessPid
}

//Tgid returns the event thread group ID
func (i ID) Tgid() uint32 {
	return i.ProcessTgid
}

//Sid is the event for Session ID changes
type Sid struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
}

//Pid returns the event process  ID
func (s Sid) Pid() uint32 {
	return s.ProcessPid
}

//Tgid returns the event thread group ID
func (s Sid) Tgid() uint32 {
	return s.ProcessTgid
}

//Ptrace is the event for ptrace events
type Ptrace struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
	TracerPid   uint32 `json:"tracrer_pid" pretty:"Tracer PID"`
	TracerTgid  uint32 `json:"tracer_tgid" pretty:"Tracer TGID"`
}

//Pid returns the event Process ID
func (p Ptrace) Pid() uint32 {
	return p.ProcessPid
}

//Tgid returns the event thread group ID
func (p Ptrace) Tgid() uint32 {
	return p.ProcessTgid
}

//Comm represents changes to the command name, /proc/$PID/comm
type Comm struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
	Comm        string `json:"command" pretty:"Command"`
}

//Pid returns the event Process ID
func (c Comm) Pid() uint32 {
	return c.ProcessPid
}

//Tgid returns the event thread group ID
func (c Comm) Tgid() uint32 {
	return c.ProcessTgid
}

//Coredump is the event for...core dumps
type Coredump struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
}

//Pid returns the event Process ID
func (c Coredump) Pid() uint32 {
	return c.ProcessPid
}

//Tgid returns the event thread group ID
func (c Coredump) Tgid() uint32 {
	return c.ProcessTgid
}

//Exit is the event for exit()
type Exit struct {
	ProcessPid  uint32 `json:"process_pid" pretty:"Process PID"`
	ProcessTgid uint32 `json:"process_tgid" pretty:"Process TGID"`
	ExitCode    uint32 `json:"exit_code" pretty:"Exit Code"`
	ExitSignal  uint32 `json:"exit_signal" pretty:"Exit Signal"`
}

//Pid returns the event Process ID
func (e Exit) Pid() uint32 {
	return e.ProcessPid
}

//Tgid returns the event thread group ID
func (e Exit) Tgid() uint32 {
	return e.ProcessTgid
}
