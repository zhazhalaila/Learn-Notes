package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"../rpcargs"
)

type Master struct {
	mu   sync.Mutex
	task []int
	done bool
}

const (
	unassign  = 0
	inprocess = 1
	completed = 2
)

func (m *Master) AssignWork(args *rpcargs.TaskRequest, reply *rpcargs.TaskResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Index = -1

	if m.done {
		log.Printf("All task done.\n")
		return nil
	}

	taskDone := true

	for index := range m.task {
		taskDone = taskDone && m.task[index] == completed
	} //检查所有任务是否完成，如果完成则不再分配任务

	if taskDone {
		log.Printf("Task done? %t\n", taskDone)
		m.done = true
		return nil
	}

	for index := range m.task {
		if m.task[index] == unassign { //检查任务是否已经分配，如果未分配则分配该任务
			m.task[index] = inprocess
			reply.Index = index
			log.Printf("Task %d has assign.\n", index)
			go m.taskWait(index)
			return nil
		}
	}

	return nil
}

func (m *Master) ReportTask(args *rpcargs.TaskReport, reply *rpcargs.TaskReportReply) error {
	m.task[args.Index] = completed
	log.Printf("Task %d has completed.\n", args.Index)
	return nil
}

func (m *Master) taskWait(index int) {
	time.Sleep(2 * time.Second) //定时检查任务完成状态
	m.mu.Lock()
	if m.task[index] != completed {
		log.Printf("Task %d error!\n", index)
		m.task[index] = unassign
	} //如果任务没有在规定时间内完成，则更改任务状态
	m.mu.Unlock()
}

func (m *Master) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(m)
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			} else {
				break
			}
		}
		l.Close()
	}()
}

func main() {
	m := Master{}
	m.task = make([]int, 10)
	m.mu = sync.Mutex{}
	m.done = false
	m.server()
	for !m.done {
		time.Sleep(1 * time.Second)
	} // 检查所有任务是否完成，如果未完成主Goroutine进行等待
}
