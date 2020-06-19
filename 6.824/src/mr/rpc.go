package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

// WorkerId
type WorkArgs struct {
	WorkerId string
}

// 任务信息
type WorkReply struct {
	IsFinished bool // 任务是否完成
	TaskId     int
	FileName   string
	MapReduce  string // 当前任务是map还是reduce
	FileNumber int    // 任务应读取和输出多少个文件
}

type CommitArgs struct {
	WorkerId  string
	TaskId    int
	MapReduce string
}

type CommitReply struct {
	isOK bool
}

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
