package cnvrt

import (
	"context"
	cnvrtv1 "github.com/netscrawler/protoss/gen/go/changeAPI"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Converter interface {
	Convert(ctx context.Context,
		amount uint32,
		targetCurrency string,
	) (convertedAmount uint32, rate float32, err error)
}

type serverAPI struct {
	cnvrtv1.UnimplementedConverterServer
	convert Converter
}

func Register(gRPC *grpc.Server, convert Converter) {
	cnvrtv1.RegisterConverterServer(gRPC, &serverAPI{convert: convert})
}

func (s *serverAPI) Convert(
	ctx context.Context,
	req *cnvrtv1.ConvertRequest) (
	*cnvrtv1.ConvertResponse, error) {
	if !isValidCurrency(req.GetTargetCurrency()) {
		return nil, status.Error(codes.InvalidArgument, "Invalid target currency")
	}

	convertedAmount, rate, err := s.convert.Convert(ctx, req.GetAmount(), req.GetTargetCurrency())
	if err != nil {
		//todo error handler
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return &cnvrtv1.ConvertResponse{
		BaseAmount:        req.GetAmount(),
		ConvertedAmount:   convertedAmount,
		ConvertedCurrency: req.GetTargetCurrency(),
		Rate:              rate,
	}, nil
}

func isValidCurrency(currency string) bool {

	currencies := map[string]bool{
		"AUD": true, "GBP": true, "BYR": true, "DKK": true, "USD": true, "EUR": true,
		"ISK": true, "KZT": true, "CAD": true, "NOK": true, "XDR": true, "SGD": true,
		"TRL": true, "UAH": true, "SEK": true, "CHF": true, "JPY": true,
	}
	if currencies[currency] {
		return true
	}
	return false
}
