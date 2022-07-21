package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// the temp output directory
const TempDir = "tmp"

// task timeout, 10 seconds
const TaskTimeout time.Duration = 10000

// Q: what is the best way to timeout ?

type TaskStatus int

const (
	NotStarted TaskStatus = iota
	Executing
	Finished
)

type Task struct {
	Type          TaskType
	Status        TaskStatus
	Index         int
	InputFileName string
	WorkerId      int
}

type Coordinator struct {
	totalReducerSize int // the total reducer size

	mu sync.Mutex
	// need to lock before read or write any following fields
	mapTasks        []Task
	reduceTasks     []Task
	leftMapTasks    int // the number of left map tasks
	leftReduceTasks int // the number of left reduce tasks
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
// func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
// 	reply.Y = args.X + 1
// 	return nil
// }

// expose to worker
func (c *Coordinator) RequestTask(args *TaskRequestArgs, reply *TaskRequestReply) error {
	// only assign reduce task when all map tasks are FINISHED

	c.mu.Lock()

	var task *Task
	if c.leftMapTasks > 0 {
		// assign map task
		task = c.selectTaskNotThreadSafe(c.mapTasks, args.WorkerId)
	} else if c.leftReduceTasks > 0 {
		// assign reduce task
		task = c.selectTaskNotThreadSafe(c.reduceTasks, args.WorkerId)
	} else {
		// assign exit task
		task = &Task{ExitTask, NotStarted, -1, "", -1}
	}

	reply.TaskType = task.Type
	reply.TaskId = task.Index
	reply.InputFilePath = task.InputFileName

	c.mu.Unlock() // the critical lock might be too large, it can be placed in the selectTask func

	go c.waitForTask(task)

	return nil
}

func (c *Coordinator) RequestReducerSize(args *ReducerSizeArgs, reply *ReducerSizeReply) error {
	reply.ReducerSize = c.totalReducerSize
	return nil
}

func (c *Coordinator) ReportTaskDone(args *TaskReportArgs, reply *TaskReportReply) error {

	c.updateTaskStatus(args.TaskType, args.TaskId, args.WorkerId)

	reply.CanExit = c.checkFinishState()

	return nil
}

func (c *Coordinator) updateTaskStatus(taskType TaskType, taskId int, workerId int) {
	var task *Task
	if taskType == MapTask {
		task = &c.mapTasks[taskId]
	} else {
		task = &c.reduceTasks[taskId]
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// workers can only report task done if the task was not re-assigned due to timeout
	if taskId == task.Index && task.Status == Executing {
		task.Status = Finished
		if taskType == MapTask && c.leftMapTasks > 0 {
			c.leftMapTasks--
		} else if taskType == ReduceTask && c.leftReduceTasks > 0 {
			c.leftReduceTasks--
		}
	}
}

// select a task that is NotStarted
// This is not a thread safe impl so need to get lock before calling
func (c *Coordinator) selectTaskNotThreadSafe(taskList []Task, workerId int) *Task {
	var task *Task

	for i := 0; i < len(taskList); i++ {
		if taskList[i].Status == NotStarted {
			task = &taskList[i]     // get the pointer
			task.Status = Executing // change to executing
			task.WorkerId = workerId
			return task
		}
	}

	// if no task, then return a noTask
	// the worker will wait for a while and request again
	return &Task{NoTask, Finished, -1, "", -1}
}

func (c *Coordinator) checkFinishState() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.leftMapTasks == 0 && c.leftReduceTasks == 0
}

func (m *Coordinator) waitForTask(task *Task) {
	// task type won't be changed after init
	// so no lock is needed here
	if task.Type != MapTask && task.Type != ReduceTask {
		return
	}

	<-time.After(time.Second * TaskTimeout)

	m.mu.Lock()
	defer m.mu.Unlock()

	if task.Status == Executing {
		task.Status = NotStarted // change back to notStarted
		task.WorkerId = -1
		fmt.Println("Task timeout, reset task status: ", *task)
	}
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	fmt.Println("Starting coordinator.")
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	return c.checkFinishState()
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	fmt.Println("Initializing coordinator.")
	c := Coordinator{}

	c.totalReducerSize = nReduce
	nMap := len(files)
	c.leftMapTasks = nMap
	c.leftReduceTasks = nReduce
	c.mapTasks = make([]Task, 0, nMap)
	c.reduceTasks = make([]Task, 0, nReduce)

	for i := 0; i < nMap; i++ {
		mTask := Task{MapTask, NotStarted, i, files[i], -1}
		c.mapTasks = append(c.mapTasks, mTask)
	}

	for i := 0; i < nReduce; i++ {
		rTask := Task{ReduceTask, NotStarted, i, "NullFile", -1}
		c.reduceTasks = append(c.reduceTasks, rTask)
	}

	c.server()

	c.cleanUp()

	return &c
}

func (c *Coordinator) cleanUp() {
	fmt.Println("Cleaning up temp dir from previous running.")
	// clean up and create temp directory
	outFiles, _ := filepath.Glob("mr-out*")
	for _, f := range outFiles {
		if err := os.Remove(f); err != nil {
			log.Fatalf("Cannot remove file %v\n", f)
		}
	}

	err := os.RemoveAll(TempDir)
	if err != nil {
		log.Fatalf("Cannot remove temp directory %v\n", TempDir)
	}
	err = os.Mkdir(TempDir, 0755)
	if err != nil {
		log.Fatalf("Cannot create temp directory %v\n", TempDir)
	}
}
