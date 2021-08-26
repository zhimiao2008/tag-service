package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weizhimiao/tag-service/pkg/bapi"
	"github.com/weizhimiao/tag-service/pkg/errcode"
	pb "github.com/weizhimiao/tag-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TagServer struct {
	auth *Auth
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	// auth 认证
	if err := t.auth.Check(ctx); err != nil {
		fmt.Printf("auth fail, err: %v \n", err)
		return nil, err
	}

	//panic("异常抛出 测试")

	// read metadata
	md,_ := metadata.FromIncomingContext(ctx)
	fmt.Printf("metadata.FromIncomingContext md: %+v \n", md)

	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ERROR_GET_TAG_LIST_FAIL)
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}
	return &tagList, nil
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
