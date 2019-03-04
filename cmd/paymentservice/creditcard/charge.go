package creditcard

import "github.com/google/uuid"

// Charger is used to charge the credit card.
type Charger interface {
	// Charge charges the credit card and returns
	// a transaction ID.
	Charge(Card) (id string, err error)
}

// NoOps is a dummy charger.
// No charges are really proceeded.
type NoOps struct{}

// Charge always returns an UUID as transaction ID.
func (*NoOps) Charge(_ Card) (string, error) {
	id, _ := uuid.NewRandom()
	return id.String(), nil
}
