package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/naming"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/weizhimiao/tag-service/global"
	"github.com/weizhimiao/tag-service/internal/middleware"
	"github.com/weizhimiao/tag-service/pkg/tracer"
	pb "github.com/weizhimiao/tag-service/proto"
	"github.com/weizhimiao/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

const SERVICE_NAME = "tag-service"

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("setupTracer err: %v", err)
	}
}

func main() {

	auth := server.Auth{
		AppKey:    "appkeysss",
		AppSecret: "secretsecret",
	}

	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	//_ = []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}

	clientConn, err := GetClientConn(ctx, SERVICE_NAME, opts)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tarServiceClient.GetTagList err: %v", err)
	}

	log.Printf("res: %v", resp)
}

func GetClientConn(ctx context.Context, serviceName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	config := clientv3.Config{
		Endpoints: []string{
			"http://127.0.0.1:2379",
		},
		DialTimeout: 60 * time.Second,
	}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	r := &naming.GRPCResolver{Client: etcdClient}

	target := fmt.Sprintf("/etcdv3://go-grogramming-tour/grpc/%s", serviceName)

	opts = append(opts, grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(r)),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_middleware.ChainUnaryClient(middleware.UnaryContextTimeout()),
				grpc_retry.UnaryClientInterceptor(
					grpc_retry.WithMax(2),
					grpc_retry.WithCodes(
						codes.Unknown,
						codes.Internal,
						codes.DeadlineExceeded,
					),
				),
				middleware.ClientTracing(),
			),
		))
	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("tag-service-client", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
