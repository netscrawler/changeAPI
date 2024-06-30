package cnvrt

import (
	"context"
	cnvrtv1 "github.com/netscrawler/protos/gen/go/changeAPI"
	"google.golang.org/grpc"
)

type serverAPI struct {
	cnvrtv1.UnimplementedConverterServer
}

func Register(gRPC *grpc.Server) {
	cnvrtv1.RegisterConverterServer(gRPC, &serverAPI{})
}

func (s serverAPI) Convert(
	ctx context.Context,
	req *cnvrtv1.ConvertRequest) (
	*cnvrtv1.ConvertResponse, error) {
	panic("implement me")
}
