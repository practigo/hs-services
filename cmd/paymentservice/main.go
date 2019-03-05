package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/practigo/hipstershop/genproto"
	cc "github.com/practigo/hipstershop/paymentservice/creditcard"
)

const (
	defaultPort = "50051"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

var charger cc.Charger

func main() {
	port := defaultPort
	if value, ok := os.LookupEnv("APP_PORT"); ok {
		port = value
	}
	port = fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	svc := &server{}
	pb.RegisterPaymentServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, health.NewServer())
	log.Infof("Payment Service listening on port %s", port)

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	charger = &cc.NoOps{} // TODO: so far it is just a mock

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// server controls RPC service responses.
type server struct{}

func (s *server) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	log.Info("[Charge] received request")
	defer log.Info("[Charge] completed request")
	// convert
	card := cc.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    req.CreditCard.CreditCardCvv,
		Exp: cc.Expiration{
			Year:  int(req.CreditCard.CreditCardExpirationYear),
			Month: int(req.CreditCard.CreditCardExpirationMonth),
		},
	}
	id, _ := charger.Charge(card, req.Amount)
	return &pb.ChargeResponse{
		TransactionId: id,
	}, nil
}
