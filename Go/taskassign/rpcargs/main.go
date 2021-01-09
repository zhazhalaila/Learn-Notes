package rpcargs

type TaskRequest struct {
}

type TaskResponse struct {
	Index int
}

type TaskReport struct {
	Index  int
	Status int
}

type TaskReportReply struct {
}
