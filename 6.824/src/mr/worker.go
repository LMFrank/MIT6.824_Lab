package mr

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)
import "log"
import "net/rpc"
import "hash/fnv"

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func genWorkerID() (uuid string) {
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x\n", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}

//
// main/mrworker.go calls this function.
// 启动无限循环，不断向master发出请求
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the master.
	// CallExample()
	workerId := genWorkerID()
	retry := 3

	for {
		args := WorkArgs{WorkerId: workerId}
		reply := WorkReply{}
		working := call("Master.Work", &args, &reply)

		//log.Println(working, reply.IsFinished)

		if reply.IsFinished || !working {
			//log.Println("finised")
			return
		}
		//log.Println("task info:", reply)
		switch reply.MapReduce {
		case "map":
			MapWork(reply, mapf)
			retry = 3
		case "reduce":
			ReduceWork(reply, reducef)
			retry = 3
		default:
			//log.Println("error reply: would retry times:", retry)
			if retry < 0 {
				return
			}
			retry--
		}

		CommitArgs := CommitArgs{
			WorkerId:  workerId,
			TaskId:    reply.TaskId,
			MapReduce: reply.MapReduce,
		}
		commitReply := CommitReply{}

		_ = call("Master.Commit", &CommitArgs, &commitReply)

		//isOK := call("Master.Commit", &CommitArgs, &commitReply)
		//log.Println("Call isOK:", isOK)
		// 设置一个轮询间隔时间，降低负载
		time.Sleep(500 * time.Millisecond)
	}
}

func MapWork(task WorkReply, mapf func(string, string) []KeyValue) {
	file, err := os.Open(task.FileName)
	if err != nil {
		log.Fatalf("cannot open %v", task.FileName)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", task.FileName)
	}
	file.Close()
	kva := mapf(task.FileName, string(content))
	sort.Sort(ByKey(kva))

	tmpName := "mr-tmp-" + strconv.Itoa(task.TaskId)
	var fileNum = make(map[int]*json.Encoder)
	for i := 0; i < task.FileNumber; i++ {
		ofile, _ := os.Create(tmpName + "-" + strconv.Itoa(i))
		fileNum[i] = json.NewEncoder(ofile)
		defer ofile.Close()
	}
	for _, kv := range kva {
		key := kv.Key
		reduceIdx := ihash(key) % task.FileNumber
		err := fileNum[reduceIdx].Encode(&kv)
		if err != nil {
			log.Fatal("Unable to write to file")
		}
	}
}

func ReduceWork(task WorkReply, reducef func(string, []string) string) {
	intermediate := []KeyValue{}

	for mapTaskNumber := 0; mapTaskNumber < task.FileNumber; mapTaskNumber++ {
		filename := "mr-tmp-" + strconv.Itoa(mapTaskNumber) + "-" + strconv.Itoa(task.TaskId)
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal("Unable to read from:", filename)
		}
		defer f.Close()
		decoder := json.NewDecoder(f)
		var kv KeyValue
		for decoder.More() {
			err := decoder.Decode(&kv)
			if err != nil {
				log.Fatal("Json decode failed,", err)
			}
			intermediate = append(intermediate, kv)
		}
	}
	sort.Sort(ByKey(intermediate))

	i := 0
	ofile, err := os.Create("mr-out-" + strconv.Itoa(task.TaskId+1))
	if err != nil {
		log.Fatal("Unable to create file:", ofile)
	}
	defer ofile.Close()
	//log.Println("complete to", task.TaskId, "start to write in to", ofile.Name())
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)
		i = j
	}
	ofile.Close()
}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
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
