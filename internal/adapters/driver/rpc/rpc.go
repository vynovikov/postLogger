package rpc

import (
	"context"
	"net"
	"os"
	"postLogger/internal/adapters/application"
	"postLogger/internal/adapters/driver/rpc/pb"
	"postLogger/internal/logger"
	"sync"

	"google.golang.org/grpc"
)

type ReceiverStruct struct {
	A application.Application
	pb.LoggerServer
	Listener net.Listener
	Server   *grpc.Server
	l        sync.Mutex
}
type Receiver interface {
	Run()
}

func NewReceiver(a application.Application) *ReceiverStruct {

	lis, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.L.Errorf("in rpc.NewReceiver failed to listen on 3200: %v\n", err)
	}

	baseServer := grpc.NewServer()

	r := &ReceiverStruct{
		Listener: lis,
		A:        a,
	}

	pb.RegisterLoggerServer(baseServer, r)
	r.Server = baseServer

	return r
}

func (r *ReceiverStruct) Run() {
	localHost := os.Getenv("HOSTNAME")
	if len(localHost) > 0 {
		logger.L.Infof("listening %s:3200", localHost)
	} else {
		logger.L.Infof("listening localhost:3200")
	}
	r.Server.Serve(r.Listener)
}

func (r *ReceiverStruct) Log(ctx context.Context, in *pb.LogReq) (*pb.LogRes, error) {
	r.l.Lock()
	defer r.l.Unlock()
	r.A.Handle(in.Ts, in.LogString)
	return &pb.LogRes{Result: true}, nil
}
