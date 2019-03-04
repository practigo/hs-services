package creditcard_test

import (
	"testing"

	cc "github.com/practigo/hipstershop/paymentservice/creditcard"
)

func TestLuhn(t *testing.T) {
	num := "79927398713"
	if !cc.Luhn(num) {
		t.Error(num + " should be valid")
	}
	// mutate one digit
	num = "79927398712"
	if cc.Luhn(num) {
		t.Error(num + " should be invalid")
	}
	num = "69927398713"
	if cc.Luhn(num) {
		t.Error(num + " should be invalid")
	}
	// actually swap any two distance-2 digits can pass
	num = "79937298713"
	if !cc.Luhn(num) {
		t.Error(num + " should be valid")
	}
	num = "99727398713"
	if !cc.Luhn(num) {
		t.Error(num + " should be valid")
	}
}

type cardTester struct {
	c *cc.Card
	t *testing.T
}

func (ct *cardTester) assert(cond bool, msg string) {
	if !cond {
		ct.t.Errorf("expected %s - error: %v (%s)", msg, ct.c.Err(), ct.c.Msg())
	}
}

func TestValidate(t *testing.T) {
	card := &cc.Card{
		Number: "4242424242424242", // from https://stripe.com/docs/testing
		Exp: cc.Expiration{
			Year:  2020,
			Month: 2,
		},
	}
	ct := &cardTester{
		c: card,
		t: t,
	}
	// working cases
	ct.assert(card.Validate(), "a valid VISA")
	card.Number = "5555555555554444"
	ct.assert(card.Validate(), "a valid MASTERCARD")
	// error cases
	card.Number = "4242424242424243"
	ct.assert(!card.Validate() && card.Err() == cc.ErrInvalid, "a invalid card")
	card.Number = "6011111111111117" // discover
	ct.assert(!card.Validate() && card.Err() == cc.ErrUnaccpted, "an unaccepted card")
	card.Number = "4242424242424242" // reset
	card.Exp.Year = 2018
	ct.assert(!card.Validate() && card.Err() == cc.ErrExpired, "an expired card")
}

func TestCharger(t *testing.T) {
	c := cc.NoOps{}
	id, err := c.Charge(cc.Card{})
	if err != nil || len(id) <= 0 {
		t.Error(id, err)
	}
}
