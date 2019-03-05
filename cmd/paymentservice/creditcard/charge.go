package creditcard

import (
	"github.com/google/uuid"

	pb "github.com/practigo/hipstershop/genproto"
)

// Charger is used to charge the credit card.
type Charger interface {
	// Charge charges the money from the credit card
	// and returns a transaction ID.
	Charge(Card, *pb.Money) (id string, err error)
}

// NoOps is a dummy charger.
// No charges are really proceeded.
type NoOps struct{}

// Charge always returns an UUID as transaction ID.
func (*NoOps) Charge(_ Card, _ *pb.Money) (string, error) {
	id, _ := uuid.NewRandom()
	return id.String(), nil
}
