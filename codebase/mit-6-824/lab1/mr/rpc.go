package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

// args to get the number of reducer
type ReducerSizeArgs struct {
	// empty
}

// reply to get the number of reducer
// it is more efficient to request it once and stored it as a global variable
type ReducerSizeReply struct {
	ReducerSize int
}

// args to request a task
type TaskRequestArgs struct {
	// how to get the worker ID ?
	WorkerId int
}

// reply for a task request
type TaskRequestReply struct {
	TaskType      TaskType
	TaskId        int
	InputFilePath string
}

// args to report a finished task
type TaskReportArgs struct {
	TaskType TaskType
	TaskId   int // this can be globally unique
	WorkerId int
}

// reply for a task finish report
type TaskReportReply struct {
	CanExit bool // tell a worker it can exit
}

// TaskType enum
type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	NoTask
	ExitTask
	UnknownTask // this should not be used
)

func (t TaskType) String() string {
	switch t {
	case MapTask:
		return "MapTask"
	case ReduceTask:
		return "ReduceTask"
	case NoTask:
		return "NoTask"
	case ExitTask:
		return "ExitTask"
	default:
		return "UnknownTask"
	}
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
