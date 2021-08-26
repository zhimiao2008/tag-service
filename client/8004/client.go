package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/weizhimiao/tag-service/internal/middleware"
	pb "github.com/weizhimiao/tag-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {

	originCtx := context.Background()
	//ctx := metadata.AppendToOutgoingContext(originCtx, "eddycjy", "go 编程之旅")
	md := metadata.New(map[string]string{
		"go":   "programming",
		"tour": "book",
	})

	ctx := metadata.NewOutgoingContext(originCtx, md)
	clientConn, err := GetClientConn(ctx, "localhost:8004", nil)
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

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithUnaryInterceptor(
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
		),
	))
	return grpc.DialContext(ctx, target, opts...)
}
