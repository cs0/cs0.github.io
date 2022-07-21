package mr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const WAIT_TIME time.Duration = 500

// number of reducers, so we can partition the temp file
var nReducer int = 0

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
// here the mrworker pass in the mapf and reducef
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// get reducer size
	reducerSizeReply, succ := requestReducerSize()
	if !succ {
		fmt.Println("Failed to get reducer count, worker exiting.")
		return
	}
	nReducer = reducerSizeReply.ReducerSize

	// a loop to keep getting task and finishing task until the coordinator tell us to stop
	for {
		// coordinate will use worker id in the request as task id?
		assignedTask, succ := requestTask()

		if !succ {
			// it can be because the coordinate exit already
			// though it can also be due to network issue, here we directly exit
			fmt.Println("Failed to get task from coordinator, worker exiting.")
			return
		}
		fmt.Printf("Worker %v get task %v with type %v\n", getWorkerId(), assignedTask.TaskId, assignedTask.TaskType)

		canExit := false
		switch assignedTask.TaskType {
		case MapTask:
			doMapTask(mapf, assignedTask.InputFilePath, assignedTask.TaskId)
			canExit, succ = reportTaskDone(assignedTask.TaskType, assignedTask.TaskId)
		case ReduceTask:
			doReduceTask(reducef, assignedTask.TaskId)
			canExit, succ = reportTaskDone(assignedTask.TaskType, assignedTask.TaskId)
		case NoTask:
			// sleep for a while and wait for the task
			// if can exit, will exit
			// TODO we can use the channel, or the condition variable to make it better
		case ExitTask:
			// exit the loop
			fmt.Println("All tasks are done, worker exiting")
			return
		}

		if canExit {
			fmt.Println("Coordinate told the worker that it can exit. working exiting.")
			return
		}
		if !succ {
			fmt.Println("Fail to report that the task is done. working exiting.")
			return
		}

		// here sleep for an interval for each around
		time.Sleep(WAIT_TIME)
	}
}

func requestReducerSize() (*ReducerSizeReply, bool) {
	args := ReducerSizeArgs{}
	reply := ReducerSizeReply{}
	succ := call("Coordinator.RequestReducerSize", &args, &reply)
	return &reply, succ
}

func requestTask() (*TaskRequestReply, bool) {
	args := TaskRequestArgs{}
	args.WorkerId = getWorkerId()
	reply := TaskRequestReply{}
	succ := call("Coordinator.RequestTask", &args, &reply)
	return &reply, succ
}

func getWorkerId() int {
	// TODO we can register worker id in the coordinate from 0
	return os.Getpid()
}

func getReducerId(key string) int {
	return ihash(key) % nReducer
}

// input: input file path
func doMapTask(mapf func(string, string) []KeyValue, inputFilePath string, mapId int) {
	// call the mapf and write the output to a temp file
	// return the temp file to coordinate
	file, err := os.Open(inputFilePath)
	checkError(err, "Cannot open file %v\n", inputFilePath)

	content, err := ioutil.ReadAll(file)
	checkError(err, "Cannot read file %v\n", inputFilePath)
	// instead of defer, we close the file aggressively
	file.Close()

	kvPairs := mapf(inputFilePath, string(content))
	// the output path is auto generated
	writeMapOutput(kvPairs, mapId)
}

func checkError(err error, format string, v ...interface{}) {
	if err != nil {
		log.Fatalf(format, v)
	}
}

// A reasonable naming convention for intermediate files is mr-X-Y,
// where X is the Map task number, and Y is the reduce task number.
// Since the output needs to partitioned for reducers
func writeMapOutput(pairs []KeyValue, mapId int) {
	// use io buffers to reduce disk I/O, which greatly improves
	// performance when running in containers with mounted volumes
	prefix := fmt.Sprintf("%v/mr-%v", TempDir, mapId)
	files := make([]*os.File, 0, nReducer)
	buffers := make([]*bufio.Writer, 0, nReducer)
	encoders := make([]*json.Encoder, 0, nReducer)

	// create temp files, the contract is reducer id starts from 0
	// use the worker id to distinguish
	for i := 0; i < nReducer; i++ {
		filePath := fmt.Sprintf("%v-%v-%v", prefix, i, getWorkerId())
		file, err := os.Create(filePath)
		checkError(err, "Cannot create file %v \n", filePath)
		buf := bufio.NewWriter(file)
		files = append(files, file)
		buffers = append(buffers, buf)
		// for encoder, instead of build it on a file, we build it on a buffer
		encoders = append(encoders, json.NewEncoder(buf))
	}

	// write map outputs to temp files
	for _, p := range pairs {
		reducerId := getReducerId(p.Key)
		err := encoders[reducerId].Encode(p)
		checkError(err, "Cannot encode % v to file \n", p)
	}

	// flush file buffer to disk
	for i, buf := range buffers {
		err := buf.Flush()
		checkError(err, "Cannot flush buffer for file: %v\n", files[i].Name())
	}

	// atomatically rename temp file
	for i, file := range files {
		file.Close()                               // need to close first
		newPath := fmt.Sprintf("%v-%v", prefix, i) // output file path "TempDir/mr-mapId-reduceId"
		err := os.Rename(file.Name(), newPath)
		checkError(err, "Cannot rename file %v \n", file.Name())
	}
}

// reducerId
func doReduceTask(reducef func(string, []string) string, reducerId int) {
	kvMap := readReduceFiles(reducerId)

	// sort the kv map by key
	keys := make([]string, 0, len(kvMap))
	for k := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create temp file "TempDir/mr-out-reduceId-workerId"
	filePath := fmt.Sprintf("%v/mr-out-%v-%v", TempDir, reducerId, getWorkerId())
	file, err := os.Create(filePath)
	checkError(err, "Cannot create file %v\n", filePath)

	// Call reduce and write to temp file
	for _, k := range keys {
		v := reducef(k, kvMap[k])
		_, err := fmt.Fprintf(file, "%v %v\n", k, v)
		checkError(err, "Cannot write mr output (%v, %v) to file", k, v)
	}

	// atomatically rename temp file
	file.Close()
	//
	newPath := fmt.Sprintf("mr-out-%v", reducerId)
	err = os.Rename(filePath, newPath)
	checkError(err, "Cannot rename file %v\n", filePath)
}

func readReduceFiles(reducerId int) map[string][]string {
	// read all the files that belongs to this reducer
	files, err := filepath.Glob(fmt.Sprintf("%v/mr-%v-%v", TempDir, "*", reducerId))
	checkError(err, "Cannot list reduce files")

	// each key can have multiple values, here we need to aggregate them together
	kvMap := make(map[string][]string)
	var kv KeyValue

	for _, filePath := range files {
		file, err := os.Open(filePath)
		checkError(err, "Cannot open file %v\n", filePath)

		dec := json.NewDecoder(file)
		for dec.More() {
			err = dec.Decode(&kv) // read the value and assign it to kv
			checkError(err, "Cannot decode from file %v\n", filePath)
			kvMap[kv.Key] = append(kvMap[kv.Key], kv.Value)
		}
	}
	return kvMap
}

// output1: a flag indicate if the worker can exit
// output2: a flag indicate if the report succeeds
func reportTaskDone(taskType TaskType, taskId int) (bool, bool) {
	// here ideally we need to report again if the report fails
	// but for the sake of simplicity, we do not use it
	args := TaskReportArgs{taskType, taskId, getWorkerId()}
	reply := TaskReportReply{}
	succ := call("Coordinator.ReportTaskDone", &args, &reply)
	return reply.CanExit, succ
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
