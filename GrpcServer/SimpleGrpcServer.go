package main

import (
	pb "20210810BFLProj/grpcSimpleService1017/proto"
	"context"
	log "fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	//1117ad
	pbtrust "20210810BFLProj/grpcSimpleService1017/GrpcClient/proto"
)
const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network string = "tcp"
)

type SimpleService struct{}
// RouteList 实现RouteList方法
func (s *SimpleService) RouteList(srv pb.StreamClient_RouteListServer) error {
	for {
		//从流中获取消息
		res, err := srv.Recv()
		if err == io.EOF {
			//发送结果，并关闭
			return srv.SendAndClose(&pb.SimpleResponse{Value: "ok====SimpleService info rpcRespval"})
		}
		if err != nil {
			return err
		}
		log.Println(res.StreamData)
	}
}
type TrustService struct{}
func (s *TrustService) TrustRequestData(ctx context.Context,trustreq *pbtrust.TrustRequest) (*pbtrust.TrustResponse,error) {
	log.Println("cur invoke rpc TrustRequestData()!")
	var retval ="5555777"
	var code int32= 44
	var score int64 =33
	var address []string
	address = append(address,"1B6011C07BAF2F69EA7627183D2F732B30680C47")
	return &pbtrust.TrustResponse{Value:&retval,Code:&code,GetRanking: &pbtrust.GetRanking{Score: &score,Address: address}},nil
}

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Println("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterStreamClientServer(grpcServer, &SimpleService{})
	//1117add
	pbtrust.RegisterTrustClientServer(grpcServer, &TrustService{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println("grpcServer.Serve err: %v", err)
	}
}