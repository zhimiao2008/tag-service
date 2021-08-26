package main

import (
	"context"
	pb "github.com/weizhimiao/tag-service/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {

	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, "localhost:8001", nil)
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
	return grpc.DialContext(ctx, target, opts...)
}
