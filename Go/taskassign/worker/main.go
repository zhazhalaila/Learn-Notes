package main

import (
	"log"
	"math/rand"
	"net/rpc"
	"time"

	"../rpcargs"
)

func randBool() bool {
	return rand.Int()%2 == 0
}

func connect() *rpc.Client {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

func main() {
	client := connect()
	log.Printf("Connect!\n")
	for {
		args := rpcargs.TaskRequest{}
		reply := rpcargs.TaskResponse{}
		client.Call("Master.AssignWork", &args, &reply)
		log.Printf("Index : %d\n", reply.Index)
		if reply.Index == -1 {
			log.Printf("All Task Done!\n")
			break
		}
		report := randBool()
		if report {
			reportArgs := rpcargs.TaskReport{}
			replyArgs := rpcargs.TaskReportReply{}
			reportArgs.Index = reply.Index
			client.Call("Master.ReportTask", &reportArgs, &replyArgs)
		}
		time.Sleep(3 * time.Second)
	}
	client.Close()
}
