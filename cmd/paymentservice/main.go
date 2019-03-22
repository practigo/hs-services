package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/practigo/hipstershop/genproto"
	"github.com/practigo/hipstershop/infra"
	cc "github.com/practigo/hs-paymentservice/creditcard"
)

const (
	defaultPort = 50051
)

func main() {
	log := infra.InitLogrus()

	addr := fmt.Sprintf(":%d", infra.AppPort(defaultPort))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := infra.NewServer(infra.WithHealth, infra.WithReflection)
	pb.RegisterPaymentServiceServer(srv, &server{
		ch: &cc.NoOps{},
		lg: log,
	})

	log.Infof("Payment Service listening on %s", addr)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// server controls RPC service responses.
type server struct {
	ch cc.Charger
	lg *logrus.Logger
}

func (s *server) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	s.lg.Info("[Charge] received request")
	defer s.lg.Info("[Charge] completed request")
	// convert
	card := cc.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    req.CreditCard.CreditCardCvv,
		Exp: cc.Expiration{
			Year:  int(req.CreditCard.CreditCardExpirationYear),
			Month: int(req.CreditCard.CreditCardExpirationMonth),
		},
	}
	id, _ := s.ch.Charge(card, req.Amount)
	return &pb.ChargeResponse{
		TransactionId: id,
	}, nil
}
